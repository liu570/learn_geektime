package hash

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

type Balancer struct {
}

func (b *Balancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	// 由于 grpc 的缺陷 我们无法拿到请求 无法根据请求特性设置hash值做负载均衡
	// 同理一致性 hash
	panic("implement me")
}

type BalancerBuilder struct {
}

func (b *BalancerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {

	panic("implement me")
}
