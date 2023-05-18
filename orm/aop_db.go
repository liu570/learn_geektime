package orm

// 该文件不是 orm 代码 只做演示使用

//import (
//	"context"
//	"database/sql"
//)
//
//type AopDBContext struct {
//	query string
//	args  []any
//}
//type AopDbResult struct {
//	rows *sql.Rows
//	err  error
//}
//
//type Handler func(actx *AopDBContext) *AopDbResult
//type Middleware func(next Handler) Handler
//type AopDB struct {
//	db *sql.DB
//	ms []Middleware
//}
//
//func (db *AopDB) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
//	var handler Handler = func(actx *AopDBContext) *AopDbResult {
//		rows, err := db.db.QueryContext(ctx, actx.query, actx.args...)
//		return &AopDbResult{
//			rows: rows,
//			err:  err,
//		}
//	}
//
//	for i := len(db.ms) - 1; i >= 0; i-- {
//		handler = db.ms[i](handler)
//	}
//
//	res := handler(&AopDBContext{})
//	return res.rows, res.err
//}
