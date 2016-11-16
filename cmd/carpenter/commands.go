package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/dev-cloverlab/carpenter/cmd/carpenter/command"
)

var GlobalFlags = []cli.Flag{
	cli.BoolFlag{
		Name:   "verbose, vv",
		Hidden: false,
		Usage:  "show verbose output (default off)",
	},
	cli.BoolFlag{
		Name:   "dry-run",
		Hidden: false,
		Usage:  "execute as dry-run mode (default off)",
	},
	cli.StringFlag{
		Name:   "schema, s",
		Usage:  "database name (requires)",
		Hidden: false,
	},
	cli.StringFlag{
		Name:   "data-source, d",
		Usage:  "data source name like '[username[:password]@][tcp[(address:port)]]' (requires)",
		Hidden: false,
	},
}

var Commands = []cli.Command{
	{
		Name:   "export",
		Usage:  "Export table structure as JSON string",
		Before: command.Before,
		Action: command.CmdExport,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:   "pretty, p",
				Usage:  "show pretty output (default off)",
				Hidden: false,
			},
		},
	},
	{
		Name:   "build",
		Usage:  "Build(Migrate) table from specified JSON string",
		Before: command.Before,
		Action: command.CmdBuild,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "dir, d",
				Usage:  "path to JSON file directory (requires)",
				Hidden: false,
			},
		},
	},
	{
		Name:   "seed",
		Usage:  "Seed CSV to table",
		Before: command.Before,
		Action: command.CmdSeed,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "dir, d",
				Usage:  "path to CSV file directory (requires)",
				Hidden: false,
			},
		},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
