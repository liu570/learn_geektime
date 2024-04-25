package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"learn_geektime/micro/proto/gen"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	us := &Server{}
	server := grpc.NewServer()
	gen.RegisterUserServiceServer(server, us)
	l, err := net.Listen("tcp", ":8081")
	err = server.Serve(l)
	t.Log(err)
}

type Server struct {
	gen.UnimplementedUserServiceServer
}

func (s Server) GetById(ctx context.Context, req *gen.GetByIdReq) (*gen.GetByIdResp, error) {
	fmt.Println(req)
	return &gen.GetByIdResp{
		User: &gen.User{
			Id:   req.Id,
			Name: "hello grpc world",
		},
	}, nil
}
