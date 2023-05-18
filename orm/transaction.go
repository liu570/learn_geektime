package orm

import (
	"context"
	"database/sql"
	"learn_geektime/orm/internal/valuer"
	"learn_geektime/orm/model"
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
	// 返回 Registry, valCreator
	getCore() core
	queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	execContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// core 统一 Tx 与 DB 都需要的数据
type core struct {
	r model.Registry
	// valCreator 来确定我们是使用 反射 还是 unsafe
	valCreator valuer.Creator
	// dialect 确定我们数据库使用哪个数据库的方言
	dialect Dialect
	// AOP 方案
	ms []Middleware
}
