package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/codegangsta/cli"
	"github.com/dev-cloverlab/carpenter/builder"
	"github.com/dev-cloverlab/carpenter/dialect/mysql"
)

func CmdBuild(c *cli.Context) {
	// Write your code here
	dirPath := c.String("dir")
	withDrop := c.Bool("with-drop")
	queries, errs := makeBuildQueries(dirPath, withDrop)
	if len(errs) > 0 {
		panic(fmt.Errorf("err: makeQueries failed for reason\n%s", strings.Join(getErrorMessages(errs), "\n")))
	}
	if err := execute(queries); err != nil {
		panic(fmt.Errorf("err: execute failed for reason %s", err))
	}
}

func makeBuildQueries(path string, withDrop bool) (queries []string, errs []error) {
	files, err := walk(path, ".json")
	if err != nil {
		return nil, []error{err}
	}
	new := mysql.Tables{}
	for _, file := range files {
		tables, err := parseJSON(file[0])
		if err != nil {
			return nil, []error{err}
		}
		new = append(new, tables...)
	}
	old, err := mysql.GetTables(db, schema)
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

	newMap := new.GroupByTableName()
	oldMap := old.GroupByTableName()
	tableNames := getTableNames(newMap, oldMap)
	wg := &sync.WaitGroup{}
	for _, tableName := range tableNames {
		oTbl, ok := oldMap[tableName]
		if !ok {
			oTbl = nil
		}
		nTbl, ok := newMap[tableName]
		if !ok {
			nTbl = nil
		}
		wg.Add(1)
		go func(o, n *mysql.Table) {
			defer wg.Done()
			queries, err := builder.Build(db, o, n, withDrop)
			if err != nil {
				errCh <- err
				return
			}
			sqlCh <- queries
		}(oTbl, nTbl)
	}
	wg.Wait()

	doneCh <- true

	return queries, errs
}

func parseJSON(filename string) (mysql.Tables, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	tables := mysql.Tables{}
	if err := json.Unmarshal(buf, &tables); err != nil {
		return nil, err
	}
	return tables, nil
}

func getTableNames(new, old map[string]*mysql.Table) []string {
	tableNames := map[string]struct{}{}
	for tableName := range new {
		tableNames[tableName] = struct{}{}
	}
	for tableName := range old {
		tableNames[tableName] = struct{}{}
	}
	ret := make([]string, 0, len(tableNames))
	for name := range tableNames {
		ret = append(ret, name)
	}
	return ret
}
