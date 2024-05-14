package prometheus

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"learn_geektime/orm"
	"time"
)

type MiddlewareBuilder struct {
	Name        string
	Subsystem   string
	ConstLabels map[string]string
	Help        string
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{}
}

func (m *MiddlewareBuilder) Build() orm.Middleware {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:        m.Name,
		Subsystem:   m.Subsystem,
		Namespace:   m.Name,
		ConstLabels: m.ConstLabels,
		Help:        m.Help,
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.90:  0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, []string{"type", "table"})
	//必须使用否则 prometheus 无用
	prometheus.MustRegister(vector)
	return func(next orm.Handler) orm.Handler {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			startTime := time.Now()
			// 同时也可以记录
			// errCounterVec 记录错误数
			// histogram
			// active query
			defer func() {
				// 记录执行时间
				vector.WithLabelValues(qc.Type, qc.Model.TableName).Observe(float64(time.Since(startTime).Milliseconds()))
			}()
			res := next(ctx, qc)
			return res
		}
	}
}
