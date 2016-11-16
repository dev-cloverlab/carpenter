package seeder

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/dev-cloverlab/carpenter/dialect/mysql"
)

var (
	defaultComparisonColumnName string = "id"
)

func Seed(db *sql.DB, old, new *mysql.Chunk, compColName *string) (queries []string, err error) {
	if old == nil && new == nil {
		return []string{}, fmt.Errorf("err: Both pointer of the specified new and old is nil.")
	}
	if reflect.DeepEqual(old, new) {
		return queries, nil
	}
	ccName := defaultComparisonColumnName
	if compColName != nil {
		ccName = *compColName
	}
	if q := willTruncate(old, new); len(q) > 0 {
		queries = append(queries, q)
	} else {
		if q, err := willDelete(old, new, ccName); err != nil {
			return []string{}, err
		} else if len(q) > 0 {
			queries = append(queries, q...)
		}
	}
	if q, err := willReplace(old, new, ccName); err != nil {
		return []string{}, err
	} else if len(q) > 0 {
		queries = append(queries, q...)
	}
	if q, err := willInsert(old, new, ccName); err != nil {
		return []string{}, err
	} else if len(q) > 0 {
		queries = append(queries, q...)
	}
	return queries, nil
}

func willTruncate(old, new *mysql.Chunk) string {
	if len(old.Seeds) != 0 && len(new.Seeds) <= 0 {
		return old.ToTrancateSQL()
	}
	return ""
}

func willInsert(old, new *mysql.Chunk, compColName string) ([]string, error) {
	if new == nil {
		return []string{}, nil
	}
	if old == nil && new != nil {
		return new.ToInsertSQL(), nil
	}
	oldMap, err := old.GetSeedGroupBy(compColName)
	if err != nil {
		return nil, err
	}
	newMap, err := new.GetSeedGroupBy(compColName)
	if err != nil {
		return nil, err
	}
	seeds := mysql.Seeds{}
	for k, seed := range newMap {
		if _, ok := oldMap[k]; ok {
			continue
		}
		seeds = append(seeds, seed)
	}
	cnk := mysql.Chunk{
		TableName:   new.TableName,
		ColumnNames: new.ColumnNames,
		Seeds:       seeds,
	}
	return cnk.ToInsertSQL(), nil
}

func willDelete(old, new *mysql.Chunk, compColName string) ([]string, error) {
	if old == nil {
		return []string{}, nil
	}
	if old != nil && new == nil {
		return []string{old.ToTrancateSQL()}, nil
	}
	oldMap, err := old.GetSeedGroupBy(compColName)
	if err != nil {
		return nil, err
	}
	newMap, err := new.GetSeedGroupBy(compColName)
	if err != nil {
		return nil, err
	}
	seeds := mysql.Seeds{}
	for k, seed := range oldMap {
		if _, ok := newMap[k]; ok {
			continue
		}
		seeds = append(seeds, seed)
	}
	cnk := mysql.Chunk{
		TableName:   new.TableName,
		ColumnNames: new.ColumnNames,
		Seeds:       seeds,
	}
	colIdx, err := cnk.GetColumnIndexBy(compColName)
	if err != nil {
		return nil, err
	}
	return cnk.ToDeleteSQL(colIdx), nil
}

func willReplace(old, new *mysql.Chunk, compColName string) ([]string, error) {
	if old == nil || new == nil {
		return []string{}, nil
	}
	oldMap, err := old.GetSeedGroupBy(compColName)
	if err != nil {
		return nil, err
	}
	newMap, err := new.GetSeedGroupBy(compColName)
	if err != nil {
		return nil, err
	}
	seeds := mysql.Seeds{}
	for k, seed := range newMap {
		if _, ok := oldMap[k]; !ok {
			continue
		}
		if newMap[k].ValueEqual(oldMap[k]) {
			continue
		}
		seeds = append(seeds, seed)
	}
	cnk := mysql.Chunk{
		TableName:   new.TableName,
		ColumnNames: new.ColumnNames,
		Seeds:       seeds,
	}
	return cnk.ToReplaceSQL(), nil
}
