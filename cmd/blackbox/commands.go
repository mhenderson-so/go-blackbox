package main

import (
	"github.com/urfave/cli"
)

var commands = []cli.Command{
	{
		Name:  "list",
		Usage: "list some details of a blackboxed repository",
		Subcommands: []cli.Command{
			{
				Name:  "admins",
				Usage: "List the blackbox admins of this repository",
				Action: func(c *cli.Context) error {
					return listAdmins()
				},
			},
			{
				Name:  "files",
				Usage: "List the blackboxed files in this repository",
				Action: func(c *cli.Context) error {
					return listFiles()
				},
			},
		},
	},
}
