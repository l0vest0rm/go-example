package main

import (
	"context"
	"fmt"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/l0vest0rm/go-example/submit-serv/submit"
)

const (
	addr = ":8080"
)

// server is used to implement helloworld.GreeterServer.
type SubmitServiceHandler struct{}

func NewSubmitServiceHandler() *SubmitServiceHandler {
	return &SubmitServiceHandler{}
}

func (t *SubmitServiceHandler) Ping(ctx context.Context) (err error) {
	return nil
}

// SayHello implements helloworld.GreeterServer
func (t *SubmitServiceHandler) Submit(ctx context.Context, req *submit.SubmitRequest) (r *submit.SubmitResponse, err error) {
	fmt.Printf("%s,%s,%s\n", req.GetURL(), req.GetTitle, req.GetBody())
	return &submit.SubmitResponse{Code: 0, Message: ""}, nil
}

func main() {
	transport, err := thrift.NewTServerSocket(addr)
	if err != nil {
		panic(err)
	}

	handler := NewSubmitServiceHandler()
	processor := submit.NewSubmitServiceProcessor(handler)
	transportFactory := thrift.NewTTransportFactory()
	protocolFactory := thrift.NewTJSONProtocolFactory()
	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)

	fmt.Println("Starting the simple server... on ", addr)
	server.Serve()
}
