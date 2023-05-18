package orm

import (
	"context"
	"database/sql"
)

// Tx 对 sql.Tx 的封装
type Tx struct {
	core
	tx *sql.Tx
}

func (tx *Tx) Commit() error {
	return tx.tx.Commit()
}

func (tx *Tx) Rollback() error {
	return tx.tx.Rollback()
}
func (tx *Tx) RollbackIfNotCommit() error {
	err := tx.tx.Rollback()
	if err == sql.ErrTxDone {
		return nil
	}
	return err
}

func (tx *Tx) getCore() core {
	return tx.core
}
func (tx *Tx) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return tx.tx.QueryContext(ctx, query, args...)
}

func (tx *Tx) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return tx.tx.ExecContext(ctx, query, args...)
}

// Session 统一 Tx 与 DB 的接口
type Session interface {
	getCore() core
	queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	execContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}
