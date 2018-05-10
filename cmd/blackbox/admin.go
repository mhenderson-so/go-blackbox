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

	newkeys, err := blackbox.AdminAdd(email, path)
	if err != nil {
		fmt.Println(err)
		return err
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
