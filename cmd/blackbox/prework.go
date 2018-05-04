package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/howeyc/gopass"
	blackbox "github.ds.stackexchange.com/mhenderson/go-blackbox/cmd/go-blackbox"
)

func preworkPassphrase() ([]byte, error) {
	fmt.Printf("Enter gpg passphrase: ")
	pass, err := gopass.GetPasswd()
	if err != nil {
		return nil, err
	}
	return []byte(pass), nil
}

func preworkDecode(filename string, passphrase []byte) ([]byte, error) {
	filename = stripGpgExtension(filename) //Tab autocompletion will put a .gpg extension on the end of a filename, but that's not what we work on in the keyring file
	setup(filepath.Dir(filename))          //We will run setup() on each item as they could possibly be in different repos/keyrings.

	if !blackbox.IsOnCryptList(filename) { //Check if the file we're trying to unblackbox is listed as a file in the keyring
		log.Fatalf("%s is not a blackbox encrypted file", filename)
	}

	gpgFilename := fmt.Sprintf("%s.gpg", filename)
	contents, err := blackbox.Decode(gpgFilename, passphrase)
	if err != nil {
		return nil, fmt.Errorf("Error decoding file: %s", err)
	}
	return contents, nil
}
