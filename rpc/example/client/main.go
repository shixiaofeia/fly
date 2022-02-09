package main

import (
	"context"
	"fly/rpc/example/pb"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("localhost:9999", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewHelloWorldServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Id: 1, Name: "fly"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("greet: %s", r.GetMessage())
}
