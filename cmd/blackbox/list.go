package main

import (
	"fmt"

	blackbox "github.com/mhenderson-so/go-blackbox/cmd/go-blackbox"
)

func listFiles() error {
	setup("")
	files, err := blackbox.ListFiles()
	if err != nil {
		return err
	}
	for _, line := range files {
		fmt.Println(line)
	}
	return nil
}

func listAdmins() error {
	setup("")
	files, err := blackbox.ListAdmins()
	if err != nil {
		return err
	}
	for _, line := range files {
		fmt.Println(line)
	}
	return nil
}
