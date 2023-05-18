package test

import (
	"github.com/google/uuid"
	"learn_geektime/web"
	"learn_geektime/web/session"
	"learn_geektime/web/session/cookie"
	"learn_geektime/web/session/memory"
	"net/http"
	"testing"
	"time"
)

func TestManager(t *testing.T) {
	s := web.NewHTTPServer()
	m := session.Manager{
		SessCtxKey: "_sess",
		Store:      memory.NewStore(30 * time.Minute),
		Propagator: cookie.NewPropagator("sessid",
			cookie.WithCookieOption(func(c *http.Cookie) {
				c.HttpOnly = true
			}),
		),
	}

	s.Post("/login", func(ctx *web.Context) {
		// 前面就是一大堆登录校验

		id := uuid.New()
		sess, err := m.InitSession(ctx, id.String())
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}

		// 然后根据自己的需要设置
		err = sess.Set(ctx.Req.Context(), "mykey", "some value")
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("you are login")
	})

	s.Get("/resource", func(ctx *web.Context) {
		sess, err := m.GetSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}
		val, err := sess.Get(ctx.Req.Context(), "mykey")
		ctx.RespData = []byte(val)
		ctx.RespStatusCode = http.StatusOK
	})

	s.Post("/logout", func(ctx *web.Context) {
		_ = m.RemoveSession(ctx)
		ctx.RespData = []byte("you are logout ")
		ctx.RespStatusCode = http.StatusOK
	})

	s.Use(func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			// AOP 方案之执行校验
			if ctx.Req.URL.Path != "/login" {
				sess, err := m.GetSession(ctx)
				if err != nil {
					ctx.RespStatusCode = http.StatusUnauthorized
					return
				}
				// 用户每次请求的时候都刷新一下 session 用户体量不大的时候使用
				_ = m.Refresh(ctx.Req.Context(), sess.ID())
			}
			next(ctx)
		}
	})

	s.Start(":8081")
}
