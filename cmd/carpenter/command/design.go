package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/codegangsta/cli"
	"github.com/dev-cloverlab/carpenter/designer"
	"github.com/dev-cloverlab/carpenter/dialect/mysql"
)

func CmdDesign(c *cli.Context) {
	// Write your code here
	dirPath := c.String("dir")
	if dirPath == "" {
		var err error
		dirPath, err = os.Getwd()
		if err != nil {
			panic(fmt.Errorf("err: os.Getwd failed for reason %s", err))
		}
	}
	pretty := c.Bool("pretty")
	separate := c.Bool("separate")

	tables, err := designer.Export(db, schema)
	if err != nil {
		panic(fmt.Errorf("err: designer.Export failed for reason %s", err))
	}

	if separate {
		for _, table := range tables {
			var j []byte
			var err error
			if pretty {
				j, err = json.MarshalIndent(mysql.Tables{table}, "", "\t")
			} else {
				j, err = json.Marshal(mysql.Tables{table})
			}
			if err != nil {
				panic(fmt.Errorf("err: json.MarshalIndent is fialed for reason %s", err))
			}
			if err := exportJson(dirPath, table.TableName, j); err != nil {
				panic(fmt.Errorf("err: exportJson is fialed for reason %s", err))
			}
		}
	} else {
		var j []byte
		var err error
		if pretty {
			j, err = json.MarshalIndent(tables, "", "\t")
		} else {
			j, err = json.Marshal(tables)
		}
		if err != nil {
			panic(fmt.Errorf("err: json.MarshalIndent is fialed for reason %s", err))
		}
		if err := exportJson(dirPath, "tables", j); err != nil {
			panic(fmt.Errorf("err: exportJson is fialed for reason %s", err))
		}
	}
}

func exportJson(dirPath, filename string, buf []byte) error {
	return ioutil.WriteFile(fmt.Sprintf("%s%s%s.json", dirPath, string(os.PathSeparator), filename), buf, os.ModePerm)
}
