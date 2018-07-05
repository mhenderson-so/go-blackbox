package main

import (
	"fmt"

	blackbox "github.com/mhenderson-so/go-blackbox/cmd/go-blackbox"
	"github.com/urfave/cli"
)

func register(args cli.Args) error {
	if len(args) == 0 { //If there are no filenames specified, error
		return fmt.Errorf("No filenames specified")
	}

	//Attempt to add each file
	for _, filename := range args {
		setup(filename)
		err := blackbox.RegisterNewFile(filename)
		if err != nil {
			return err
		}
	}
	return nil
}
