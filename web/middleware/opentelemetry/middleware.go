package opentelemetry

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"learn_geektime/web"
)

type MiddlewareBuilder struct {
	//tracer
	tracer trace.Tracer
}

func (m *MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			spanCtx, span := m.tracer.Start(ctx.Req.Context(), "Unknown")
			defer span.End()
			span.SetAttributes(attribute.String("http.method", ctx.Req.Method))
			//请求路径
			// [:1024] 防止攻击者传入超长的url导致我们的span被压垮
			span.SetAttributes(attribute.String("http.path", ctx.Req.URL.Path[:1024]))
			span.SetAttributes(attribute.String("peer.hostname", ctx.Req.Host))
			span.SetAttributes(attribute.String("http.url", ctx.Req.URL.String()))
			span.SetAttributes(attribute.String("http.scheme", ctx.Req.URL.Scheme))
			span.SetAttributes(attribute.String("span.kind", "server"))
			span.SetAttributes(attribute.String("component", "web"))
			span.SetAttributes(attribute.String("peer.address", ctx.Req.RemoteAddr))
			span.SetAttributes(attribute.String("http.proto", ctx.Req.Proto))

			//将spanCtx传入请求中的ctx来实现tracing
			// TODO 此地会频繁的创建 context 对性能的影响很大
			ctx.Req = ctx.Req.WithContext(spanCtx)
			// ctx.Ctx = spanCtx
			next(ctx)
			// 判断是否找到路由 而不是404
			if ctx.MatchedRoute != "" {
				span.SetName(ctx.MatchedRoute)
			}
			span.SetAttributes(attribute.Int("http.status",
				ctx.RespStatusCode))

		}
	}
}
