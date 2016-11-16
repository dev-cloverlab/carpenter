package seeder

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/dev-cloverlab/carpenter/builder"
	"github.com/dev-cloverlab/carpenter/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db     *sql.DB
	schema = "test"
)

func init() {
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("root@/%s", schema))
	if err != nil {
		panic(err)
	}
	new, err := getTables("./_test/table1.json")
	if err != nil {
		panic(err)
	}
	createSQL, err := builder.Build(db, nil, new[0])
	if err != nil {
		panic(err)
	}
	for _, sql := range createSQL {
		if _, err := db.Exec(sql); err != nil {
			panic(err)
		}
	}
}

func TestMain(m *testing.M) {
	code := m.Run()
	db.Exec("drop table if exists `seed_test`")
	os.Exit(code)
}

func TestInsert(t *testing.T) {
	colName := "int"
	oldChunk, err := mysql.GetChunk(db, "seed_test", &colName)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().Format(mysql.TimeFmt)
	newChunk := makeChunk("seed_test", oldChunk.ColumnNames, mysql.Seeds{
		makeSeed([]interface{}{float64(10), "stringA", now, nil}),
		makeSeed([]interface{}{float64(20), "stringB", now, nil}),
	})

	expected := []string{
		"insert into `seed_test`(`int`,`string`,`time`,`null`)\n" +
			"values\n" +
			fmt.Sprintf("(10,\"stringA\",\"%v\",null),\n", now) +
			fmt.Sprintf("(20,\"stringB\",\"%v\",null)", now),
	}
	compString := "int"
	actual, err := Seed(db, oldChunk, newChunk, &compString)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("err: create: unexpected SQL returned.\nactual:\n%s\nexpected:\n%s\n", actual, expected)
	}
	for _, sql := range actual {
		if _, err := db.Exec(sql); err != nil {
			t.Fatal(err)
		}
	}
}

func TestReplace(t *testing.T) {
	colName := "int"
	oldChunk, err := mysql.GetChunk(db, "seed_test", &colName)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().Format(mysql.TimeFmt)
	newChunk := makeChunk("seed_test", oldChunk.ColumnNames, mysql.Seeds{
		makeSeed([]interface{}{float64(10), "stringC", now, nil}),
		makeSeed([]interface{}{float64(20), "stringB", now, nil}),
	})

	expected := []string{
		"replace into `seed_test`(`int`,`string`,`time`,`null`)\n" +
			"values\n" +
			fmt.Sprintf("(10,\"stringC\",\"%v\",null)", now),
	}
	compString := "int"
	actual, err := Seed(db, oldChunk, newChunk, &compString)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("err: create: unexpected SQL returned.\nactual:\n%s\nexpected:\n%s\n", actual, expected)
	}
	for _, sql := range actual {
		if _, err := db.Exec(sql); err != nil {
			t.Fatal(err)
		}
	}
}

func TestDelete(t *testing.T) {
	colName := "int"
	oldChunk, err := mysql.GetChunk(db, "seed_test", &colName)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().Format(mysql.TimeFmt)
	newChunk := makeChunk("seed_test", oldChunk.ColumnNames, mysql.Seeds{
		makeSeed([]interface{}{float64(20), "stringB", now, nil}),
	})

	expected := []string{
		"delete from `seed_test` where `int` in (\n" +
			"10\n" +
			")",
	}
	compString := "int"
	actual, err := Seed(db, oldChunk, newChunk, &compString)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("err: create: unexpected SQL returned.\nactual:\n%s\nexpected:\n%s\n", actual, expected)
	}
	for _, sql := range actual {
		if _, err := db.Exec(sql); err != nil {
			t.Fatal(err)
		}
	}
}

func TestTruncate(t *testing.T) {
	colName := "int"
	oldChunk, err := mysql.GetChunk(db, "seed_test", &colName)
	if err != nil {
		t.Fatal(err)
	}
	newChunk := makeChunk("seed_test", oldChunk.ColumnNames, mysql.Seeds{})

	expected := []string{
		"truncate table `seed_test`",
	}
	compString := "int"
	actual, err := Seed(db, oldChunk, newChunk, &compString)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("err: create: unexpected SQL returned.\nactual:\n%s\nexpected:\n%s\n", actual, expected)
	}
	for _, sql := range actual {
		if _, err := db.Exec(sql); err != nil {
			t.Fatal(err)
		}
	}
}

func makeSeed(columnData []interface{}) mysql.Seed {
	return mysql.Seed{
		ColumnData: columnData,
	}
}

func makeChunk(tableName string, columnNames []string, seeds mysql.Seeds) *mysql.Chunk {
	return &mysql.Chunk{
		TableName:   tableName,
		ColumnNames: columnNames,
		Seeds:       seeds,
	}
}

func getTables(filename string) (mysql.Tables, error) {
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
