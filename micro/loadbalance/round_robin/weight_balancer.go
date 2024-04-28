package round_robin

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"math"
	"sync"
)

type WeightBalancer struct {
	connections []*weightSubConn
	mutex       sync.Mutex
}

func (w *WeightBalancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if w.connections == nil || len(w.connections) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	var resConn *weightSubConn
	var totalWeight uint32
	for _, conn := range w.connections {
		conn.mutex.Lock()
		totalWeight += conn.efficientWeight
		conn.currentWeight += conn.efficientWeight
		if resConn == nil || resConn.currentWeight < conn.currentWeight {
			resConn = conn
		}
		conn.mutex.Unlock()
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

type WeightBalanceBuilder struct{}

func (w *WeightBalanceBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	connections := make([]*weightSubConn, 0, len(info.ReadySCs))
	for sub, subInfo := range info.ReadySCs {
		weight := subInfo.Address.Attributes.Value("weight").(uint32)
		wConn := weightSubConn{
			conn:            sub,
			weight:          weight,
			efficientWeight: weight,
			currentWeight:   weight,
		}
		connections = append(connections, &wConn)
	}
	return &WeightBalancer{
		connections: connections,
	}
}

type weightSubConn struct {
	conn            balancer.SubConn
	mutex           sync.Mutex
	weight          uint32
	currentWeight   uint32
	efficientWeight uint32
}
