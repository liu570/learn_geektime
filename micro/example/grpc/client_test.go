package grpc

import (
	"context"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"learn_geektime/micro/proto/gen"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	cc, err := grpc.Dial("localhost:8081", grpc.WithInsecure())
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	client := gen.NewUserServiceClient(cc)
	resp, err := client.GetById(ctx, &gen.GetByIdReq{Id: 123})
	require.NoError(t, err)
	t.Fatal(resp)
}
