package round_robin

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/balancer"
	"testing"
)

func TestBalancer_Pick(t *testing.T) {
	tests := []struct {
		name string
		b    *Balancer

		wantErr          error
		wantSubConn      SubConn
		wantBalanceIndex int32
	}{
		{
			name: "start",
			b: &Balancer{
				connections: []balancer.SubConn{
					SubConn{name: "127.0.0.1:8080"},
					SubConn{name: "127.0.0.1:8081"},
				},
				length: 2,
				index:  -1,
			},
			wantErr:          nil,
			wantSubConn:      SubConn{name: "127.0.0.1:8080"},
			wantBalanceIndex: 0,
		},
		{
			name: "end",
			b: &Balancer{
				connections: []balancer.SubConn{
					SubConn{name: "127.0.0.1:8080"},
					SubConn{name: "127.0.0.1:8081"},
				},
				length: 2,
				index:  1,
			},
			wantErr:          nil,
			wantSubConn:      SubConn{name: "127.0.0.1:8080"},
			wantBalanceIndex: 2,
		},

		{
			name:    "no connections",
			b:       &Balancer{},
			wantErr: balancer.ErrNoSubConnAvailable,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.b.Pick(balancer.PickInfo{})
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantSubConn.name, res.SubConn.(SubConn).name)
			assert.NotNil(t, res.Done)
			assert.Equal(t, tc.wantBalanceIndex, tc.b.index)

		})
	}
}

type SubConn struct {
	balancer.SubConn
	name string
}
