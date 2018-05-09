package main

import (
	"os"

	blackbox "github.ds.stackexchange.com/mhenderson/go-blackbox/cmd/go-blackbox"
)

func initialize() error {
	cwd, _ := os.Getwd()
	return blackbox.Initialize(cwd)
}
