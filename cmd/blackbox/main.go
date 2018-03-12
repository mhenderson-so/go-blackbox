package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
	blackbox "github.ds.stackexchange.com/mhenderson/go-blackbox/cmd/go-blackbox"
	"github.ds.stackexchange.com/mhenderson/go-blackbox/version"
)

var ()

func main() {
	bb := cli.NewApp()
	bb.Name = "BlackBox"
	bb.Usage = "Safely store secrets in Git/Mercurial/Subversion"
	bb.Version = version.GetVersionInfo()

	bb.Commands = commands // From commands.go
	err := bb.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func setup(filename string) {
	err := blackbox.InitBlackbox(filename)
	if err != nil {
		log.Fatal(err)
	}
}

// stripGpgExtension strips the GPG extension from a filename only if it's present
func stripGpgExtension(filename string) string {
	len := len(filename)
	if filename[len-3:] == "gpg" {
		return filename[0 : len-4]
	}
	return filename
}
