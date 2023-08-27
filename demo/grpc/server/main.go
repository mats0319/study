package main

import (
	"context"
	"fmt"
	"github.com/mats9693/study/demo/grpc/proto/impl"
	"google.golang.org/grpc"
	"log"
	"net"
)

type serverImpl struct {
	rpc_impl.UnimplementedIUserServer
}

func (s *serverImpl) Login(_ context.Context, req *rpc_impl.LoginReq) (*rpc_impl.LoginRes, error) {
	fmt.Println("> Node: get login req.", req.UserName, req.Password)

	return &rpc_impl.LoginRes{IsSuccess: true}, nil
}

func main() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":9693"))
	if err != nil {
		log.Fatalln("listen tcp failed, error:", err)
	}

	s := grpc.NewServer()
	rpc_impl.RegisterIUserServer(s, &serverImpl{})

	err = s.Serve(listen)
	if err != nil {
		log.Fatalln("serve failed, error:", err)
	}
}
