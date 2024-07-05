package repository

import (
	"context"
	"fmt"
	"github.com/eugene-vodyanko/orgonaut/internal/model"
	"github.com/eugene-vodyanko/orgonaut/pkg/oracle"
	ora "github.com/sijms/go-ora/v2"
	"log/slog"
	"time"
)

type Repository struct {
	*oracle.Oracle
	schema string
}

func NewRepository(schema string, db *oracle.Oracle) *Repository {
	return &Repository{
		schema: schema,
		Oracle: db,
	}
}

// GetRecords receives the changed rows in the database.
//
// The implementation is based on the underlying pl/sql package (see org$gate_api).
// The data is grouped into batches to increase throughput.
// The data is encoded in XML format using the high-performance Oracle dbms_xmlgen core package (written in C).
// For efficient transmission over the network, data is also compressed using the gzip algorithm.
func (r *Repository) GetRecords(ctx context.Context, task *model.Task) ([]*model.Record, error) {
	start := time.Now()

	rowset, err := getGZipXmlRowSet(ctx, task, r.schema, r.Oracle)
	if err != nil {
		return nil, fmt.Errorf("db - get xml rowset error: %w", err)
	}

	updRecords, err := makeRecords(rowset.updatedRows)
	if err != nil {
		return nil, fmt.Errorf("db - convert updated rows error: %w", err)
	}

	delRecords, err := makeRecords(rowset.deletedRows)
	if err != nil {
		return nil, fmt.Errorf("db - convert deleted rows error: %w", err)
	}

	elapsed := time.Now()

	slog.Debug("db - get records",
		"elapsed", elapsed.Sub(start),
		"group_id", task.GroupId,
		"part_id", task.PartId,
		"upd_amount", len(updRecords),
		"del_amount", len(delRecords),
	)

	return append(updRecords, delRecords...), nil
}

type rowSet struct {
	updatedRows []byte
	deletedRows []byte
}

func getGZipXmlRowSet(ctx context.Context, task *model.Task, schema string, oracle *oracle.Oracle) (*rowSet, error) {
	query := "begin " +
		schema +
		".org$gate_api.getNextEvents(" +
		"  p_group_id => :1" +
		", p_part_id => :2" +
		", p_rows => :3" +
		", p_qry_columns => :4" +
		", p_qry_from => :5" +
		", p_qry_pk_column => :6" +
		", r_upd_rows_dump => :7" +
		", r_upd_rows_count => :8" +
		", r_del_rows_dump => :9" +
		", r_del_rows_count => :10" +
		"); " +
		"end;"

	tx, err := getTx(ctx, oracle.Db)
	if err != nil {
		return nil, err
	}

	var rowset rowSet
	var updRowsDump ora.Blob
	var updRowsCount int
	var delRowsDump ora.Blob
	var delRowsCount int

	_, err = tx.ExecContext(ctx, query,
		// eg: "test_tab"
		task.GroupId,
		// eg: 2
		task.PartId,
		// eg: 1000
		task.BatchSize,
		// eg: "*"
		task.Query.Columns,
		// eg: "select t.* from test_tab t"
		task.Query.From,
		// eg: "id"
		task.Query.PkColumn,

		// output
		ora.Out{Dest: &updRowsDump, Size: 1000},
		&updRowsCount,
		ora.Out{Dest: &delRowsDump, Size: 1000},
		&delRowsCount,
	)

	if err != nil {
		return nil, err
	}

	if updRowsDump.Data != nil {
		rowset.updatedRows = updRowsDump.Data
	}

	if delRowsDump.Data != nil {
		rowset.deletedRows = delRowsDump.Data
	}

	return &rowset, nil
}
