package mock

import (
	"context"
	"learn_geektime/orm"
	"time"
)

// 在数据库查询中我们通常使用sqlmock 来模拟sql返回语句 ， 但其实我们也可以通过手写middleware来模拟mock

type MiddlewareBuilder struct {
}

func (m MiddlewareBuilder) Build() orm.Middleware {
	return func(next orm.Handler) orm.Handler {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {

			val := ctx.Value(mockKey{})
			// 如果用户设置了 mock 中间件，我这个 middleware 就不会发起真的查询
			if val != nil {
				mock := val.(*Mock)
				// 模拟你的数据库查询很慢
				if mock.Sleep > 0 {
					time.Sleep(mock.Sleep)
				}
				return mock.Result
			}
			return next(ctx, qc)
		}
	}
}

type Mock struct {
	Sleep  time.Duration
	Result *orm.QueryResult
}

type mockKey struct {
}
