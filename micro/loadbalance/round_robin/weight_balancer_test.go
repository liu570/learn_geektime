package round_robin

import (
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/balancer"
	"testing"
)

func TestWeightBalancer_Pick(t *testing.T) {

	weightBalancer := &WeightBalancer{
		connections: []*weightSubConn{
			{
				conn: SubConn{
					name: "weight-5",
				},
				weight:          5,
				efficientWeight: 5,
				currentWeight:   5,
			},
			{
				conn: SubConn{
					name: "weight-4",
				},
				weight:          4,
				efficientWeight: 4,
				currentWeight:   4,
			},
			{
				conn: SubConn{
					name: "weight-3",
				},
				weight:          3,
				efficientWeight: 3,
				currentWeight:   3,
			},
		},
	}
	//for i := 0; i < 12; i++ {
	//	res, err := weightBalancer.Pick(balancer.PickInfo{})
	//	require.NoError(t, err)
	//	for _, connection := range weightBalancer.connections {
	//		fmt.Print(connection.currentWeight, "\t")
	//	}
	//	fmt.Println("\t\t\t", i+1, "\t", res.SubConn.(SubConn).name)
	//}
	res, err := weightBalancer.Pick(balancer.PickInfo{})
	require.NoError(t, err)
	require.Equal(t, "weight-5", res.SubConn.(SubConn).name)
	res, err = weightBalancer.Pick(balancer.PickInfo{})
	require.NoError(t, err)
	require.Equal(t, "weight-4", res.SubConn.(SubConn).name)
	res, err = weightBalancer.Pick(balancer.PickInfo{})
	require.NoError(t, err)
	require.Equal(t, "weight-3", res.SubConn.(SubConn).name)
	res, err = weightBalancer.Pick(balancer.PickInfo{})
	require.NoError(t, err)
	require.Equal(t, "weight-5", res.SubConn.(SubConn).name)
	res, err = weightBalancer.Pick(balancer.PickInfo{})
	require.NoError(t, err)
	require.Equal(t, "weight-4", res.SubConn.(SubConn).name)
	res, err = weightBalancer.Pick(balancer.PickInfo{})
	require.NoError(t, err)
	require.Equal(t, "weight-5", res.SubConn.(SubConn).name)
}
