package registry

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/sync/errgroup"
	"learn_geektime/micro"
	"learn_geektime/micro/proto/gen"
	"learn_geektime/micro/registry/etcd"
	"testing"
	"time"
)

func TestServer(t *testing.T) {

	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	require.NoError(t, err)
	r, err := etcd.NewRegistry(etcdClient)
	require.NoError(t, err)
	group := "A"
	var eg errgroup.Group
	for i := 0; i < 5; i++ {
		time.Sleep(time.Millisecond * 200)
		group = string(byte('A' + i%2))
		server, err := micro.NewServer("user-service", micro.ServerWithRegistry(r), micro.ServerWithGroup(group))
		require.NoError(t, err)
		us := &UserServiceServer{
			group: group,
			addr:  fmt.Sprintf(":808%d", i+1),
		}
		fmt.Println(group, " ", fmt.Sprintf(":808%d", i+1))
		gen.RegisterUserServiceServer(server, us)
		eg.Go(func() error {
			return server.Start(fmt.Sprintf(":808%d", i+1))
		})
	}
	err = eg.Wait()
	t.Log(err)
}

type UserServiceServer struct {
	group string
	addr  string
	gen.UnimplementedUserServiceServer
}

func (s UserServiceServer) GetById(ctx context.Context, req *gen.GetByIdReq) (*gen.GetByIdResp, error) {
	return &gen.GetByIdResp{
		User: &gen.User{
			Id:   req.Id,
			Name: fmt.Sprintf("res:%v group:%s addr:%s time %20s", req, s.group, s.addr, time.Now().String()),
		},
	}, nil
}
