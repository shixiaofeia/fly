package controller

import (
	"context"
	"fly/example/rpc/example/pb"
	"fly/pkg/logging"
	"fmt"
)

type HelloController struct{}

// SayHello
func (*HelloController) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	msg := fmt.Sprintf("hello id: %d, name: %s", in.Id, in.Name)
	logging.Infof("SayHello: %s", msg)
	return &pb.HelloReply{Message: msg}, nil
}
