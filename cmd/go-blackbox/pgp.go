package blackbox

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/openpgp"
)

// Decode takes a blackbox-encrypted payload and decrypts it with the users private key.
// It loads the users `.gnupg/secring.gpg` file and uses the passphrase to attempt to
// decode the private keys contained in the keyring.
func Decode(filepath string, passphrase []byte) ([]byte, error) {
	keyringFileBuffer, err := os.Open(PrivateKeyringPath())
	if err != nil {
		return nil, err
	}
	defer keyringFileBuffer.Close()
	entityList, err := openpgp.ReadKeyRing(keyringFileBuffer)
	if err != nil {
		return nil, err
	}

	if len(passphrase) == 0 {
		return nil, fmt.Errorf("No passphrase supplied")
	}

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	md, err := openpgp.ReadMessage(file, entityList, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		for _, key := range keys {
			if key.PrivateKey != nil {
				key.PrivateKey.Decrypt(passphrase)
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

// PublicKeyringPath returns the path to the users public keyring
func PublicKeyringPath() string {
	return filepath.Join(blackboxData, "pubring.gpg")
}

// PrivateKeyringPath returns the path to the users private keyring
func PrivateKeyringPath() string {
	home, _ := homedir.Dir()
	return filepath.Join(home, ".gnupg", "secring.gpg")
}
