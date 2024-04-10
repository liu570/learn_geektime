package rpc

import "context"

type Service interface {
	Name() string
}
type Proxy interface {
	Invoke(ctx context.Context, req *Request) (*Response, error)
}

// Request rpc 调用所需要的调用信息
type Request struct {
	ServiceName string
	MethodName  string
	Arg         []byte
}

type Response struct {
	data []byte
}
