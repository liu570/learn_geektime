package web

import (
	"fmt"
	"testing"
)

func TestServer(t *testing.T) {
	s := NewHTTPServer()
	// 组建AOP方面的责任链
	//s.Use(repeat_body.Middleware(), accesslog.MiddlewareBuilder{}.Build())
	s.Get("/", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello world"))
	})
	s.Post("/user", func(ctx *Context) {
		u := &User{}
		err := ctx.BindJSON(u)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(u)
	})

	s.Post("/fileupload", (&FileUploader{}).Handle())

	s.Get("/user", func(ctx *Context) {
		//age, err := ctx.PathValueV1("age").ToInt64()
		ctx.Resp.Write([]byte("hello user"))
	})
	//-----------------------------------------------------
	// web框架一般情况下是不需要线程安全的因为我们都是先注册好了路由才开始启动服务器的
	s.Start("8081")
}

type User struct {
	Name string
}
