package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"learn_geektime/micro/rpc/message"
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

		req := message.DecodeReq(reqBs)
		if err != nil {
			return err
		}
		// 还原了调用信息，在后面我需要发起业务调用
		resp, err := s.Invoke(context.Background(), req)
		if err != nil {
			//	这个可能是你的业务error
			//	处理业务 error 不需要直接中断,直接传回业务error
			resp.Error = []byte(err.Error())
		}

		respBs := message.EncodeResp(resp)
		_, err = conn.Write(respBs)
		if err != nil {
			return err
		}
	}
}

func (s *Server) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	service, ok := s.services[req.ServiceName]
	resp := &message.Response{
		MesssageID:  req.MesssageID,
		Version:     req.Version,
		Compression: req.Compression,
		Serialize:   req.Serialize,
	}

	if !ok {
		return resp, errors.New("micro:调用的服务不存在")
	}

	respData, err := service.Invoke(ctx, req.MethodName, req.Data)
	resp.Data = respData
	if err != nil {
		return resp, err
	}
	return resp, nil
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
		err = results[1].Interface().(error)
	}
	res, er := json.Marshal(results[0].Interface())
	if er != nil {
		return nil, er
	}
	return res, err
}
