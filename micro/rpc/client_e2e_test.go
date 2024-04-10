package rpc

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestInitClientProxy(t *testing.T) {
	server := NewServer("tcp", ":8081")
	service := &UserServiceServer{}
	server.RegisterServer(service)
	go func() {
		err := server.Start()
		t.Log(err)
	}()
	time.Sleep(time.Second * 1)
	UsService := &UserService{}
	client, err := NewClient(":8081")
	require.NoError(t, err)
	err = client.InitService(UsService)
	require.NoError(t, err)

	testCases := []struct {
		name     string
		mock     func()
		msg      string
		wantErr  error
		wantResp *GetByIdResp
	}{
		{
			name: "no err",
			mock: func() {
				service.Msg = "hello world"
				service.Err = nil
			},
			wantResp: &GetByIdResp{
				Msg: "hello world",
			},
		},
		{
			name: "err",
			mock: func() {
				service.Err = errors.New("mock error")
				service.Msg = ""
			},
			wantResp: &GetByIdResp{},
			wantErr:  errors.New("mock error"),
		},
		{
			name: "both ",
			mock: func() {
				service.Err = errors.New("mock error")
				service.Msg = "hello world"
			},
			wantResp: &GetByIdResp{
				Msg: "hello world",
			},
			wantErr: errors.New("mock error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			resp, err := UsService.GetById(context.Background(), &GetByIdReq{Id: 123})
			require.Equal(t, tc.wantErr, err)
			require.Equal(t, tc.wantResp, resp)

		})
	}

}
