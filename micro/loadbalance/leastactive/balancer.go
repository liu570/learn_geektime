package leastactive

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"math"
	"sync/atomic"
)

type Balance struct {
	connections []*activeConn
}

func (b *Balance) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	res := &activeConn{
		cnt: math.MaxInt,
	}
	for _, connection := range b.connections {
		if connection.cnt < res.cnt {
			res = connection
		}
	}
	atomic.AddUint32(&res.cnt, 1)
	return balancer.PickResult{
		SubConn: res.conn,
		Done: func(info balancer.DoneInfo) {
			atomic.AddUint32(&res.cnt, -1)
		},
	}, nil
}

type BalancerBuilder struct {
}

func (b *BalancerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	connections := make([]activeConn, 0, len(info.ReadySCs))
	for c := range info.ReadySCs {
		connections = append(connections, activeConn{
			conn: c,
		})
	}
	return &Balance{connections: make([]*activeConn, 0)}
}

type activeConn struct {
	cnt  uint32
	conn balancer.SubConn
}
