package random

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
	"learn_geektime/micro/route"
	"math/rand"
)

type WeightBalancer struct {
	connections []weightSubConn
	totalWeight uint32
	filter      route.Filter
}

func (w *WeightBalancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	candidates := make([]weightSubConn, 0, len(w.connections))
	var totalWeight uint32
	for _, conn := range w.connections {
		if w.filter != nil && !w.filter(info, conn.address) {
			continue
		}
		candidates = append(candidates, conn)
		totalWeight += conn.weight
	}
	if len(candidates) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	targetWeight := rand.Intn(int(totalWeight) + 1)
	var index int
	for i, conn := range candidates {
		targetWeight -= int(conn.weight)
		if targetWeight <= 0 {
			index = i
			break
		}
	}
	return balancer.PickResult{
		SubConn: candidates[index].subConn,
		Done: func(info balancer.DoneInfo) {

		},
	}, nil

}

type WeightBalancerBuilder struct {
	Filter route.Filter
}

func (w *WeightBalancerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]weightSubConn, len(info.ReadySCs))
	var totalWeight uint32
	for conn, sc := range info.ReadySCs {
		weight := sc.Address.Attributes.Value("weight").(uint32)
		conns = append(conns, weightSubConn{
			subConn: conn,
			weight:  weight,
			address: sc.Address,
		})
		totalWeight += weight
	}
	return &WeightBalancer{
		connections: conns,
		totalWeight: totalWeight,
		filter:      w.Filter,
	}
}

type weightSubConn struct {
	subConn balancer.SubConn
	weight  uint32
	address resolver.Address
}
