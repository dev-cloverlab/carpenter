package command

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/codegangsta/cli"
	"github.com/dev-cloverlab/carpenter/dialect/mysql"
	"github.com/dev-cloverlab/carpenter/exporter"
)

func CmdExport(c *cli.Context) {
	// Write your code here
	dirPath := c.String("dir")
	if dirPath == "" {
		var err error
		dirPath, err = os.Getwd()
		if err != nil {
			panic(fmt.Errorf("err: os.Getwd failed for reason %s", err))
		}
	}

	re := c.String("regexp")
	tableNameRegexp := regexp.MustCompile(re)
	tables, err := mysql.GetTables(db, schema)
	if err != nil {
		panic(fmt.Errorf("err: mysql.GetTables failed for reason %s", err))
	}

	errs := []error{}
	errCh := make(chan error)
	succeedCh := make(chan string)
	doneCh := make(chan bool)
	go func() {
		for {
			select {
			case err := <-errCh:
				errs = append(errs, err)
			case tableName := <-succeedCh:
				if verbose {
					fmt.Println(tableName)
				}
			case <-doneCh:
				if len(errs) > 0 {
					msg := make([]string, 0, len(errs))
					for _, err := range errs {
						msg = append(msg, err.Error())
					}
					panic(fmt.Errorf("err: %s", strings.Join(msg, "\n")))
				}
				return
			}
		}
	}()

	wg := sync.WaitGroup{}
	for _, table := range tables {
		tableName := table.TableName
		if !tableNameRegexp.MatchString(tableName) {
			continue
		}
		wg.Add(1)
		go func(d *sql.DB, s, t string) {
			defer wg.Done()
			csv, err := exporter.Export(d, s, t)
			if err != nil {
				errCh <- err
				return
			}
			if err := ioutil.WriteFile(fmt.Sprintf("%s%s%s.csv", dirPath, string(os.PathSeparator), t), []byte(csv), os.ModePerm); err != nil {
				errCh <- err
				return
			}
			succeedCh <- t
		}(db, schema, tableName)
	}
	wg.Wait()
	doneCh <- true
}
