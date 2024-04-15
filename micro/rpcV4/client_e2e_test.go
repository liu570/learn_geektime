package rpc

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"learn_geektime/micro/proto/gen"
	"learn_geektime/micro/rpc/serialize/proto"
	"testing"
	"time"
)

func TestInitClientProxy(t *testing.T) {
	server := NewServer("tcp", ":8081")
	service := &UserServiceServer{}
	server.RegisterService(service)
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

func TestInitClientProto(t *testing.T) {
	server := NewServer("tcp", ":8081")
	service := &UserServiceServer{}
	server.RegisterService(service)
	server.RegisterSerializer(&proto.Serializer{})
	go func() {
		err := server.Start()
		t.Log(err)
	}()
	time.Sleep(time.Second * 1)
	UsService := &UserService{}
	client, err := NewClient(":8081", ClientWithSerializer(&proto.Serializer{}))
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
			resp, err := UsService.GetByIdProto(context.Background(), &gen.GetByIdReq{Id: 123})
			require.Equal(t, tc.wantErr, err)
			if resp != nil && resp.User != nil {
				require.Equal(t, tc.wantResp.Msg, resp.User.Name)
			}

		})
	}

}

func TestOneway(t *testing.T) {
	server := NewServer("tcp", ":8081")
	service := &UserServiceServer{}
	server.RegisterService(service)
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
			name: "oneway",
			mock: func() {
				service.Msg = "hello world"
				service.Err = errors.New("mock error")
			},
			wantResp: &GetByIdResp{},
			wantErr:  errors.New("micro:oneway调用,不处理任何响应结果"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			ctx := ContextWithOneway(context.Background())
			resp, err := UsService.GetById(ctx, &GetByIdReq{Id: 123})
			require.Equal(t, tc.wantErr, err)
			require.Equal(t, tc.wantResp, resp)

		})
	}

}
