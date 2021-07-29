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
	return &pb.UserAllReply{Users: ConvertUsers(data)}, nil
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

// impl: ms *MessagesServer pb.MessagesServer

func (ms *MessagesServer) Pagination(ctx context.Context, req *pb.MessagesPaginationRequest) (*pb.MessagePaginationReply, error) {
	data, err := ms.MessagesService.MessagesPagination(services.MessagesPaginationInput{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	return &pb.MessagePaginationReply{
		Messages: ConvertMessages(data),
	}, nil
}

func (ms *MessagesServer) Create(ctx context.Context, req *pb.MessageCreateRequest) (*pb.MessageCreateReply, error) {
	data, err := ms.MessagesService.CreateMessage(services.MessageCreateInput{
		SenderUsername: req.SenderUsername,
		Content:        req.Content,
	})
	if err != nil {
		return nil, err
	}
	return &pb.MessageCreateReply{Message: ConvertMessage(*data)}, nil
}
