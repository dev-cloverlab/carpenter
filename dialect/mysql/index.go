package mysql

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
)

type (
	IndexColumn struct {
		Table        string
		NonUniue     int8
		KeyName      string
		SeqInIndex   int32
		ColumnName   string
		Collation    string
		SubPart      JsonNullString
		Packed       JsonNullString
		Null         string
		IndexType    string
		Comment      string
		IndexComment string
	}
	Index   []IndexColumn
	Indices []Index
)

func (m Index) IsPrimaryKey() bool {
	return m[0].KeyName == "PRIMARY"
}

func (m Index) IsUniqueKey() bool {
	return m[0].NonUniue != 1
}

func (m Index) GetKeyName() string {
	return m[0].KeyName
}

func (m Index) ColumnNames() []string {
	names := make([]string, 0, len(m))
	for _, info := range m {
		names = append(names, Quote(info.ColumnName))
	}
	return names
}

func (m Index) KeyNamesWithSubPart() []string {
	names := make([]string, 0, len(m))
	for _, info := range m {
		if info.SubPart.Valid {
			names = append(names, fmt.Sprintf("%s(%s)", Quote(info.ColumnName), info.SubPart.String))
		} else {
			names = append(names, Quote(info.ColumnName))
		}
	}
	return names
}

func (m Indices) ToSQL() []string {
	indices := m.getSortedIndices(m.GetSortedKeys())
	indexSQLs := make([]string, 0, len(m))
	for _, index := range indices {
		sql := ""
		comment := index[0].Comment
		switch {
		case index.IsPrimaryKey():
			sql = fmt.Sprintf("primary key (%s)", strings.Join(index.ColumnNames(), ","))
		case !index.IsPrimaryKey() && index.IsUniqueKey():
			sql = fmt.Sprintf("unique key %s (%s)", Quote(index.GetKeyName()), strings.Join(index.ColumnNames(), ","))
		default:
			sql = fmt.Sprintf("key %s (%s)", Quote(index.GetKeyName()), strings.Join(index.KeyNamesWithSubPart(), ","))
		}
		if comment != "" {
			sql = fmt.Sprintf("%s comment %s", sql, QuoteString(comment))
		}
		indexSQLs = append(indexSQLs, sql)
	}
	return indexSQLs
}

func (m Indices) ToAddSQL() []string {
	sqls := m.ToSQL()
	for i, sql := range sqls {
		sqls[i] = fmt.Sprintf("add %s", sql)
	}
	return sqls
}

func (m Indices) ToDropSQL() []string {
	idxMap := map[string]struct{}{}
	for _, index := range m {
		idxMap[index.GetKeyName()] = struct{}{}
	}
	sqls := make([]string, 0, len(idxMap))
	for keyName := range idxMap {
		sqls = append(sqls, fmt.Sprintf("drop key %s", Quote(keyName)))
	}
	return sqls
}

func (m Indices) GroupByKeyName() map[string]Indices {
	nameMap := make(map[string]Indices, len(m))
	for _, index := range m {
		if _, ok := nameMap[index.GetKeyName()]; !ok {
			nameMap[index.GetKeyName()] = Indices{}
		}
		nameMap[index.GetKeyName()] = append(nameMap[index.GetKeyName()], index)
	}
	return nameMap
}

func (m Indices) GetSortedKeys() []string {
	keys := make([]string, 0, len(m))
	for _, index := range m {
		keys = append(keys, index.GetKeyName())
	}
	sort.Strings(keys)
	return keys
}

func (m Indices) getSortedIndices(keys []string) Indices {
	idxMap := map[string]Index{}
	for _, index := range m {
		idxMap[index.GetKeyName()] = index
	}
	indices := make(Indices, 0, len(idxMap))
	for _, key := range keys {
		index := idxMap[key]
		if index.IsPrimaryKey() {
			indices = append(indices, index)
		}
	}
	for _, key := range keys {
		index := idxMap[key]
		if !index.IsPrimaryKey() && index.IsUniqueKey() {
			indices = append(indices, index)
		}
	}
	for _, key := range keys {
		index := idxMap[key]
		if !index.IsPrimaryKey() && !index.IsUniqueKey() {
			indices = append(indices, index)
		}
	}
	return indices
}

func GetIndices(db *sql.DB, table string) (Indices, error) {
	query := fmt.Sprintf("show index from `%s`", table)
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("err: db.Query faild `%s' for reason %s", query, err)
	}
	defer rows.Close()

	dummy := JsonNullInt64{}
	idxMap := map[string]Index{}
	for rows.Next() {
		idxCol := IndexColumn{}
		if err := rows.Scan(
			&idxCol.Table,
			&idxCol.NonUniue,
			&idxCol.KeyName,
			&idxCol.SeqInIndex,
			&idxCol.ColumnName,
			&idxCol.Collation,
			&dummy,
			&idxCol.SubPart,
			&idxCol.Packed,
			&idxCol.Null,
			&idxCol.IndexType,
			&idxCol.IndexComment,
			&idxCol.Comment,
		); err != nil {
			return nil, err
		}
		if _, ok := idxMap[idxCol.KeyName]; !ok {
			idxMap[idxCol.KeyName] = Index{}
		}
		idxMap[idxCol.KeyName] = append(idxMap[idxCol.KeyName], idxCol)
	}
	indices := make(Indices, 0, len(idxMap))
	for _, index := range idxMap {
		indices = append(indices, index)
	}
	return indices.getSortedIndices(indices.GetSortedKeys()), nil
}
