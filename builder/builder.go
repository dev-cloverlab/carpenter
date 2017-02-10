package builder

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/dev-cloverlab/carpenter/dialect/mysql"
)

func Build(db *sql.DB, old, new *mysql.Table) (queries []string, err error) {
	if old == nil && new == nil {
		return queries, fmt.Errorf("err: Both pointer of the specified new and old is nil.")
	}
	if old != nil && new != nil && old.TableName != new.TableName {
		return queries, fmt.Errorf("err: Table name of the specified new and old is a difference")
	}
	if reflect.DeepEqual(old, new) {
		return queries, nil
	}
	if q := willCreate(old, new); len(q) > 0 {
		queries = append(queries, q)
	}
	if q := willDrop(old, new); len(q) > 0 {
		queries = append(queries, q)
	}
	if q := willAlterTableCharacterSet(old, new); len(q) > 0 {
		queries = append(queries, q)
	}
	if q := willAlterColumnCharacterSet(old, new); len(q) > 0 {
		queries = append(queries, q...)
	}
	if q := willAlter(old, new); len(q) > 0 {
		queries = append(queries, q)
	}
	return queries, nil
}

func willCreate(old, new *mysql.Table) string {
	if old == nil && new != nil {
		return new.ToCreateSQL()
	}
	return ""
}

func willDrop(old, new *mysql.Table) string {
	if old != nil && new == nil {
		return old.ToDropSQL()
	}
	return ""
}

func willAlterTableCharacterSet(old, new *mysql.Table) string {
	if old == nil || new == nil {
		return ""
	}

	alter := []string{}
	if old.GetCharset() != new.GetCharset() {
		alter = append(alter, new.ToConvertCharsetSQL())
	}
	old.TableCollation = new.TableCollation

	return new.ToAlterSQL(alter)
}

func willAlterColumnCharacterSet(old, new *mysql.Table) []string {
	newCols := new.Columns.GroupByColumnName()
	oldCols := old.Columns.GroupByColumnName()
	sqls := []string{}
	for _, colName := range new.Columns.GetSortedColumnNames() {
		if _, ok := oldCols[colName]; !ok {
			continue
		}
		newCol := newCols[colName]
		oldCol := oldCols[colName]
		if oldCol.CompareCharacterSet(newCol) {
			continue
		}
		oldCols[colName].CollationName = newCol.CollationName
		sqls = append(sqls, new.ToAlterSQL([]string{newCol.ToModifyCharsetSQL()}))
	}
	return sqls
}

func willAlter(old, new *mysql.Table) string {
	if old == nil || new == nil {
		return ""
	}
	if reflect.DeepEqual(old, new) {
		return ""
	}

	alter := []string{}
	alter = append(alter, willDropIndex(old, new)...)
	alter = append(alter, willDropColumn(old, new)...)
	alter = append(alter, willAddColumn(old, new)...)
	alter = append(alter, willAddIndex(old, new)...)
	alter = append(alter, willModifyColumn(old, new)...)
	return new.ToAlterSQL(alter)
}

func willAddColumn(old, new *mysql.Table) []string {
	cols := mysql.Columns{}
	for _, column := range new.Columns {
		if old.Columns.Contains(column) {
			continue
		}
		cols = append(cols, column)
	}
	return cols.ToAddSQL(new.Columns)
}

func willDropColumn(old, new *mysql.Table) []string {
	cols := mysql.Columns{}
	for _, column := range old.Columns {
		if new.Columns.Contains(column) {
			continue
		}
		cols = append(cols, column)
	}
	return cols.ToDropSQL()
}

func willModifyColumn(old, new *mysql.Table) []string {
	newCols := new.Columns.GroupByColumnName()
	oldCols := old.Columns.GroupByColumnName()
	sqls := []string{}
	for _, colName := range new.Columns.GetSortedColumnNames() {
		if _, ok := oldCols[colName]; !ok {
			continue
		}
		newCol := newCols[colName]
		oldCol := oldCols[colName]
		oldTableSchema := oldCol.TableSchema
		oldColumnKey := oldCol.ColumnKey
		oldPrivileges := oldCol.Privileges
		oldCol.TableSchema = newCol.TableSchema
		oldCol.ColumnKey = newCol.ColumnKey
		oldCol.Privileges = newCol.Privileges
		if !reflect.DeepEqual(oldCol, newCol) {
			sql := newCol.ToModifySQL()
			sql = fmt.Sprintf("%s %s", sql, newCol.AppendPos(new.Columns))
			sqls = append(sqls, sql)
		}
		oldCol.TableSchema = oldTableSchema
		oldCol.ColumnKey = oldColumnKey
		oldCol.Privileges = oldPrivileges
	}
	return sqls
}

func willAddIndex(old, new *mysql.Table) []string {
	newIndicesMap := new.Indices.GroupByKeyName()
	oldIndicesMap := old.Indices.GroupByKeyName()
	sqls := []string{}
	for _, keyName := range new.Indices.GetSortedKeys() {
		if _, ok := oldIndicesMap[keyName]; !ok {
			sqls = append(sqls, newIndicesMap[keyName].ToAddSQL()...)
			continue
		}
		newIndices := newIndicesMap[keyName]
		oldIndices := oldIndicesMap[keyName]
		newIndices.ResetCardinality()
		oldIndices.ResetCardinality()
		if reflect.DeepEqual(oldIndices, newIndices) {
			continue
		}
		sqls = append(sqls, oldIndices.ToDropSQL()...)
		sqls = append(sqls, newIndices.ToAddSQL()...)
	}
	return sqls
}

func willDropIndex(old, new *mysql.Table) []string {
	newIndicesMap := new.Indices.GroupByKeyName()
	oldIndicesMap := old.Indices.GroupByKeyName()
	sqls := []string{}
	for _, keyName := range old.Indices.GetSortedKeys() {
		if _, ok := newIndicesMap[keyName]; ok {
			continue
		}
		sqls = append(sqls, oldIndicesMap[keyName].ToDropSQL()...)
	}
	return sqls
}
