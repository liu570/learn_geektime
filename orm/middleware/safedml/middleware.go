package safedml

import (
	"context"
	"errors"
	"fmt"
	"learn_geektime/orm"
	"strings"
)

// 强制查询语句 带 where
// 1、SELECT、DELETE、UPDATE 必须带 where
// 2、DELETE、UPDATE 必须带 where 因为 SELECT 语句不涉及修改数据 所以可要可不要

type MiddlewareBuilder struct {
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{}
}

func (m *MiddlewareBuilder) Build() orm.Middleware {
	return func(next orm.Handler) orm.Handler {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			if qc.Type == "SELECT" || qc.Type == "INSERT" {
				return next(ctx, qc)
			}
			q, err := qc.Builder.Build()
			if err != nil {
				return &orm.QueryResult{
					Err: err,
				}
			}
			if strings.Contains(q.SQL, "WHERE") {
				return &orm.QueryResult{
					Err: errors.New(fmt.Sprintf("orm: safedml_middleware SQL type %s can not contains `where`", qc.Type)),
				}
			}
			res := next(ctx, qc)
			return res
		}
	}
}
