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
	// 权重信息 用于表示服务器的处理能力
	Weight uint32
	// 分组信息 用于表示服务实例的组别
	Group string

	// 其它信息也可以是使用 map来实现， 但是上面更加直观
	//Attributes map[string]string
}
type Event struct {
}
