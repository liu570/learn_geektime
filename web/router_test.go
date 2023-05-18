package web

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_router_addRoute(t *testing.T) {

	tests := []struct {
		name string
		// 输入
		method     string
		path       string
		handleFunc HandleFunc
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail/:order_sn",
		},
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		{
			method: http.MethodGet,
			path:   "/order/*/detail",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
		// 通配符测试用例
	}

	var handleFunc HandleFunc = func(context *Context) {

	}
	println(handleFunc)

	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path:    "/",
				handler: handleFunc,
				children: map[string]*node{
					"user": &node{
						path:    "user",
						handler: handleFunc,
						children: map[string]*node{
							"home": &node{path: "home", handler: handleFunc},
						},
					},
					"order": &node{
						path: "order",
						starChild: &node{
							path:    "*",
							handler: handleFunc,
							children: map[string]*node{
								"detail": &node{
									path:    "detail",
									handler: handleFunc,
								},
							},
						},
						children: map[string]*node{
							"detail": &node{
								path:    "detail",
								handler: handleFunc,
								paramChild: &node{
									path:    ":order_sn",
									handler: handleFunc,
								},
							},
						},
					},
				},
			},
			http.MethodPost: &node{
				path: "/",
				children: map[string]*node{
					"login": &node{path: "login", handler: handleFunc},
					"order": &node{
						path: "order",
						children: map[string]*node{
							"create": &node{path: "create", handler: handleFunc},
						},
					},
				},
			},
		},
	}

	res := &router{
		trees: map[string]*node{},
	}

	for _, tc := range tests {
		res.addRoute(tc.method, tc.path, handleFunc)
	}

	println(wantRouter)
	// 这里直接断言整棵路由树，所以不需要循环断言
	assert.Equal(t, wantRouter, res)

	findcases := []struct {
		name   string
		method string
		path   string

		found    bool
		wantPath string
	}{
		{
			name:     " /",
			method:   http.MethodGet,
			path:     "/",
			found:    true,
			wantPath: "/",
		},
		{
			name:     " /user",
			method:   http.MethodGet,
			path:     "/user",
			found:    true,
			wantPath: "user",
		},
		{
			name:     " /order/*/detail",
			method:   http.MethodGet,
			path:     "/order/*/detail",
			found:    true,
			wantPath: "detail",
		},
		{
			name:   " /order",
			method: http.MethodGet,
			path:   "/order",
			found:  false,
		},
		{
			name:     " /order/detail",
			method:   http.MethodGet,
			path:     "/order/detail",
			found:    true,
			wantPath: "detail",
		},
		{
			name:     " /order/*",
			method:   http.MethodGet,
			path:     "/order/abc",
			found:    true,
			wantPath: "*",
		},
		{
			name:     " /order/detail/:order_sn",
			method:   http.MethodGet,
			path:     "/order/detail/:order_sn",
			found:    true,
			wantPath: ":order_sn",
		},
	}

	for _, tc := range findcases {
		t.Run(tc.name, func(t *testing.T) {
			mt, ok := res.findRoute(http.MethodGet, tc.path)
			assert.Equal(t, tc.found, ok)
			if !ok {
				return
			}
			assert.Equal(t, tc.wantPath, mt.n.path)
			assert.NotNil(t, mt.n.handler)
		})

	}

}
