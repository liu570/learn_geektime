package orm

import (
	"context"
	"database/sql"
)

var (
	_ Session = &DB{}
	_ Session = &Tx{}
)

// Session 统一 Tx 与 DB 的接口
type Session interface {
	// 构建 core 的数据
	getCore() core
	queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	execContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

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

// RollbackIfNotCommit 用于在用户自己手动管理事务的时候 如果业务报错需要回滚事务，则用户需要多地管理事务，所以这里提供一个方法
// 只要用户没有提交事务 我们就直接 rollback 方便用户回滚
func (tx *Tx) RollbackIfNotCommit() error {
	err := tx.tx.Rollback()
	if err == sql.ErrTxDone {
		// 已提交错误 直接忽略
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
