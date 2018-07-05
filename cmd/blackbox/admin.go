package main

import (
	"fmt"

	blackbox "github.com/mhenderson-so/go-blackbox/cmd/go-blackbox"
	"github.com/urfave/cli"
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

func admincleanup() error {
	setup("")
	removed, err := blackbox.AdminCleanup()
	if err != nil {
		return err
	}
	if len(removed) > 0 {
		fmt.Println("Orphaned admins removed:")
		for _, admin := range removed {
			fmt.Println("    ", admin)
		}
	}

	return nil
}
