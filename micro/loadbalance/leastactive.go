package loadbalance

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

type Balance struct {
	connections []*activeConn
}

func (b *Balance) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	//var res *activeConn
	//for _, connection := range b.connections {
	//
	//}
	// TODO 未完成
	panic("implement me")
}

type BalancerBuilder struct {
}

func (b *BalancerBuilder) Build(base.PickerBuildInfo) balancer.Picker {
	//TODO implement me
	panic("implement me")
}

type activeConn struct {
	cnt  int
	conn balancer.SubConn
}
