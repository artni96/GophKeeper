package user

import (
	"context"
	"errors"
	"fmt"
	"sync"

	pb "github.com/artni96/GophKeeper/api/proto/users"
	"github.com/artni96/GophKeeper/internal/server/config"
	"github.com/artni96/GophKeeper/internal/server/model/user"
	userrepo "github.com/artni96/GophKeeper/internal/server/repository/user"
	userserv "github.com/artni96/GophKeeper/internal/server/service/user"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Handler represents the User handler instance.
type Handler struct {
	pb.UnimplementedUserServiceServer
	streams     map[uuid.UUID][]chan *pb.UpdateNotification
	mu          sync.RWMutex
	UserService userserv.ServiceI
	Logger      *zap.Logger
	cfg         *config.Config
}

// NewHandler initializes and returns the User Handler instance.
func NewHandler(
	userService userserv.ServiceI,
	logger *zap.Logger,
	streams map[uuid.UUID][]chan *pb.UpdateNotification,
	cfg *config.Config,
) *Handler {
	return &Handler{
		UserService: userService,
		Logger:      logger,
		streams:     streams,
		cfg:         cfg,
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
	userCredentials, err := h.UserService.Login(ctx, loginData)
	if err != nil {
		if errors.Is(err, userrepo.ErrUserNotFound) || errors.Is(err, userserv.ErrWrongUserOrPassword) {
			return nil, status.Error(codes.Unauthenticated, "login failed")
		}
		return nil, status.Errorf(codes.Internal, "internal grpc error")
	}
	
	resp.SetToken(userCredentials.Token)
	var keyspb []*pb.UserKey
	for _, key := range userCredentials.Keys {
		i := &pb.UserKey{}
		i.SetEncryptedKey(key.EncryptedKey)
		i.SetKeyId(key.KeyID)
		i.SetSalt(key.Salt)
		i.SetIsActive(key.IsActive)
		keyspb = append(keyspb, i)
	}
	resp.SetUserKeys(keyspb)
	return resp, nil
}

// SeekUpdates is streaming handler that notifies connected clients about user's data modification.
func (h *Handler) SeekUpdates(req *pb.SeekUpdateRequest, stream pb.UserService_SeekUpdatesServer) error {

	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Error(codes.Unauthenticated, "metadata is not provided")
	}
	if len(md["authorization"]) == 0 {
		h.Logger.Info("no authorization header in the request")
		return status.Error(codes.Unauthenticated, "authorization required")
	}

	token := md.Get("authorization")[0]
	if token == "" {
		h.Logger.Info("authorization header is empty")
		return status.Errorf(codes.Unauthenticated, "invalid token")
	}

	userID := userserv.GetUserIDFromJWT(token, h.cfg)
	if userID == uuid.Nil {
		h.Logger.Info("user is not authorized via jwt token")
		return status.Errorf(codes.Unauthenticated, "invalid token")

	}

	notificationChan := make(chan *pb.UpdateNotification, 10)

	h.mu.Lock()
	h.streams[userID] = append(h.streams[userID], notificationChan)
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		streams := h.streams[userID]
		for i, ch := range streams {
			if ch == notificationChan {
				h.streams[userID] = append(streams[:i], streams[i+1:]...)
				break
			}
		}
		if len(h.streams[userID]) == 0 {
			delete(h.streams, userID)
		}
		h.mu.Unlock()
		close(notificationChan)
	}()

	for {
		select {
		case update := <-notificationChan:
			fmt.Println("test on server side")
			if err := stream.Send(update); err != nil {
				return err
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}
