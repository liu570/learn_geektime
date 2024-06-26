package querylog

import (
	"context"
	"learn_geektime/orm"
	"log"
)

type MiddlewareBuilder struct {
	logFunc func(sql string, args ...any)
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		logFunc: func(sql string, args ...any) {
			log.Printf("orm: sql: %s, args: %v", sql, args)
		},
	}
}

func (m *MiddlewareBuilder) LogFunc(logFunc func(sql string, args ...any)) *MiddlewareBuilder {
	m.logFunc = logFunc
	return m
}

func (m *MiddlewareBuilder) Build() orm.Middleware {
	return func(next orm.Handler) orm.Handler {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {

			q, err := qc.Builder.Build()
			if err != nil {
				return &orm.QueryResult{
					Err: err,
				}
			}
			m.logFunc(q.SQL, q.Args...)
			res := next(ctx, qc)
			return res
		}
	}
}
