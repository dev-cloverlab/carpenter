package command

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/dev-cloverlab/carpenter/designer"
)

func CmdDesign(c *cli.Context) {
	// Write your code here
	pretty := c.Bool("pretty")
	json, err := designer.Export(db, pretty, schema)
	if err != nil {
		panic(fmt.Errorf("err: designer.Export failed for reason %s", err))
	}
	fmt.Printf("%s", json)
}
