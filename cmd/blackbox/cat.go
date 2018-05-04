package main

import (
	"fmt"

	"github.com/urfave/cli"
)

func cat(args cli.Args) error {
	if len(args) == 0 { //If there are no filenames specified, error
		return fmt.Errorf("No filenames specified")
	}

	passphrase, err := preworkPassphrase()
	if err != nil {
		return err
	}

	//Attempt to start decryption of each file in turn
	for _, filename := range args {
		content, err := preworkDecode(filename, passphrase)
		if err != nil {
			return err
		}
		fmt.Printf(string(content))
	}
	return nil
}
