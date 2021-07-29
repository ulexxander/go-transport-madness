package grpc

import (
	"context"

	"github.com/ulexxander/transport-madness/transport/grpc/pb"
	"google.golang.org/grpc"
)

type Responder struct {
	Server *grpc.Server
}

func (r *Responder) Setup() {
	pb.RegisterTestServer(r.Server, &TestServer{})
}

type TestServer struct {
	pb.UnimplementedTestServer
}

// impl: ts *TestServer pb.TasksServer

func (ts *TestServer) SomeCall(ctx context.Context, _ *pb.Void) (*pb.TestReply, error) {
	return &pb.TestReply{
		Text: "hello",
	}, nil
}
