package service

import (
	"context"

	// "os/user"

	"github.com/jmoiron/sqlx"
	pb "github.com/project1/user-service/genproto"
	l "github.com/project1/user-service/pkg/logger"
	cl "github.com/project1/user-service/service/grpc_client"
	"github.com/project1/user-service/storage"

	// "golang.org/x/tools/go/analysis/passes/nilfunc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"golang.org/x/crypto/bcrypt"
)

//UserService ...
type UserService struct {
	storage storage.IStorage
	logger  l.Logger
	client  cl.GrpcClientI
}

//NewUserService ...
func NewUserService(db *sqlx.DB, log l.Logger, client cl.GrpcClientI) *UserService {
	return &UserService{
		storage: storage.NewStoragePg(db),
		logger:  log,
		client:  client,
	}
}

func(s *UserService) LogIn(ctx context.Context, req *pb.LogInRequest) (*pb.User, error) {
	user, err := s.storage.User().LogIn(req)
	if err != nil {
		s.logger.Error("pasword", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while getting by Id user")
	}

	posts, err := s.client.PostService().GetAllUserPosts(ctx, &pb.GetUserPostsrequest{UserId: user.Id})
	
	if err != nil {
		s.logger.Error("failed while getting user posts", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while getting user posts")
	}
	user.Posts = posts.Posts
	
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		s.logger.Error("password comparison error", l.Error(err))
		return nil, status.Error(codes.Internal, "password comparison error")
	}

	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.User) (*pb.User, error) {

	user, err := s.storage.User().CreateUser(req)
	if err != nil {
		s.logger.Error("failed while creating user", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while creating user")
	}
	if req.Posts != nil {
		for _, post := range req.Posts {
			post.UserId = user.Id
			createdPosts, err := s.client.PostService().CreatePost(context.Background(), post)
			if err != nil {
				s.logger.Error("failed while inserting user post", l.Error(err))
				return nil, status.Error(codes.Internal, "failed while inserting user post")
			}
			user.Posts = append(user.Posts, createdPosts)
		}
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.User) (*pb.UpdateUserResponse, error) {
	id, err := s.storage.User().UpdateUser(req)
	if err != nil {
		s.logger.Error("failed while updating user", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while updating user")
	}
	return &pb.UpdateUserResponse{
		Id: id,
	}, nil
}

func (s *UserService) GetUserById(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.User, error) {
	user, err := s.storage.User().GetUserById(req.Id)
	if err != nil {
		s.logger.Error("failed while getting by Id user", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while getting by Id user")
	}

	posts, err := s.client.PostService().GetAllUserPosts(ctx, &pb.GetUserPostsrequest{UserId: req.Id})

	if err != nil {
		s.logger.Error("failed while getting user posts", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while getting user posts")
	}

	user.Posts = posts.Posts
	return user, nil
}

func (s *UserService) GetAllUser(ctx context.Context, req *pb.Empty) (*pb.GetAllResponse, error) {
	users, err := s.storage.User().GetAllUser()
	if err != nil {
		s.logger.Error("failed while getting All users", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while getting All users")
	}

	for _, user := range users {
		posts, err := s.client.PostService().GetAllUserPosts(
			ctx,
			&pb.GetUserPostsrequest{
				UserId: user.Id,
			},
		)
		if err != nil {
			s.logger.Error("failed while getting user posts", l.Error(err))
			return nil, status.Error(codes.Internal, "failed while getting user posts")
		}

		user.Posts = posts.Posts
	}

	return &pb.GetAllResponse{
		Users: users,
	}, err
}

func (s *UserService) GetUserFromPost(ctx context.Context, req *pb.GetUserFromPostRequest) (*pb.GetUserFromPostResponse, error) {
	user, err := s.storage.User().GetUserFromPost(req.UserId)
	if err != nil {
		s.logger.Error("failed while getting a user", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while getting a user")
	}

	return user, nil
}

func (s *UserService) GetAllUserPosts(ctx context.Context, req *pb.GetUserPostsrequest) (*pb.GetUserPosts, error) {
	res, err := s.client.PostService().GetAllUserPosts(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (s *UserService) UserList(ctx context.Context, req *pb.UserListRequest) (*pb.UserListResponse, error) {
	users, count, err := s.storage.User().UserList(req.Limit, req.Page)

	if err != nil {
		s.logger.Error("failed while getting list of users", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while getting list of users")
	}

	for _, user := range users {
		posts, err := s.client.PostService().GetAllUserPosts(
			ctx,
			&pb.GetUserPostsrequest{
				UserId: user.Id,
			},
		)

		if err != nil {
			s.logger.Error("failed while getting list of user postd", l.Error(err))
			return nil, status.Error(codes.Internal, "failed while getting list of user posts")
		}

		user.Posts = posts.Posts
	}

	return &pb.UserListResponse{
		User:  users,
		Count: count,
	}, nil
}

func (s *UserService) CheckField(ctx context.Context, req *pb.UserCheckRequest) (*pb.UserCheckResponse, error) {
	
	bl, err := s.storage.User().CheckFeild(req.Field, req.Value)
	
	if err != nil {
		s.logger.Error("CheckFeild FUNC ERROR", l.Error(err))
		return nil, status.Error(codes.Internal, "CheckFeild FUNC ERROR")
	}

	return &pb.UserCheckResponse{Response: bl}, nil
}
