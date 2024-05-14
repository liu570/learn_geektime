package prometheus

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"learn_geektime/web"
	"time"
)

type MiddlewareBuilder struct {
	Name        string
	Subsystem   string
	ConstLabels map[string]string
	Help        string
}

func (m *MiddlewareBuilder) Build() web.Middleware {
	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:        m.Name,
		Subsystem:   m.Subsystem,
		ConstLabels: m.ConstLabels,
		Help:        m.Help,
	}, []string{"pattern", "method", "status"})
	//必须使用否则 prometheus 无用
	prometheus.MustRegister(summaryVec)
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			startTime := time.Now()
			defer func() {
				endTime := time.Now()
				duration := endTime.Sub(startTime).Milliseconds()
				summary := summaryVec.WithLabelValues(ctx.MatchedRoute, ctx.Req.Method,
					fmt.Sprintf("%d", ctx.RespStatusCode))
				summary.Observe(float64(duration))
			}()
			next(ctx)
		}
	}
}
