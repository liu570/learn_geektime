package no_delete

import (
	"context"
	"errors"
	"learn_geektime/orm"
)

type MiddlewareBuilder struct {
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{}
}

func (m *MiddlewareBuilder) Build() orm.Middleware {
	return func(next orm.Handler) orm.Handler {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			if qc.Type == "DElETE" {
				return &orm.QueryResult{
					Err: errors.New("orm: no_delete_middleware no delete 禁止 DELETE 语句"),
				}
			}
			res := next(ctx, qc)
			return res
		}
	}
}
