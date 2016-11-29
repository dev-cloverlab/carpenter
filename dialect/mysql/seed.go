package mysql

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type (
	Chunk struct {
		TableName   string
		ColumnNames []string
		Seeds       Seeds
	}

	Seeds []Seed

	Seed struct {
		ColumnData []interface{}
	}
)

var (
	TimeFmt         string = "2006-01-02 15:04:05"
	defaultBulkSize int    = 2000
	trancateSQLFmt  string = `truncate table %s`
	insertSQLFmt    string = `insert into %s(%s)
values
%s`
	deleteSQLFmt string = `delete from %s where %s in (
%s
)`
	replaceSQLFmt string = `replace into %s(%s)
values
%s`
)

func (m *Chunk) GetColumnIndexBy(columnName string) (int, error) {
	for i, cName := range m.ColumnNames {
		if columnName == cName {
			return i, nil
		}
	}
	return 0, fmt.Errorf("err: Specified columnName `%s' is not found in this table %s", columnName, m.TableName)
}

func (m *Chunk) GetSeedGroupBy(columnName string) (map[interface{}]Seed, error) {
	colIdx, err := m.GetColumnIndexBy(columnName)
	if err != nil {
		return nil, err
	}
	cdMap := make(map[interface{}]Seed, len(m.Seeds))
	for _, seed := range m.Seeds {
		if len(seed.ColumnData) <= 0 {
			continue
		}
		cdMap[seed.ColumnData[colIdx]] = seed
	}
	return cdMap, nil
}

func (m *Chunk) GetFormatedTableName() string {
	return Quote(m.TableName)
}

func (m *Chunk) ToTrancateSQL() string {
	return fmt.Sprintf(trancateSQLFmt, m.GetFormatedTableName())
}

func (m *Chunk) ToInsertSQL() []string {
	columnStr := strings.Join(QuoteMulti(m.ColumnNames), ",")

	seedSize := defaultBulkSize
	if seedSize > len(m.Seeds) {
		seedSize = len(m.Seeds)
	}
	queries := make([]string, 0, int(len(m.Seeds)/defaultBulkSize)+1)
	seeds := Seeds{}
	for _, seed := range m.Seeds {
		seeds = append(seeds, seed)
		if len(seeds) >= defaultBulkSize {
			queries = append(queries, fmt.Sprintf(insertSQLFmt, m.GetFormatedTableName(), columnStr, strings.Join(seeds.ToValueSQL(), ",\n")))
			seeds = Seeds{}
		}
	}
	if len(seeds) > 0 {
		queries = append(queries, fmt.Sprintf(insertSQLFmt, m.GetFormatedTableName(), columnStr, strings.Join(seeds.ToValueSQL(), ",\n")))
	}
	return queries
}

func (m *Chunk) ToDeleteSQL(colIdx int) []string {
	columnStr := Quote(m.ColumnNames[colIdx])
	seedSize := defaultBulkSize
	if seedSize > len(m.Seeds) {
		seedSize = len(m.Seeds)
	}
	queries := make([]string, 0, int(len(m.Seeds)/defaultBulkSize)+1)
	seeds := Seeds{}
	for _, seed := range m.Seeds {
		seeds = append(seeds, seed)
		if len(seeds) >= defaultBulkSize {
			queries = append(queries, fmt.Sprintf(deleteSQLFmt, m.GetFormatedTableName(), columnStr, strings.Join(seeds.ToColumnValues(colIdx), ",\n")))
			seeds = Seeds{}
		}
	}
	if len(seeds) > 0 {
		queries = append(queries, fmt.Sprintf(deleteSQLFmt, m.GetFormatedTableName(), columnStr, strings.Join(seeds.ToColumnValues(colIdx), ",\n")))
	}
	return queries
}

func (m *Chunk) ToReplaceSQL() []string {
	columnStr := strings.Join(QuoteMulti(m.ColumnNames), ",")

	seedSize := defaultBulkSize
	if seedSize > len(m.Seeds) {
		seedSize = len(m.Seeds)
	}
	queries := make([]string, 0, int(len(m.Seeds)/defaultBulkSize)+1)
	seeds := Seeds{}
	for _, seed := range m.Seeds {
		seeds = append(seeds, seed)
		if len(seeds) >= defaultBulkSize {
			queries = append(queries, fmt.Sprintf(replaceSQLFmt, m.GetFormatedTableName(), columnStr, strings.Join(seeds.ToValueSQL(), ",")))
			seeds = Seeds{}
		}
	}
	if len(seeds) > 0 {
		queries = append(queries, fmt.Sprintf(replaceSQLFmt, m.GetFormatedTableName(), columnStr, strings.Join(seeds.ToValueSQL(), ",")))
	}
	return queries
}

func (m Seeds) ToValueSQL() []string {
	sqls := make([]string, 0, len(m))
	for _, seed := range m {
		sqls = append(sqls, fmt.Sprintf("(%s)", seed.ToValueSQL()))
	}
	return sqls
}

func (m Seed) ToValueSQL() string {
	str := make([]string, 0, len(m.ColumnData))
	for _, data := range m.ColumnData {
		str = append(str, toString(data))
	}
	return strings.Join(str, ",")
}

func (m Seeds) ToColumnValues(colIdx int) []string {
	values := make([]string, 0, len(m))
	for _, seed := range m {
		values = append(values, seed.ToColumnValue(colIdx))
	}
	return values
}

func (m Seed) ToColumnValue(colIdx int) string {
	return toString(m.ColumnData[colIdx])
}

func (m Seed) ValueEqual(seed Seed) bool {
	cnt := len(m.ColumnData)
	for i := 0; i < cnt; i++ {
		if toString(m.ColumnData[i]) != toString(seed.ColumnData[i]) {
			return false
		}
	}
	return true
}

func toString(data interface{}) (str string) {
	switch data.(type) {
	case nil:
		str = "null"
	case string:
		str = QuoteString(Unescape(data.(string)))
	case time.Time:
		str = QuoteString(data.(time.Time).Format(TimeFmt))
	default:
		str = fmt.Sprintf("%v", data)
	}
	return str
}

func GetChunk(db *sql.DB, table string, colName *string) (*Chunk, error) {
	cntCol := "*"
	if colName != nil {
		cntCol = Quote(*colName)
	}
	res, err := db.Exec(fmt.Sprintf("select count(%s) from %s", cntCol, Quote(table)))
	if err != nil {
		return nil, err
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(fmt.Sprintf("select * from %s", table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	cnk := &Chunk{
		TableName:   table,
		ColumnNames: columns,
		Seeds:       make(Seeds, 0, cnt),
	}
	colLen := len(columns)
	holderPtrs := make([]interface{}, colLen)
	for rows.Next() {
		holders := make([]interface{}, colLen)
		for i := range columns {
			holderPtrs[i] = &holders[i]
		}
		if err := rows.Scan(holderPtrs...); err != nil {
			return nil, err
		}
		for i := range columns {
			var v interface{}
			var err error
			if b, ok := holders[i].([]byte); ok {
				if v, err = strconv.ParseFloat(string(b), 64); err != nil {
					v = string(b)
				}
			} else {
				v = holders[i]
			}
			holders[i] = v
		}
		cnk.Seeds = append(cnk.Seeds, Seed{
			ColumnData: holders,
		})
	}
	return cnk, nil
}
