package repository

import (
	"context"
	"github.com/eugene-vodyanko/orgonaut/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var task = &model.Task{
	GroupId:   "test_tab",
	PartId:    5,
	BatchSize: 1000,
	Query: model.Query{
		From:     "select t.* from test_tab t",
		Columns:  "*",
		PkColumn: "id",
	},
}

func TestRepository_getXmlRowSet(t *testing.T) {
	start := time.Now()

	db, schema, teardown := TestOra(t)
	defer teardown()

	tm := NewTxManager(db.Db)

	ctx := context.Background()
	err := tm.WithinTransaction(ctx, func(txCtx context.Context) error {
		rs, err := getGZipXmlRowSet(txCtx, task, schema, db)

		if err == nil {
			t.Logf("upd_size: %d, del_size: %d", len(rs.updatedRows), len(rs.deletedRows))
		}

		return err
	})

	assert.NoError(t, err)
	elapsed := time.Now()

	t.Logf("elapsed: %v", elapsed.Sub(start))
}

func TestRepository_GetRecords(t *testing.T) {
	start := time.Now()

	db, schema, teardown := TestOra(t)
	defer teardown()

	repo := NewRepository(schema, db)
	tm := NewTxManager(db.Db)

	ctx := context.Background()
	err := tm.WithinTransaction(ctx, func(txCtx context.Context) error {

		records, err := repo.GetRecords(txCtx, task)

		if err == nil {
			t.Logf("records size: %d", len(records))
		}

		return err
	})
	assert.NoError(t, err)
	elapsed := time.Now()

	t.Logf("elapsed: %v", elapsed.Sub(start))
}
