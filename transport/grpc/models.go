package grpc

import (
	"github.com/ulexxander/transport-madness/models"
	"github.com/ulexxander/transport-madness/transport/grpc/pb"
)

func ConvertUser(mu models.User) *pb.User {
	return &pb.User{
		Username:  mu.Username,
		CreatedAt: mu.CreatedAt.String(),
	}
}

func ConvertUsers(mu []models.User) *pb.UserAllReply {
	var pu []*pb.User
	for _, u := range mu {
		pu = append(pu, ConvertUser(u))
	}
	return &pb.UserAllReply{Users: pu}
}
