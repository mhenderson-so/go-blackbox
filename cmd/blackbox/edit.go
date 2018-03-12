package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/urfave/cli"
	blackbox "github.ds.stackexchange.com/mhenderson/go-blackbox/cmd/go-blackbox"
)

// editStart decrypts a list of files. They will need to be re-encrypted manually.
func editStart(args cli.Args) error {
	if len(args) == 0 { //If there are no filenames specified, error
		return fmt.Errorf("No filenames specified")
	}

	//Attempt to start decryption of each file in turn
	for _, filename := range args {
		filename = stripGpgExtension(filename) //Tab autocompletion will put a .gpg extension on the end of a filename, but that's not what we work on in the keyring file
		setup(filepath.Dir(filename))          //We will run setup() on each item as they could possibly be in different repos/keyrings.

		if !blackbox.IsOnCryptList(filename) { //Check if the file we're trying to unblackbox is listed as a file in the keyring
			log.Fatalf("%s is not a blackbox encrypted file", filename)
		}
		fmt.Println(filename, "is a blackbox file")
	}
	return nil
}
