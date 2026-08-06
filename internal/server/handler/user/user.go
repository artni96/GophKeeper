package user

import (
	"context"
	"errors"

	pb "github.com/artni96/GophKeeper/api/proto/users"
	"github.com/artni96/GophKeeper/internal/server/model/user"
	userrepo "github.com/artni96/GophKeeper/internal/server/repository/user"
	userserv "github.com/artni96/GophKeeper/internal/server/service/user"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler represents the User handler instance.
type Handler struct {
	pb.UnimplementedUserServiceServer
	UserService userserv.ServiceI
	Logger      *zap.Logger
}

// NewHandler initializes and returns the User Handler instance.
func NewHandler(userService userserv.ServiceI, logger *zap.Logger) *Handler {
	return &Handler{
		UserService: userService,
		Logger:      logger,
	}
}

// CreateUser is a handler to create a new User entity.
func (h *Handler) CreateUser(ctx context.Context, req *pb.UserCreateRequest) (*pb.UserCreateResponse, error) {
	resp := &pb.UserCreateResponse{}
	userToCreate := user.UserCreateRequest{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}
	err := h.UserService.Create(ctx, userToCreate)
	if err != nil {
		if errors.Is(err, userrepo.ErrUserAlreadyExists) {
			h.Logger.Info(err.Error(), zap.String("username", req.GetUsername()))
			return nil, status.Errorf(codes.AlreadyExists, "user '%s' already exists", userToCreate.Username)
		}
		return nil, status.Errorf(codes.Internal, "internal grpc error")
	}
	resp.SetResult("created")
	return resp, nil
}

// Login is a handler to provide the user-authorization flow.
func (h *Handler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	resp := &pb.LoginResponse{}
	loginData := user.LoginRequest{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}
	token, err := h.UserService.Login(ctx, loginData)
	if err != nil {
		if errors.Is(err, userrepo.ErrUserNotFound) || errors.Is(err, userserv.ErrWrongUserOrPassword) {
			return nil, status.Error(codes.Unauthenticated, "login failed")
		}
		return nil, status.Errorf(codes.Internal, "internal grpc error")
	}
	resp.SetToken(token)
	return resp, nil
}
