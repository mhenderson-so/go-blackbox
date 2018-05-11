package main

import (
	"fmt"

	"github.com/urfave/cli"
	blackbox "github.ds.stackexchange.com/mhenderson/go-blackbox/cmd/go-blackbox"
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
