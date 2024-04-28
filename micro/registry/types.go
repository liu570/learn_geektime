package registry

import (
	"context"
	"io"
)

type Registry interface {
	Register(ctx context.Context, si ServiceInstance) error
	UnRegister(ctx context.Context, si ServiceInstance) error
	ListServices(ctx context.Context, serviceName string) ([]ServiceInstance, error)
	Subscribe(serviceName string) (<-chan Event, error)
	io.Closer
}
type ServiceInstance struct {
	Name string
	// 这里表示是你的定位信息
	Address string

	// 其它信息
	// 定位信息 用于
	// TODO 服务端定位信息传递未完成
	Weight uint32
}
type Event struct {
}
