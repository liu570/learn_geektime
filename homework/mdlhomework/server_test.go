package mdlhomework

import (
	"fmt"
	"log"
	"testing"
)

func TestHTTPServer(t *testing.T) {
	s := NewHTTPServer()

	s.Get("/user", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello user"))
	})

	s.Get("/user/*", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello user *"))
	})

	s.Get("/user/*/name", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello user * *"))
	})

	s.Post("/user/:id/:name", func(ctx *Context) {
		value, err := ctx.PathValueV1("id").ToInt64()
		if err != nil {
			t.Fatal(err)
		}
		ctx.Resp.Write([]byte(fmt.Sprintf("hello user %v", value)))
	})

	err := s.Start(":8081")
	if err != nil {
		log.Println(err)
	}
}

type User struct {
	name string
}
