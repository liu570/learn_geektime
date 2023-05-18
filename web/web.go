package web

import (
	"net"
	"net/http"
)

// HTTPSServer 使用教程
func Start() {
	var s Server = &HTTPSServer{
		// 通过此地确认Server 接口的 实现结构体HTTPServer
		Server: &HTTPServer{},
	}
	var h1 HandleFunc
	var h2 HandleFunc
	s.AddRoute(http.MethodPost, "/user", func(context *Context) {
		// 循环调用多个handlefunc
		h1(context)
		h2(context)
	})

	s.Start(":8082")
}

type Server interface {
	http.Handler

	// 监听端口并启动服务器
	Start(addr string) error

	// AddRoute 注册路由的核心抽象
	AddRoute(method, path string, handler HandleFunc)

	// 如果中间件提供注册多个路由的方法
	// 中间件没法考虑用户如何调度的问题，因为不同的业务需求不同
	// 中间发生中断该如何处理，也是需要考虑的问题
	// 同时 用户也可能一个路由都不传
	//AddRoutes(method, path string, handler ...HandleFunc)
}

// 用于确保 HTTPServer 肯定实现了 Server 接口
var _ Server = &HTTPServer{}

type ServerOption func(server *HTTPServer)

// HTTPServer 监管http请求
type HTTPServer struct {
	router
	ms []Middleware
	// 模板引擎
	tplEngine TemplateEngine
}

// AddRoute 实现该方法
func (m *HTTPServer) AddRoute(method, path string, handler HandleFunc) {
	m.addRoute(method, path, handler)
}

func (m *HTTPServer) Get(path string, handler HandleFunc) {
	m.AddRoute(http.MethodGet, path, handler)
}

func (m *HTTPServer) Post(path string, handler HandleFunc) {
	m.AddRoute(http.MethodPost, path, handler)
}

func (m *HTTPServer) Start(addr string) error {
	// 端口启动前
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		return err
	}
	// 端口启动后

	// 可以在这里注册本服务器到你的管理平台
	// 比如说你注册到 etcd，然后你打开管理界面，你就能看到这个实例
	// 10.0.0.1：8081

	return http.Serve(listener, m)
	// 这个是阻塞的
	//return http.ListenAndServe(addr, m)
}

// ServeHTTP 为 Handle接口的方法代表 HTTPServer实现了Handle接口
func (m *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:        request,
		Resp:       writer,
		tplEngine:  m.tplEngine,
		UserValues: make(map[string]any),
	}
	// 构造责任链
	// 这里就是AOP方案，使用多个middleware 来处理各种非业务逻辑
	root := m.serve
	for i := len(m.ms) - 1; i >= 0; i-- {
		root = m.ms[i](root)
	}
	// 创建一个flushMiddleware 使得返回前端页面的响应头，和响应body是我们context维护的数据
	var flushMdl Middleware = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
			m.writeResp(ctx)
		}
	}
	root = flushMdl(root)

	// 构建到终点后执行你的业务逻辑
	root(ctx)
}

func (m *HTTPServer) serve(ctx *Context) {
	// before route
	mi, ok := m.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || mi.n == nil || mi.n.handler == nil {
		ctx.RespStatusCode = http.StatusNotFound
		ctx.RespData = []byte("Not Found")
		return
	}
	ctx.pathParams = mi.pathParams
	ctx.MatchedRoute = mi.n.route
	//  before execute
	mi.n.handler(ctx)
	// after execute

	// after route

}

func (m *HTTPServer) Group(prefix string) *Group {
	return &Group{
		prefix: prefix,
		s:      m,
	}
}

func (m *HTTPServer) Use(ms ...Middleware) {
	m.ms = ms
	//m.ms = append(m.ms, ms...)
}

func (s *HTTPServer) writeResp(ctx *Context) {
	ctx.Resp.WriteHeader(ctx.RespStatusCode)
	ctx.Resp.Write(ctx.RespData)
}

func NewHTTPServer(opts ...ServerOption) *HTTPServer {
	s := &HTTPServer{
		router: newRouter(),
	}
	for _, opts := range opts {
		opts(s)
	}
	return s
}

func ServerWithTemplateEngine(engine TemplateEngine) ServerOption {
	return func(server *HTTPServer) {
		server.tplEngine = engine
	}
}

// HTTPSServer 典型装饰器模式
type HTTPSServer struct {
	Server
	CertFile string
	KeyFile  string
}

func (m *HTTPSServer) Start(addr string) error {
	return http.ListenAndServeTLS(addr, m.CertFile, m.KeyFile, m)
}

type Group struct {
	prefix string
	s      Server
}

func (g *Group) AddRoute(method, path string, handler HandleFunc) {
	g.s.AddRoute(method, g.prefix+path, handler)
}
