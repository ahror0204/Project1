package repo

import (
	pb "github.com/project1/post-service/genproto"
)

//PostStorageI ...
type PostStorageI interface {
	CreatePost(*pb.Post) (*pb.Post, error)
	GetPostById(id string) (*pb.Post, error)
	GetAllUserPosts(userID string) ([]*pb.Post, error)
	GetUserByPostId(postID string) (*pb.GetUserByPostIdResponse, error)
	CreatePostUser(user *pb.User) error
}
