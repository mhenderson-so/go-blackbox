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
		Name:  "initialize",
		Usage: "Initialize a new repository for use with blackbox",
		Action: func(c *cli.Context) error {
			return initialize()
		},
	},
	{
		Name:   "initialise", //Not doing this as an alias as I want it to be hidden
		Hidden: true,
		Action: func(c *cli.Context) error {
			return initialize()
		},
	},
	{
		Name:  "cat",
		Usage: "Cat a blackboxed file",
		Action: func(c *cli.Context) error {
			return cat(c.Args())
		},
	},
	{
		Name:  "addadmin",
		Usage: "Add an administrator to this repository",
		Action: func(c *cli.Context) error {
			return addadmin(c.Args())
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
			{
				Name:  "end",
				Usage: "Encrypt a file once you have finished editing it in another program",
				Action: func(c *cli.Context) error {
					return editEnd(c.Args())
				},
			},
		},
	},
	{
		Name:  "register",
		Usage: "Take a previously unencrypted file and enrol it in blackbox",
		Action: func(c *cli.Context) error {
			return register(c.Args())
		},
	},

	{
		Name:  "cleanup",
		Usage: "Cleans up the public keychain and removes orphaned entries",
		Action: func(c *cli.Context) error {
			return admincleanup()
		},
	},
}
