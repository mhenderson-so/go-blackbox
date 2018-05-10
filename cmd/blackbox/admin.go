package main

import (
	"fmt"

	"github.com/urfave/cli"
	blackbox "github.ds.stackexchange.com/mhenderson/go-blackbox/cmd/go-blackbox"
)

func addadmin(args cli.Args) error {
	setup("")
	email := args.Get(0)
	path := args.Get(1)

	if email == "" {
		return fmt.Errorf("No keyname (email address) provided")
	}

	newkeys, err := blackbox.AdminAdd(email, path)
	if err != nil {
		return err
	}
	if len(newkeys) == 0 {
		fmt.Println("No changes made")
		return nil
	}
	fmt.Println("Imported:")
	for _, newkey := range newkeys {
		fmt.Println("   ", newkey)
	}
	return nil
}

func removeadmin(args cli.Args) error {
	return nil
}
