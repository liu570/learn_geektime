package random

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"math/rand"
)

type WeightBalancer struct {
	connections []weightSubConn
	totalWeight uint32
}

func (w *WeightBalancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(w.connections) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	targetWeight := rand.Intn(int(w.totalWeight) + 1)
	var index int
	for i, conn := range w.connections {
		targetWeight -= int(conn.weight)
		if targetWeight <= 0 {
			index = i
			break
		}
	}
	return balancer.PickResult{
		SubConn: w.connections[index].subConn,
		Done: func(info balancer.DoneInfo) {

		},
	}, nil

}

type WeightBalancerBuilder struct{}

func (w *WeightBalancerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]weightSubConn, len(info.ReadySCs))
	var totalWeight uint32
	for conn, sc := range info.ReadySCs {
		weight := sc.Address.Attributes.Value("weight").(uint32)
		conns = append(conns, weightSubConn{
			subConn: conn,
			weight:  weight,
		})
		totalWeight += weight
	}
	return &WeightBalancer{
		connections: conns,
		totalWeight: totalWeight,
	}
}

type weightSubConn struct {
	subConn balancer.SubConn
	weight  uint32
}
