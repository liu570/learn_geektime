package mdlhomework

import (
	"errors"
	"net/http"
)

type Server interface {
	Start(addr string)
	AddRoute(method string, path string, handleFunc HandleFunc)
}

type HTTPServer struct {
	router
}

func (s *HTTPServer) Start(addr string) error {
	if addr[0] != ':' {
		return errors.New("lossed ':' in your address")
	}
	return http.ListenAndServe(addr, s)
}

func (s *HTTPServer) AddRoute(method string, path string, handleFunc HandleFunc) {
	s.addRoute(method, path, handleFunc)
}

func (s *HTTPServer) Get(path string, handlerFunc HandleFunc) {
	s.AddRoute(http.MethodGet, path, handlerFunc)
}

func (s *HTTPServer) Post(path string, handlerFunc HandleFunc) {
	s.AddRoute(http.MethodPost, path, handlerFunc)
}

func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	s.serve(ctx)
}

func (s *HTTPServer) serve(ctx *Context) {
	mi, ok := s.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok {
		ctx.Resp.WriteHeader(http.StatusNotFound)
		ctx.Resp.Write([]byte("Not Found"))
		return
	}
	n := mi.node
	ctx.pathParams = mi.pathParams
	n.handle(ctx)
}

func NewHTTPServer() *HTTPServer {
	s := &HTTPServer{
		router: newRouter(),
	}
	return s
}
