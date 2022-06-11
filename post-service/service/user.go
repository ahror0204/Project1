package service

import (
	"context"

	"github.com/jmoiron/sqlx"
	pb "github.com/project1/post-service/genproto"
	l "github.com/project1/post-service/pkg/logger"
	cl "github.com/project1/post-service/service/grpc_client"
	"github.com/project1/post-service/storage"

	// "github.com/project1/post-service/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//PostService ...
type PostService struct {
	Storage storage.IStorage
	Logger  l.Logger
	Client  cl.GrpcClientI
}

type PostI interface {
	CreatePost(ctx context.Context, req *pb.Post) (*pb.Post, error)
	GetPostById(ctx context.Context, req *pb.GetPostByIdRequest) (*pb.Post, error)
	GetAllUserPosts(ctx context.Context, req *pb.GetUserPostsrequest) (*pb.GetUserPosts, error)
	GetUserByPostId(ctx context.Context, req *pb.GetUserByPostIdRequest) (*pb.GetUserByPostIdResponse, error)
}

//NewPostService ...
func NewPostService(db *sqlx.DB, log l.Logger, client cl.GrpcClientI) *PostService {
	return &PostService{
		Storage: storage.NewStoragePg(db),
		Logger:  log,
		Client:  client,
	}
}

func (s *PostService) CreatePost(ctx context.Context, req *pb.Post) (*pb.Post, error) {
	
	user, err := s.Storage.Post().CreatePost(req)
	if err != nil {
		s.Logger.Error("failed while inserting post", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while inserting post")
	}

	return user, nil
}

func (s *PostService) GetPostById(ctx context.Context, req *pb.GetPostByIdRequest) (*pb.Post, error) {
	user, err := s.Storage.Post().GetPostById(req.UserId)
	if err != nil {
		s.Logger.Error("failed get post", l.Error(err))
		return nil, status.Error(codes.Internal, "failed get user")
	}

	return user, err
}

func (s *PostService) GetAllUserPosts(ctx context.Context, req *pb.GetUserPostsrequest) (*pb.GetUserPosts, error) {
	posts, err := s.Storage.Post().GetAllUserPosts(req.UserId)
	if err != nil {
		s.Logger.Error("failed get all user posts", l.Error(err))
		return nil, status.Error(codes.Internal, "failed get all user posts")
	}

	// user, err := s.client.UserServise().GetUserById(ctx, &pb.GetUserByIdRequest{
	// 	Id: req.UserId,
	// })

	// if err != nil {
	// 	s.Logger.Error("failed get a user by user_id in posts", l.Error(err))
	// 	return nil, status.Error(codes.Internal, "failed get a user by user_id in posts")
	// }

	// user.Posts = posts

	return &pb.GetUserPosts{
		Posts: posts,
	}, err
}

func (s *PostService) GetUserByPostId(ctx context.Context, req *pb.GetUserByPostIdRequest) (*pb.GetUserByPostIdResponse, error) {
	post, err := s.Storage.Post().GetUserByPostId(req.Post_Id)
	if err != nil {
		s.Logger.Error("failed get a post", l.Error(err))
		return nil, status.Error(codes.Internal, "failed get a post")
	}

	user, err := s.Client.UserServise().GetUserById(ctx, &pb.GetUserByIdRequest{
		Id: post.UserId,
	})

	if err != nil {
		s.Logger.Error("failed get a user by user_id in posts", l.Error(err))
		return nil, status.Error(codes.Internal, "failed get a user by user_id in posts")
	}

	post.UserFirstname = user.FirstName
	post.UserLastname = user.LastName

	return post, err
}
