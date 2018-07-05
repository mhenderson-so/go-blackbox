package main

import (
	"bufio"
	"fmt"
	"os"

	blackbox "github.com/mhenderson-so/go-blackbox/cmd/go-blackbox"
)

func initialize() error {
	fmt.Print("Enable blackbox for this repo? (yes/no) ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	//Why do we do it this way? Only to maintain compatibility with the existing
	//blackbox functionality, as it works this way
	if text != "yes" {
		if text == "no" || text == "n" || text == "" {
			return nil
		}
	}

	//Try to initialize
	cwd, _ := os.Getwd()
	err := blackbox.Initialize(cwd)
	if err != nil {
		fmt.Println("Unable to initialize blackbox:", err)
		return err
	}

	//If we initialized succesfully then give the user some feedback
	setup("")
	fmt.Println("VCS_TYPE:", blackbox.VCSType)
	fmt.Println()
	fmt.Println("You need to manually check in your new blackbox files")

	return nil
}
