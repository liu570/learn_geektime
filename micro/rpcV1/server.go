package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"reflect"
)

type Server struct {
	network string
	addr    string

	services map[string]reflectionStub
}

func NewServer(network string, addr string) *Server {
	return &Server{
		network:  network,
		addr:     addr,
		services: make(map[string]reflectionStub, 16),
	}
}
func (s *Server) RegisterServer(service Service) {
	s.services[service.Name()] = reflectionStub{
		s:     service,
		value: reflect.ValueOf(service),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen(s.network, s.addr)
	if err != nil {
		// 常见错误为端口被占用
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go func() {
			if er := s.handleConn(conn); er != nil {
				_ = conn.Close()
			}
		}()
	}
}

// handleConn 我们可以任务 一个请求包含两个部分
// 1. 长度字段：8字节
// 2. 请求数据
func (s *Server) handleConn(conn net.Conn) error {
	for {
		reqBs, err := ReadMsg(conn)
		if err != nil {
			return err
		}

		req := &Request{}
		err = json.Unmarshal(reqBs, req)
		if err != nil {
			return err
		}
		// 还原了调用信息，在后面我需要发起业务调用
		resp, err := s.Invoke(context.Background(), req)
		if err != nil {
			//	这个可能是你的业务error
			//	暂时不知道怎么回传 error 所以简单记录下
			return err
		}

		res := EncodeMsg(resp.data)
		_, err = conn.Write(res)
		if err != nil {
			return err
		}
	}
}

func (s *Server) Invoke(ctx context.Context, req *Request) (*Response, error) {
	service, ok := s.services[req.ServiceName]
	if !ok {
		return nil, errors.New("micro:调用的服务不存在")
	}
	resp, err := service.Invoke(ctx, req.MethodName, req.Arg)
	if err != nil {
		return nil, err
	}
	return &Response{
		data: resp,
	}, nil
}

type reflectionStub struct {
	s     Service
	value reflect.Value
}

func (r *reflectionStub) Invoke(ctx context.Context, methodName string, data []byte) ([]byte, error) {
	method := r.value.MethodByName(methodName)
	in := make([]reflect.Value, 2)
	in[0] = reflect.ValueOf(context.Background())
	inReq := reflect.New(method.Type().In(1).Elem())
	err := json.Unmarshal(data, inReq.Interface())
	if err != nil {
		return nil, err
	}
	in[1] = inReq

	results := method.Call(in)
	if results[1].Interface() != nil {
		return nil, results[1].Interface().(error)
	}

	return json.Marshal(results[0].Interface())
}
