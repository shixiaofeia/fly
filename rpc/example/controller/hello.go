package controller

import (
	"context"
	"fly/pkg/logging"
	"fly/rpc/example/pb"
	"fmt"
)

type HelloController struct{}

// SayHello
func (*HelloController) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	msg := fmt.Sprintf("hello id: %d, name: %s", in.Id, in.Name)
	logging.Log.Info("SayHello", msg)
	return &pb.HelloReply{Message: msg}, nil
}
