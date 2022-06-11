package repo

import (
	pb "github.com/project1/user-service/genproto"
)

//UserStorageI ...
type UserStorageI interface {
	CreateUser(*pb.User) (*pb.User, error)
	UpdateUser(*pb.User) (string, error)
	GetUserById(id string) (*pb.User, error)
	GetAllUser() ([]*pb.User, error)
	GetUserFromPost(userID string) (*pb.GetUserFromPostResponse, error)
	UserList(limit, page int64) ([]*pb.User, int64, error)
}
