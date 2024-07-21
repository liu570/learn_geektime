package round_robin

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
	"learn_geektime/micro/route"
	"math"
	"sync"
)

type WeightBalancer struct {
	connections []*weightSubConn
	mutex       sync.Mutex
	filter      route.Filter
}

func (w *WeightBalancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	var resConn *weightSubConn
	var totalWeight uint32
	for _, conn := range w.connections {
		if w.filter != nil && !w.filter(info, conn.address) {
			continue
		}
		conn.mutex.Lock()
		totalWeight += conn.efficientWeight
		conn.currentWeight += conn.efficientWeight
		if resConn == nil || resConn.currentWeight < conn.currentWeight {
			resConn = conn
		}
		conn.mutex.Unlock()
	}
	if resConn == nil {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	resConn.mutex.Lock()
	resConn.currentWeight -= totalWeight
	resConn.mutex.Unlock()
	return balancer.PickResult{
		SubConn: resConn.conn,
		Done:    w.done(resConn),
	}, nil
}

func (w *WeightBalancer) done(resConn *weightSubConn) func(info balancer.DoneInfo) {
	return func(info balancer.DoneInfo) {
		resConn.mutex.Lock()
		if info.Err != nil && resConn.efficientWeight == 0 {
			return
		}
		if info.Err == nil && resConn.efficientWeight == math.MaxUint32 {
			return
		}
		if info.Err != nil {
			resConn.efficientWeight--
		} else {
			resConn.efficientWeight++
		}
		resConn.mutex.Unlock()
	}
}

type WeightBalanceBuilder struct {
	Filter route.Filter
}

func (w *WeightBalanceBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	connections := make([]*weightSubConn, 0, len(info.ReadySCs))
	for sub, subInfo := range info.ReadySCs {
		weight := subInfo.Address.Attributes.Value("weight").(uint32)
		wConn := weightSubConn{
			conn:            sub,
			weight:          weight,
			efficientWeight: weight,
			currentWeight:   weight,
			address:         subInfo.Address,
		}
		connections = append(connections, &wConn)
	}
	return &WeightBalancer{
		connections: connections,
		filter:      w.Filter,
	}
}

type weightSubConn struct {
	conn            balancer.SubConn
	address         resolver.Address
	mutex           sync.Mutex
	weight          uint32
	currentWeight   uint32
	efficientWeight uint32
}
