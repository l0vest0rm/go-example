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
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (t *server) Submit(ctx context.Context, req *submit.SubmitRequest) (rsp *submit.SubmitResponse, err error) {
	fmt.Printf("%s,%s,%s\n", req.GetUrl(), req.GetTitle(), req.GetBody())
	return &submit.SubmitResponse{Code: 0, Message: ""}, nil
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
