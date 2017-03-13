package mysql

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
)

type (
	Table struct {
		TableCatalog   string
		TableSchema    string
		TableName      string
		TableType      string
		Engine         string
		Version        int
		RowFormat      string
		TableRows      int
		AvgRowLength   int
		DataLength     int
		MaxDataLength  int
		IndexLength    int
		DataFree       int
		AutoIncrement  JsonNullInt64
		TableCollation string
		CheckSum       JsonNullString
		CreateOptions  string
		TableComment   string
		Columns        Columns
		Indices        Indices
		Partitions     Partitions
	}
	Tables []*Table
)

var (
	createSQLFmt string = `create table if not exists %s (
	%s
) engine=%s default charset=%s %s`
	dropSQLFmt  string = `drop table if exists %s`
	alterSQLFmt string = `alter table %s
	%s`
)

func (m *Table) IsPartitioned() bool {
	return m.CreateOptions == "partitioned"
}

func (m *Table) GetFormatedTableName() string {
	return Quote(m.TableName)
}

func (m *Table) GetCharset() string {
	seg := strings.Split(m.TableCollation, "_")
	return seg[0]
}

func (m *Table) ToAlterSQL(sqls []string) string {
	if len(sqls) <= 0 {
		return ""
	}
	return fmt.Sprintf(alterSQLFmt, m.GetFormatedTableName(), strings.Join(sqls, ",\n	"))
}

func (m Tables) GetFormatedTableNames() []string {
	names := make([]string, 0, len(m))
	for _, table := range m {
		names = append(names, table.GetFormatedTableName())
	}
	return names
}

func (m *Table) ToCreateSQL() string {
	columnSQLs := m.Columns.ToSQL()
	indexSQLs := m.Indices.ToSQL()
	partitionSQL := m.Partitions.ToSQL()
	sqls := make([]string, 0, len(columnSQLs)+len(indexSQLs))
	sqls = append(columnSQLs, indexSQLs...)
	return fmt.Sprintf(createSQLFmt, m.GetFormatedTableName(), strings.Join(sqls, ",\n	"), m.Engine, m.GetCharset(), partitionSQL)
}

func (m *Table) ToDropSQL() string {
	return fmt.Sprintf(dropSQLFmt, m.GetFormatedTableName())
}

func (m *Table) ToConvertCharsetSQL() string {
	return fmt.Sprintf("convert to character set %s", m.GetCharset())
}

func (m Tables) Contains(t *Table) bool {
	for _, v := range m {
		if v.TableName == t.TableName {
			return true
		}
	}
	return false
}

func (m Tables) GroupByTableName() map[string]*Table {
	nameMap := make(map[string]*Table, len(m))
	for _, table := range m {
		nameMap[table.TableName] = table
	}
	return nameMap
}

func (m Tables) GetSortedTableNames() []string {
	names := make([]string, 0, len(m))
	for _, table := range m {
		names = append(names, table.TableName)
	}
	sort.Strings(names)
	return names
}

func GetTables(db *sql.DB, schema string, tableNames ...string) (Tables, error) {
	var rows *sql.Rows
	var err error

	selectCols := []string{
		"TABLE_CATALOG",
		"TABLE_SCHEMA",
		"TABLE_NAME",
		"TABLE_TYPE",
		"ENGINE",
		"VERSION",
		"ROW_FORMAT",
		"TABLE_ROWS",
		"AVG_ROW_LENGTH",
		"DATA_LENGTH",
		"MAX_DATA_LENGTH",
		"INDEX_LENGTH",
		"DATA_FREE",
		"AUTO_INCREMENT",
		"TABLE_COLLATION",
		"CHECKSUM",
		"CREATE_OPTIONS",
		"TABLE_COMMENT",
	}
	query := fmt.Sprintf(`select %s from information_schema.tables where TABLE_SCHEMA=%s`, strings.Join(selectCols, ","), QuoteString(schema))
	if len(tableNames) > 0 {
		tn := make([]string, 0, len(tableNames))
		for _, t := range tableNames {
			tn = append(tn, QuoteString(t))
		}
		query = fmt.Sprintf(`select %s from information_schema.tables where TABLE_SCHEMA=%s and TABLE_NAME in (%s)`, strings.Join(selectCols, ","), QuoteString(schema), strings.Join(tn, ","))
	}
	rows, err = db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("err: db.Query `%s' failed for reason %s", query, err)
	}
	defer rows.Close()

	tables := []*Table{}
	for rows.Next() {
		table := &Table{}
		if err := rows.Scan(
			&table.TableCatalog,
			&table.TableSchema,
			&table.TableName,
			&table.TableType,
			&table.Engine,
			&table.Version,
			&table.RowFormat,
			&table.TableRows,
			&table.AvgRowLength,
			&table.DataLength,
			&table.MaxDataLength,
			&table.IndexLength,
			&table.DataFree,
			&table.AutoIncrement,
			&table.TableCollation,
			&table.CheckSum,
			&table.CreateOptions,
			&table.TableComment,
		); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	columns, err := GetColumns(db, schema)
	if err != nil {
		return nil, err
	}
	for i, table := range tables {
		indices, err := GetIndices(db, table.TableName)
		if err != nil {
			return nil, err
		}
		c := []*Column{}
		for _, v := range columns {
			if table.TableName != v.TableName {
				continue
			}
			c = append(c, v)
		}
		if table.IsPartitioned() {
			partitions, err := GetPartitions(db, table.TableSchema, table.TableName)
			if err != nil {
				return nil, err
			}
			tables[i].Partitions = partitions
		}
		tables[i].Columns = c
		tables[i].Indices = indices
	}
	return tables, nil
}
