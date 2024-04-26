package round_robin

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

type WeightBalancer struct {
	connections []weightSubConn
	length      int32
	index       int32
}

func (w WeightBalancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	//TODO implement me
	panic("implement me")
}

type WeightBalanceBuilder struct{}

func (w *WeightBalanceBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	connections := make([]weightSubConn, 0, len(info.ReadySCs))
	for c := range info.ReadySCs {
		// TODO:加权轮询
		wConn := weightSubConn{
			conn: c,
		}
		connections = append(connections, wConn)
	}
	return &WeightBalancer{
		connections: connections,
		index:       -1,
		length:      int32(len(connections)),
	}
}

type weightSubConn struct {
	conn            balancer.SubConn
	weight          int
	currentWeight   int
	efficientWeight int
}
