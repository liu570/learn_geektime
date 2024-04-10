package rpc

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestInitClientProxy(t *testing.T) {
	server := NewServer("tcp", ":8081")
	server.RegisterServer(&UserServiceServer{})
	go func() {
		err := server.Start()
		t.Log(err)
	}()
	time.Sleep(time.Second * 1)

	UsClient := &UserService{}
	err := InitClientProxy(":8081", UsClient)
	require.NoError(t, err)
	resp, err := UsClient.GetById(context.Background(), &GetByIdReq{Id: 123})
	require.NoError(t, err)
	require.Equal(t, &GetByIdResp{
		Msg: "hello liu",
	}, resp)

}
