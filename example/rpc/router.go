package main

import (
	"fly/example/rpc/example/controller"
	"fly/example/rpc/example/pb"

	"google.golang.org/grpc"
)

// Index
func Index(s *grpc.Server) {
	hc := &controller.HelloController{}
	pb.RegisterHelloWorldServiceServer(s, hc)
}
