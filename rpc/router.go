package rpc

import (
	"fly/rpc/example/controller"
	"fly/rpc/example/pb"
	"google.golang.org/grpc"
)

func Index(s *grpc.Server) {
	hc := &controller.HelloController{}
	pb.RegisterHelloWorldServiceServer(s, hc)
}
