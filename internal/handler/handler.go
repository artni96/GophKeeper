package handler

import (
	"context"
	"errors"

	pb "github.com/artni96/GophKeeper/api/proto/users"
	"github.com/artni96/GophKeeper/internal/model"
	userrepo "github.com/artni96/GophKeeper/internal/repository/user"
	userserv "github.com/artni96/GophKeeper/internal/service/user"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	UserService userserv.ServiceI
	Logger      *zap.Logger
}

func NewUserHandler(userService userserv.ServiceI, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		UserService: userService,
		Logger:      logger,
	}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *pb.UserCreateRequest) (*pb.UserCreateResponse, error) {
	resp := &pb.UserCreateResponse{}
	userToCreate := model.UserCreateRequest{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}
	err := h.UserService.Create(ctx, userToCreate)
	if err != nil {
		if errors.As(err, &userrepo.ErrUserAlreadyExists) {
			h.Logger.Info(err.Error(), zap.String("username", req.GetUsername()))
			return nil, status.Errorf(codes.AlreadyExists, "user '%s' already exists", userToCreate.Username)
		}
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	resp.SetResult("created")
	return resp, nil
}

func (h *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	resp := &pb.LoginResponse{}
	loginData := model.LoginRequest{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}
	token, err := h.UserService.Login(ctx, loginData)
	if err != nil {
		if errors.As(err, &userrepo.ErrUserNotFound) || errors.As(err, &userserv.ErrWrongUserOrPassword) {
			return nil, status.Error(codes.Unauthenticated, "login failed")
		}
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	resp.SetToken(token)
	return resp, nil
}
