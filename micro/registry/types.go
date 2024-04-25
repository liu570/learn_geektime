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
}
type Event struct {
}
