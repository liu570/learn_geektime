package registry

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"
	"learn_geektime/micro"
	"learn_geektime/micro/proto/gen"
	"learn_geektime/micro/registry/etcd"
	"testing"
)

func TestServer(t *testing.T) {
	us := &UserServiceServer{}
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	require.NoError(t, err)
	r, err := etcd.NewRegistry(etcdClient)
	require.NoError(t, err)
	server, err := micro.NewServer("user-service", ":8081", micro.ServerWithRegistry(r))
	require.NoError(t, err)

	gen.RegisterUserServiceServer(server, us)
	err = server.Start()
	t.Log(err)
}

type UserServiceServer struct {
	gen.UnimplementedUserServiceServer
}

func (s UserServiceServer) GetById(ctx context.Context, req *gen.GetByIdReq) (*gen.GetByIdResp, error) {
	fmt.Println(req)
	return &gen.GetByIdResp{
		User: &gen.User{
			Id:   req.Id,
			Name: "hello grpc world",
		},
	}, nil
}
