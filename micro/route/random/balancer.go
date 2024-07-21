package random

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
	"learn_geektime/micro/route"
	"math/rand"
)

type Balancer struct {
	connections []subConn
	length      int
	filter      route.Filter
}

func (b *Balancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	candidates := make([]subConn, 0, len(b.connections))
	for _, conn := range b.connections {
		if b.filter != nil && !b.filter(info, conn.address) {
			continue
		}
		candidates = append(candidates, conn)
	}
	if len(candidates) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	return balancer.PickResult{
		SubConn: candidates[rand.Intn(len(candidates))].conn,
		Done: func(info balancer.DoneInfo) {

		},
	}, nil

}

type BalancerBuilder struct {
	Filter route.Filter
}

func (b *BalancerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	connections := make([]subConn, 0, len(info.ReadySCs))
	for conn, subInfo := range info.ReadySCs {
		connections = append(connections, subConn{
			conn:    conn,
			address: subInfo.Address,
		})
	}
	return &Balancer{
		connections: connections,
		length:      len(connections),
		filter:      b.Filter,
	}
}

type subConn struct {
	conn    balancer.SubConn
	address resolver.Address
}
