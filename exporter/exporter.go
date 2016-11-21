package exporter

import (
	"database/sql"
	"strings"

	"github.com/dev-cloverlab/carpenter/dialect/mysql"
)

func Export(db *sql.DB, schema string, tableName string) (string, error) {
	cnk, err := mysql.GetChunk(db, tableName, nil)
	if err != nil {
		return "", err
	}
	csv := make([]string, 0, len(cnk.Seeds)+1)
	csv = append(csv, strings.Join(cnk.ColumnNames, ","))
	for _, seed := range cnk.Seeds {
		cols := make([]string, 0, len(cnk.ColumnNames))
		for i := range seed.ColumnData {
			cols = append(cols, seed.ToColumnValue(i))
		}
		csv = append(csv, strings.Join(cols, ","))
	}
	return strings.Join(csv, "\n"), nil
}
