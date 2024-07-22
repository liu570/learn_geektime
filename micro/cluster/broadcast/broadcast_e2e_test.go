package broadcast

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"learn_geektime/micro"
	"learn_geektime/micro/proto/gen"
	"learn_geektime/micro/registry/etcd"
	"testing"
	"time"
)

func TestUseBoardCast(t *testing.T) {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	require.NoError(t, err)
	r, err := etcd.NewRegistry(etcdClient)
	require.NoError(t, err)
	go func() {
		var eg errgroup.Group
		for i := 0; i < 5; i++ {
			time.Sleep(time.Millisecond * 200)
			server, err := micro.NewServer("user-service", micro.ServerWithRegistry(r))
			require.NoError(t, err)
			us := &UserServiceServer{
				addr: fmt.Sprintf(":808%d", i+1),
			}
			gen.RegisterUserServiceServer(server, us)
			eg.Go(func() error {
				return server.Start(fmt.Sprintf(":808%d", i+1))
			})
		}
		err = eg.Wait()
		t.Log(err)
	}()
	time.Sleep(time.Millisecond * 2000)

	client, err := micro.NewClient(micro.ClientWithInsecure(), micro.ClientWithRegistry(r, time.Second*3))
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	ctx = context.WithValue(ctx, "group", "A")
	defer cancel()
	bd := NewClusterBuilder(r, "user-service")
	cc, err := client.Dial(ctx, "user-service", grpc.WithUnaryInterceptor(bd.buildUnaryInterceptor()))
	require.NoError(t, err)
	uc := gen.NewUserServiceClient(cc)
	resp, err := uc.GetById(ctx, &gen.GetByIdReq{Id: 123})
	require.NoError(t, err)
	t.Log(resp)
}

type UserServiceServer struct {
	addr string
	gen.UnimplementedUserServiceServer
}

func (s UserServiceServer) GetById(ctx context.Context, req *gen.GetByIdReq) (*gen.GetByIdResp, error) {
	return &gen.GetByIdResp{
		User: &gen.User{
			Id:   req.Id,
			Name: fmt.Sprintf("hello grpc world addr:%s", s.addr),
		},
	}, nil
}
