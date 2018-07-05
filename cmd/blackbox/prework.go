package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/mhenderson-so/gpgagent"
	blackbox "github.com/mhenderson-so/go-blackbox/cmd/go-blackbox"
	"golang.org/x/crypto/openpgp"
)

// preworkPassphrase provides the passphrase that will be passed into go-blackbox later for
// decryption. It starts to attempting a gpg-agent connection. If that succeeds, use the gpg-agent
// connection. If it doesn't succeed (not running, not supported) then prompt for a passphrase
// on the console.
func preworkPassphrase() ([]byte, error) {
	//Attempt to decode via gpg-agent first.
	conn, err := gpgagent.NewGpgAgentConn()
	if err == nil {
		defer conn.Close()           //We have a valid gpg-agent connection
		return preworkGPGAgent(conn) //So use it to try and find a passphrase
	}

	//We do not have a gpg-agent connection at this stage, so manually prompt
	//for the passphrase
	fmt.Printf("Enter gpg passphrase: ")
	pass, err := gopass.GetPasswd()
	if err != nil {
		return nil, err
	}
	return []byte(pass), nil
}

// preworkGPGAgent is where the gpg-agent is used to fetch a passphrase. It goes through all
// the private keys in the users gnupg home directory to find ones that have private keys. It
// then requests an unlock from gpg-agent. Once it finds a key it can unlock, it returns the
// passphrase that the user entered. gpg-agent might ask for the passphrase interactively,
// or it might return a stored passphrase from a previous authentication.
func preworkGPGAgent(conn *gpgagent.Conn) ([]byte, error) {
	keyringFileBuffer, err := os.Open(blackbox.PrivateKeyringPath()) //Open the users private keyvault
	if err != nil {
		return nil, err
	}
	defer keyringFileBuffer.Close()
	entityList, err := openpgp.ReadKeyRing(keyringFileBuffer) //Read the keys out of the keyvault
	if err != nil {
		return nil, err
	}

	for _, key := range entityList.DecryptionKeys() { //Find the keys that can be used for decryption
		cacheID := strings.ToUpper(hex.EncodeToString(key.PublicKey.Fingerprint[:])) //Cache ID is used for gpg-agent to store/retrieve the passphrase
		var names []string
		for _, ident := range key.Entity.Identities { //This is for the prompt - fetch the list of identities tied to this certificate
			names = append(names, ident.Name)
		}
		//Our prompt message so the user knows which key we are trying to decrypt
		desc := fmt.Sprintf("You need a passphrase to unlock the secret key for user:\n%s", strings.Join(names, "\n"))

		//Build a request to pass to gpg-agent
		request := gpgagent.PassphraseRequest{
			CacheKey: cacheID,
			Desc:     desc,
			Prompt:   "Passphrase:",
		}

		//Pass the request to gpg-agent. If we fail for any reason, remove this cached passphrase. Otherwise
		//gpg-agent will continually return it, even if it's a useless passphrase.
		passphrase, err := conn.GetPassphrase(&request)
		if err != nil {
			conn.RemoveFromCache(cacheID)
			return nil, err
		}
		err = key.PrivateKey.Decrypt([]byte(passphrase)) //Use this passphrase to attempt key decryption
		if err != nil {
			conn.RemoveFromCache(cacheID)
			return nil, err
		}
		return []byte(passphrase), nil
	}
	return nil, fmt.Errorf("Unable to find key")
}

// preworkDecode actually takes care of decoding the file with go-blackbox. It does the checking
// of the file extensions, makes sure it's actually an encrypted file, and then passes back the
// byte array containing the contents of the decrypted file. It expects a passphrase already, which
// should have been received by preworkPassphrase().
func preworkDecode(filename string, passphrase []byte) ([]byte, error) {
	filename = stripGpgExtension(filename) //Tab autocompletion will put a .gpg extension on the end of a filename, but that's not what we work on in the keyring file
	setup(filepath.Dir(filename))          //We will run setup() on each item as they could possibly be in different repos/keyrings.

	if !blackbox.IsOnCryptList(filename) { //Check if the file we're trying to unblackbox is listed as a file in the keyring
		log.Fatalf("%s is not a blackbox encrypted file", filename)
	}

	gpgFilename := fmt.Sprintf("%s.gpg", filename)
	encryptedContent, err := ioutil.ReadFile(gpgFilename)
	if err != nil {
		return nil, err
	}

	contents, err := blackbox.Decode(encryptedContent, passphrase)
	if err != nil {
		return nil, fmt.Errorf("Error decoding file: %s", err)
	}
	return contents, nil
}
