package round_robin

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
	"learn_geektime/micro/route"
	"sync/atomic"
)

type Balancer struct {
	index       int32
	length      int32
	connections []subConn
	filter      route.Filter
}

func (b *Balancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	connections := make([]subConn, 0, len(b.connections))
	for _, conn := range b.connections {
		if b.filter != nil && !b.filter(info, conn.addr) {
			continue
		}
		connections = append(connections, conn)
	}
	if len(connections) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	subConn := connections[atomic.AddInt32(&b.index, 1)%int32(len(connections))]
	return balancer.PickResult{
		SubConn: subConn.c,
		Done: func(info balancer.DoneInfo) {

		},
	}, nil
}

type Builder struct {
	Filter route.Filter
}

func (b *Builder) Build(info base.PickerBuildInfo) balancer.Picker {
	connections := make([]subConn, 0, len(info.ReadySCs))
	for k, v := range info.ReadySCs {
		connections = append(connections, subConn{
			c:    k,
			addr: v.Address,
		})
	}
	return &Balancer{
		connections: connections,
		index:       -1,
		length:      int32(len(connections)),
		filter:      b.Filter,
	}
}

type subConn struct {
	c    balancer.SubConn
	addr resolver.Address
}
