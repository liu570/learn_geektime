package random

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"math/rand"
)

type Balancer struct {
	connections []balancer.SubConn
	length      int
}

func (b *Balancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(b.connections) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	return balancer.PickResult{
		SubConn: b.connections[rand.Intn(b.length)],
		Done: func(info balancer.DoneInfo) {

		},
	}, nil

}

type BalancerBuilder struct {
}

func (b *BalancerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	connections := make([]balancer.SubConn, 0, len(info.ReadySCs))
	for conn, _ := range info.ReadySCs {
		connections = append(connections, conn)
	}
	return &Balancer{
		connections: connections,
		length:      len(connections),
	}
}
