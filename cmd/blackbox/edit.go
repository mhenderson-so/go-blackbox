package main

import (
	"fmt"
	"io/ioutil"
	"os"

	blackbox "github.com/mhenderson-so/go-blackbox/cmd/go-blackbox"
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
		filename = stripGpgExtension(filename)
		if _, err := os.Stat(filename); err == nil {
			return fmt.Errorf("SKIPPING: %s Will not overwrite existing files", filename)
		}
		content, err := preworkDecode(filename, passphrase)
		if err != nil {
			return err
		}
		ioutil.WriteFile(filename, content, 0644)
	}
	return nil
}

func editEnd(args cli.Args) error {
	if len(args) == 0 { //If there are no filenames specified, error
		return fmt.Errorf("No filenames specified")
	}

	//Attempt to start encryption of each file in turn
	for _, filename := range args {
		plainFile := stripGpgExtension(filename)
		gpgFile := fmt.Sprintf("%s.gpg", plainFile)
		setup(plainFile)

		//Read our plaintext file in
		plaintext, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		//Encode it
		encodedbytes, err := blackbox.Encode(plaintext)
		if err != nil {
			return err
		}

		//Write it out
		err = ioutil.WriteFile(gpgFile, encodedbytes, 0644)
		if err != nil {
			return err
		}

		//Clean up
		//os.Remove(plainFile)
	}

	return nil
}
