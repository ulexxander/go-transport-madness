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

func ConvertUsers(mu []models.User) []*pb.User {
	var pu []*pb.User
	for _, u := range mu {
		pu = append(pu, ConvertUser(u))
	}
	return pu
}

func ConvertMessage(mm models.Message) *pb.Message {
	return &pb.Message{
		SenderUsername: mm.SenderUsername,
		Content:        mm.Content,
		CreatedAt:      mm.CreatedAt.String(),
	}
}

func ConvertMessages(mm []models.Message) []*pb.Message {
	var pm []*pb.Message
	for _, m := range mm {
		pm = append(pm, ConvertMessage(m))
	}
	return pm
}
