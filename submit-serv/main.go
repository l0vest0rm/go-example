package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/l0vest0rm/go-example/submit-serv/submit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":8080"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (t *server) Post(ctx context.Context, req *submit.PostRequest) (rsp *submit.PostResponse, err error) {
	fmt.Printf("%s,%s,%s\n", req.GetUrl(), req.GetTitle(), req.GetBody())
	return &submit.PostResponse{Code: 0, Message: ""}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	submit.RegisterSubmitServiceServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
