package command

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/dev-cloverlab/carpenter/exporter"
)

func CmdExport(c *cli.Context) {
	// Write your code here
	pretty := c.Bool("pretty")
	json, err := exporter.Export(db, pretty, schema)
	if err != nil {
		panic(fmt.Errorf("err: Export failed for reason %s", err))
	}
	fmt.Printf("%s", json)
}
