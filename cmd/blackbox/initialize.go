package main

import (
	"bufio"
	"fmt"
	"os"

	blackbox "github.ds.stackexchange.com/mhenderson/go-blackbox/cmd/go-blackbox"
)

func initialize() error {
	cwd, _ := os.Getwd()
	fmt.Print("Enable blackbox for this repo? (yes/no) ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	if text != "yes" {
		if text == "no" || text == "n" || text == "" {
			return nil
		}
	}

	err := blackbox.Initialize(cwd)

	if err != nil {
		fmt.Println("Unable to initialize blackbox:", err)
	}

	setup("")
	fmt.Println("VCS_TYPE:", blackbox.VCSType)
	fmt.Println()
	fmt.Println("You need to manually check in your new blackbox files")
	return nil
}
