package rpc

import (
	"context"
	"log"
)

type UserService struct {
	GetById func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
}

func (u UserService) Name() string {
	return "user-service"
}

type GetByIdReq struct {
	Id int
}

type GetByIdResp struct {
	Msg string
}

type UserServiceServer struct {
	Msg string
	Err error
}

func (u *UserServiceServer) Name() string {
	return "user-service"
}

func (u *UserServiceServer) GetById(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error) {
	log.Println("UserServiceServer", req)
	return &GetByIdResp{
		Msg: u.Msg,
	}, u.Err
}
