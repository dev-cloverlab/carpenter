package mysql

import (
	"database/sql"
	"fmt"
	"strings"
)

type (
	Partition struct {
		TableCatalog                string
		TableSchema                 string
		TableName                   string
		PartitionName               string
		SubpartitionName            JsonNullString
		PartitionOrdinalPosition    string
		SubpartitionOrdinalPosition JsonNullString
		PartitionMethod             string
		SubpartitionMethod          JsonNullString
		PartitionExpression         string
		SubpartitionExpression      JsonNullString
		PartitionDescription        JsonNullString
		TableRows                   int
		AvgRowLength                int
		DataLength                  int
		MaxDataLength               JsonNullInt64
		IndexLength                 int
		DataFree                    int
		Checksum                    JsonNullString
		PartitionComment            string
		Nodegroup                   string
		TablespaceName              JsonNullString
	}
	Partitions []*Partition
)

const (
	PartitionMethodLinearKey  = "LINEAR KEY"
	PartitionMethodLinearHash = "LINEAR HASH"
	PartitionMethodRange      = "RANGE COLUMNS"
)

func (m Partitions) ToSQL() string {
	if m == nil {
		return ""
	}
	// only for support "LINEAR KEY" or "LINEAR HASH" or "RANGE COLUMNS" partition method
	switch m[0].PartitionMethod {
	case PartitionMethodLinearKey:
		return fmt.Sprintf("partition by linear key (%s) partitions %d", m[0].PartitionExpression, len(m))
	case PartitionMethodLinearHash:
		return fmt.Sprintf("partition by linear hash (%s) partitions %d", m[0].PartitionExpression, len(m))
	case PartitionMethodRange:
		sqls := make([]string, 0, len(m))
		for _, partition := range m {
			sqls = append(sqls, fmt.Sprintf("\t\tpartition %s values less than (%s)", partition.PartitionName, partition.PartitionDescription.String))
		}
		return fmt.Sprintf("partition by range columns (%s) (\n%s\n\t)", m[0].PartitionExpression, strings.Join(sqls, ",\n"))
	default:
		return ""
	}
}

func GetPartitions(db *sql.DB, schema string, tableName string) (Partitions, error) {
	var rows *sql.Rows
	var err error

	selectCols := []string{
		"TABLE_CATALOG",
		"TABLE_NAME",
		"PARTITION_NAME",
		"SUBPARTITION_NAME",
		"PARTITION_ORDINAL_POSITION",
		"SUBPARTITION_ORDINAL_POSITION",
		"PARTITION_METHOD",
		"SUBPARTITION_METHOD",
		"PARTITION_EXPRESSION",
		"SUBPARTITION_EXPRESSION",
		"PARTITION_DESCRIPTION",
		"TABLE_ROWS",
		"AVG_ROW_LENGTH",
		"DATA_LENGTH",
		"MAX_DATA_LENGTH",
		"INDEX_LENGTH",
		"DATA_FREE",
		"CHECKSUM",
		"PARTITION_COMMENT",
		"NODEGROUP",
		"TABLESPACE_NAME",
	}
	query := fmt.Sprintf(`select %s from information_schema.partitions where TABLE_SCHEMA=%s and TABLE_NAME=%s`, strings.Join(selectCols, ","), QuoteString(schema), QuoteString(tableName))
	rows, err = db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("err: db.Query `%s' failed for reason %s", query, err)
	}
	defer rows.Close()

	partitions := Partitions{}
	for rows.Next() {
		partition := &Partition{}
		if err := rows.Scan(
			&partition.TableCatalog,
			&partition.TableName,
			&partition.PartitionName,
			&partition.SubpartitionName,
			&partition.PartitionOrdinalPosition,
			&partition.SubpartitionOrdinalPosition,
			&partition.PartitionMethod,
			&partition.SubpartitionMethod,
			&partition.PartitionExpression,
			&partition.SubpartitionExpression,
			&partition.PartitionDescription,
			&partition.TableRows,
			&partition.AvgRowLength,
			&partition.DataLength,
			&partition.MaxDataLength,
			&partition.IndexLength,
			&partition.DataFree,
			&partition.Checksum,
			&partition.PartitionComment,
			&partition.Nodegroup,
			&partition.TablespaceName,
		); err != nil {
			return nil, err
		}
		partitions = append(partitions, partition)
	}

	return partitions, nil
}
