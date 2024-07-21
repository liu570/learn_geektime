package registry

import (
	"context"
	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"
	"learn_geektime/micro"
	"learn_geektime/micro/proto/gen"
	"learn_geektime/micro/registry/etcd"
	"learn_geektime/micro/route"
	"learn_geektime/micro/route/round_robin"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	require.NoError(t, err)
	r, err := etcd.NewRegistry(etcdClient)

	require.NoError(t, err)
	client, err := micro.NewClient(micro.ClientWithInsecure(), micro.ClientWithRegistry(r, time.Second*3),
		micro.ClientWithPickerBuilder("GROUP_ROUND_ROBIN", &round_robin.Builder{
			Filter: route.GroupFilterBuilder{}.Build(),
		}))
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	ctx = context.WithValue(ctx, "group", "A")
	defer cancel()
	cc, err := client.Dial(ctx, "user-service")
	require.NoError(t, err)
	uc := gen.NewUserServiceClient(cc)
	for i := 0; i < 10; i++ {
		resp, err := uc.GetById(ctx, &gen.GetByIdReq{Id: 123})
		require.NoError(t, err)
		t.Log(resp)
	}
}
