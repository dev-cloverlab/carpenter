package exporter

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

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
}

func TestMain(m *testing.M) {
	code := m.Run()
	db.Exec("drop table if exists `export_test`")
	os.Exit(code)
}

func TestExport(t *testing.T) {
	_, err := db.Exec("create table if not exists `export_test` (\n" +
		"	`id` int(11) unsigned not null auto_increment,\n" +
		"	`name` varchar(64) not null default'',\n" +
		"	`email` varchar(255) not null default'',\n" +
		"	`gender` tinyint(4) not null,\n" +
		"	`country_code` int(11) not null,\n" +
		"	`comment` text,\n" +
		"	`created_at` datetime not null,\n" +
		"	primary key (`id`),\n" +
		"	unique key `k1` (`email`),\n" +
		"	key `k2` (`name`),\n" +
		"	key `k3` (`gender`,`country_code`)\n" +
		") engine=InnoDB default charset=utf8",
	)
	if err != nil {
		t.Fatal(err)
	}

	a, err := Export(db, false, schema, "export_test")
	if err != nil {
		t.Fatal(err)
	}
	actual := []*mysql.Table{}

	decoder := json.NewDecoder(bytes.NewBuffer(a))
	decoder.UseNumber()
	if err := decoder.Decode(&actual); err != nil {
		t.Error(err)
	}

	e, err := ioutil.ReadFile("./_test/expected.json")
	if err != nil {
		t.Fatal(err)
	}
	expected := []*mysql.Table{}
	if err := json.Unmarshal(e, &expected); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("err: unexpected JSON returned.\nactual: %s\nexpected: %s", a, e)
	}
}
