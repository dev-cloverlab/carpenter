package command

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/codegangsta/cli"
	"github.com/dev-cloverlab/carpenter/dialect/mysql"
	"github.com/dev-cloverlab/carpenter/seeder"
)

func CmdSeed(c *cli.Context) {
	// Write your code here
	dirPath := c.String("dir")
	queries, errs := makeSeedQueries(dirPath, nil)
	if len(errs) > 0 {
		panic(fmt.Errorf("err: makeSeedQueries failed for reason\n%s", strings.Join(getErrorMessages(errs), "\n")))
	}
	if err := execute(queries); err != nil {
		panic(fmt.Errorf("err: execute failed for reason %s", err))
	}
}

func makeSeedQueries(path string, colName *string) (queries []string, errs []error) {
	files, err := walk(path, ".csv")
	if err != nil {
		return nil, []error{err}
	}

	errCh := make(chan error)
	sqlCh := make(chan []string)
	doneCh := make(chan bool)
	go func() {
		for {
			select {
			case err := <-errCh:
				errs = append(errs, err)
			case query := <-sqlCh:
				queries = append(queries, query...)
			case <-doneCh:
				return
			}
		}
	}()

	wg := &sync.WaitGroup{}
	for tableName, file := range files {
		wg.Add(1)
		go func(t string, fs []string, c *string) {
			defer wg.Done()

			var colNames []string
			var err error
			seeds := mysql.Seeds{}
			for _, f := range fs {
				var s mysql.Seeds
				colNames, s, err = parseCSV(t, f)
				if err != nil {
					errCh <- fmt.Errorf("err: parseCSV %s failed for reason %s", t, err)
					return
				}
				seeds = append(seeds, s...)
			}
			new := makeChunk(t, colNames, seeds)
			old, err := mysql.GetChunk(db, t, c)
			if err != nil {
				errCh <- fmt.Errorf("err: mysql.GetChunk %s failed for reason %s", t, err)
				return
			}
			queries, err := seeder.Seed(db, old, new, c)
			if err != nil {
				errCh <- fmt.Errorf("err: seeder.Seed %s failed for reason %s", t, err)
			}
			sqlCh <- queries
		}(tableName, file, colName)
	}
	wg.Wait()

	doneCh <- true

	return queries, errs
}

func parseCSV(tableName, filename string) (columnNames []string, seeds mysql.Seeds, err error) {
	fp, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer fp.Close()

	reader := csv.NewReader(fp)
	reader.LazyQuotes = true

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		if len(columnNames) <= 0 {
			columnNames = record
		} else {
			columnData := make([]interface{}, 0, len(record))
			for _, r := range record {
				var v interface{}
				var err error
				if r == "NULL" || r == "null" || r == "Null" {
					v = nil
				} else if v, err = strconv.ParseFloat(r, 64); err != nil {
					v = r
				} else {
					if len(string(r)) > 1 && string(string(r)[0]) == "0" {
						v = string(r)
					}
				}
				columnData = append(columnData, v)
			}
			seeds = append(seeds, mysql.Seed{
				ColumnData: columnData,
			})
		}
	}
	return columnNames, seeds, nil
}

func makeChunk(tableName string, columnNames []string, seeds mysql.Seeds) *mysql.Chunk {
	return &mysql.Chunk{
		TableName:   tableName,
		ColumnNames: columnNames,
		Seeds:       seeds,
	}
}
