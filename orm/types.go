// 此地为 orm 的核心接口定义
package orm

import (
	"context"
)

// Queier represent SELECT 语句
type Queier[T any] interface {
	Get(ctx context.Context) (*T, error)
	GetMulti(ctx context.Context) ([]*T, error)
}

// Executor represent UPDATE, DELETE, INSERT
type Executor interface {
	Exec(ctx context.Context) Result
}

type QueryBuilder interface {
	// Build 用于构建 SQL 语句
	Build() (*Query, error)
}

// db.Exec
// db.QueryRow
// db.Query
type Query struct {
	SQL  string
	Args []any
}
