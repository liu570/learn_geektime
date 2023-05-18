package repeat_body

import (
	"io"
	"learn_geektime/web"
)

func Middleware() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			ctx.Req.Body = io.NopCloser(ctx.Req.Body)
			next(ctx)
		}
	}
}
