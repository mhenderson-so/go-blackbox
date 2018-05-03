package blackbox

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/howeyc/gopass"
	homedir "github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/openpgp"
)

// Decode takes a blackbox-encrypted payload and decrypts it with the users private key.
// It loads the users `.gnupg/secring.gpg` file and uses the passphrase to attempt to
// decode the private keys contained in the keyring.
func Decode(filepath, passphrase string) ([]byte, error) {
	keyringFileBuffer, err := os.Open(privateKeyringPath())
	if err != nil {
		return nil, err
	}
	defer keyringFileBuffer.Close()
	entityList, err := openpgp.ReadKeyRing(keyringFileBuffer)
	if err != nil {
		return nil, err
	}

	var passbytes []byte
	if passphrase != "" {
		passbytes = []byte(passphrase)
	} else {
		fmt.Printf("Enter gpg passphrase: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			return nil, err
		}
		passbytes = []byte(pass)
	}

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	md, err := openpgp.ReadMessage(file, entityList, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		for _, key := range keys {
			if key.PrivateKey != nil {
				err := key.PrivateKey.Decrypt(passbytes)
				if err == nil {
					for idName := range key.Entity.Identities {
						fmt.Println("Found private key for", idName)
					}
				}
			}
		}
		return nil, nil
	}, nil)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func publicKeyringPath() string {
	return filepath.Join(blackboxData, "pubring.gpg")
}

func privateKeyringPath() string {
	home, _ := homedir.Dir()
	return filepath.Join(home, ".gnupg", "secring.gpg")
}
