package micro

import (
	"context"
	"google.golang.org/grpc"
	"learn_geektime/micro/registry"
	"net"
	"time"
)

type ServerOption func(server *Server)

type Server struct {
	name            string
	listener        net.Listener
	registry        registry.Registry
	registerTimeout time.Duration
	*grpc.Server
}

func NewServer(name string, opts ...ServerOption) (*Server, error) {
	server := &Server{
		name:            name,
		Server:          grpc.NewServer(),
		registerTimeout: time.Second * 3,
	}
	for _, opt := range opts {
		opt(server)
	}
	return server, nil
}

func (s *Server) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.listener = listener
	// 开始注册，但这里需要判断，不是所有的 server 都需要注册
	if s.registry != nil {
		ctx, cancel := context.WithTimeout(context.Background(), s.registerTimeout)
		defer cancel()
		err := s.registry.Register(ctx, registry.ServiceInstance{
			Name: s.name,
			// 这里 ip 端口地址 如果是在容器中使用，容器的 ip 地址需要映射不可直接这样使用
			Address: listener.Addr().String(),
		})
		if err != nil {
			return err
		}
	}
	err = s.Serve(listener)
	return err
}

func (s *Server) Close() error {
	// 这里关闭的时候需要先 从注册中心中解注册才能关闭链接
	if s.registry != nil {
		err := s.registry.Close()
		return err
	}
	// 这里关闭链接 可以使用 grpc 的优雅退出
	// server 里面有 listener 所以会关闭 链接
	s.GracefulStop()
	return nil
}

func ServerWithRegistry(r registry.Registry) ServerOption {
	return func(server *Server) {
		server.registry = r
	}
}
