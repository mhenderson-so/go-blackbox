package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/urfave/cli"
)

// editStart decrypts a list of files. They will need to be re-encrypted manually.
func editStart(args cli.Args) error {
	if len(args) == 0 { //If there are no filenames specified, error
		return fmt.Errorf("No filenames specified")
	}

	passphrase, err := preworkPassphrase()
	if err != nil {
		return err
	}

	//Attempt to start decryption of each file in turn
	for _, filename := range args {
		if _, err := os.Stat(filename); err == nil {
			return fmt.Errorf("SKIPPING: %s Will not overwrite non-empty files", filename)
		}
		content, err := preworkDecode(filename, passphrase)
		if err != nil {
			return err
		}
		ioutil.WriteFile(filename, content, 0644)
	}
	return nil
}