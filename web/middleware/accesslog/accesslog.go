package accesslog

import (
	"encoding/json"
	"io"
	"learn_geektime/web"
)

type MiddlewareBuilder struct {
	logFunc func(accesslog []byte)
}

func (m MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			body, err := io.ReadAll(ctx.Req.Body)
			l := accessLog{
				Method: ctx.Req.Method,
				Body:   string(body),
			}
			bs, err := json.Marshal(l)
			if err == nil {
				m.logFunc(bs)
			}
			next(ctx)
			m.logFunc(bs)
		}
	}
}

type accessLog struct {
	Method string
	Body   string
}
