package orm

import (
	"context"
	"learn_geektime/orm/model"
)

// QueryContext Query语句中传递上下文
type QueryContext struct {
	// 标明用户是使用在 UPDATE,DELETE,SELECT,INSERT 中的哪一个
	Type    string
	Builder QueryBuilder

	// 使用 openTelemetry 的时候需要用到
	Model     *model.Model
	TableName string
}

// QueryResult Query语句中的返回值
type QueryResult struct {
	// SELECT 语句，你的返回值是 T 或是 []T
	// UPDATE,SELECT,DELETE 返回值是 Result
	Result any
	// 查询本身出的错误
	Err error
}

type Handler func(ctx context.Context, qc *QueryContext) *QueryResult

type Middleware func(next Handler) Handler
