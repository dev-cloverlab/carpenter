package exporter

import (
	"database/sql"
	"encoding/json"

	"github.com/dev-cloverlab/carpenter/dialect/mysql"
)

func Export(db *sql.DB, pretty bool, schema string, tableNames ...string) ([]byte, error) {
	tables, err := mysql.GetTables(db, schema, tableNames...)
	if err != nil {
		return nil, err
	}
	var buf []byte
	if pretty {
		buf, err = json.MarshalIndent(tables, "", "\t")
	} else {
		buf, err = json.Marshal(tables)
	}
	if err != nil {
		return nil, err
	}
	return buf, nil
}
