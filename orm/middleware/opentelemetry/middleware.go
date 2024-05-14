package opentelemetry

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"learn_geektime/orm"
)

type MiddlewareBuilder struct {
	//tracer
	tracer trace.Tracer
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{}
}

func (m *MiddlewareBuilder) Build() orm.Middleware {
	return func(next orm.Handler) orm.Handler {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			tbl := qc.Model.TableName
			spanCtx, span := m.tracer.Start(ctx, fmt.Sprintf("%s-%s", qc.Type, tbl))
			defer span.End()
			q, err := qc.Builder.Build()
			if err != nil {
				return &orm.QueryResult{
					Err: err,
				}
			}
			if q != nil {
				span.SetAttributes(attribute.String("sql", q.SQL))
			}
			span.SetAttributes(attribute.String("table", tbl))
			span.SetAttributes(attribute.String("component", "orm"))
			res := next(spanCtx, qc)
			if res.Err != nil {
				span.RecordError(res.Err)
			}
			return res
		}
	}
}
