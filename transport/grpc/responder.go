package grpc

import (
	"context"

	"github.com/ulexxander/transport-madness/services"
	"github.com/ulexxander/transport-madness/transport/grpc/pb"
	"google.golang.org/grpc"
)

type Responder struct {
	Server *grpc.Server
	*UsersServer
	*MessagesServer
}

func (r *Responder) Setup() {
	pb.RegisterUsersServer(r.Server, r.UsersServer)
	pb.RegisterMessagesServer(r.Server, r.MessagesServer)
}

type UsersServer struct {
	pb.UnimplementedUsersServer
	UsersService *services.UsersService
}

// impl: us *UsersServer pb.UsersServer

func (us *UsersServer) All(ctx context.Context, _ *pb.Void) (*pb.UserAllReply, error) {
	data := us.UsersService.UsersAll()
	return ConvertUsers(data), nil
}

func (us *UsersServer) Create(ctx context.Context, req *pb.UserCreateRequest) (*pb.UserCreateReply, error) {
	data, err := us.UsersService.CreateUser(services.UserCreateInput{
		Username: req.Username,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UserCreateReply{User: ConvertUser(*data)}, nil
}

type MessagesServer struct {
	pb.UnimplementedMessagesServer
	MessagesService *services.MessagesService
}
