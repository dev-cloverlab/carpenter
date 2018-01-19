package command

import (
	"database/sql"
	"time"

	"fmt"

	"github.com/codegangsta/cli"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var schema string
var verbose bool
var dryrun bool
var maxIdleConns int
var maxOpenConns int

func Before(c *cli.Context) error {
	verbose = c.GlobalBool("verbose")
	dryrun = c.GlobalBool("dry-run")
	schema = c.GlobalString("schema")
	maxIdleConns = c.GlobalInt("max-idle-conns")
	maxOpenConns = c.GlobalInt("max-open-conns")

	if len(schema) <= 0 {
		return fmt.Errorf("err: Specify required `--schema' option")
	}
	datasource := c.GlobalString("data-source")
	if len(datasource) <= 0 {
		return fmt.Errorf("err: Specify required `--data-soruce' option")
	}
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%s/%s?charset=utf8", datasource, schema))
	if err != nil {
		return fmt.Errorf("err: db.Open is failed for reason %v", err)
	}
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(time.Minute)
	return nil
}
