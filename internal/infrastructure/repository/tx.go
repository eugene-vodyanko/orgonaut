package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type txKey struct{}

type TxManager struct {
	db *sql.DB
}

func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{db}
}

// WithinTransaction runs function within transaction
//
// The transaction commits when function were finished without error and rollback in other
func (tm *TxManager) WithinTransaction(ctx context.Context, txFunc func(ctx context.Context) error) error {
	tx, err := tm.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction error: %w", err)
	}

	err = txFunc(injectTx(ctx, tx))
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// injectTx injects transaction to context
func injectTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// extractTx extracts transaction from context
func extractTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return nil
}

func getTx(ctx context.Context, db *sql.DB) (*sql.Tx, error) {
	tx := extractTx(ctx)
	if tx != nil {
		return tx, nil
	}

	return db.Begin()
}
