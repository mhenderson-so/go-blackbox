package main

import (
	"github.com/urfave/cli"
)

var commands = []cli.Command{
	{
		Name:  "list",
		Usage: "List details of a blackboxed repository",
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
	{
		Name:  "edit",
		Usage: "Edit a blackboxed file",
		Subcommands: []cli.Command{
			{
				Name:  "start",
				Usage: "Decrypt a file so you can start editing it in another program",
				Action: func(c *cli.Context) error {
					return editStart(c.Args())
				},
			},
		},
	},
}
