package user

import (
	"context"
	"fmt"
	"time"

	userspb "github.com/artni96/GophKeeper/api/proto/users"
	"github.com/artni96/GophKeeper/internal/client/config"
	"github.com/artni96/GophKeeper/internal/client/model/common"
	"github.com/artni96/GophKeeper/internal/client/model/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Service implements the User client service to manage user-related business logic.
type Service struct {
	cfg              *config.Config
	client           userspb.UserServiceClient
	state            *config.State
	notificationChan chan common.Notification
	isBeingUpdated   bool
}

// NewService initializes and returns the new User service instance.
func NewService(
	cfg *config.Config,
	state *config.State,
	conn *grpc.ClientConn,
	notificationChan chan common.Notification,
	isBeingUpdated bool,
) *Service {
	return &Service{
		cfg:              cfg,
		state:            state,
		client:           userspb.NewUserServiceClient(conn),
		notificationChan: notificationChan,
		isBeingUpdated:   isBeingUpdated,
	}
}

// Login authorizes a login by its login and password and retrieves its JWT token.
func (s *Service) Login(ctx context.Context, userEntity user.LoginRequest) error {
	req := &userspb.LoginRequest{}
	req.SetUsername(userEntity.Login)
	req.SetPassword(userEntity.Password)

	token, err := s.client.Login(ctx, req)
	if err != nil {
		return err
	}
	s.cfg.Token = token.GetToken()
	s.state.IsOnline = true
	return nil
}

// Register make an RPC call on the server to create a new User entity.
func (s *Service) Register(ctx context.Context, userEntity user.RegistrationRequest) error {
	req := &userspb.UserCreateRequest{}
	req.SetUsername(userEntity.Login)
	req.SetPassword(userEntity.Password)

	_, err := s.client.CreateUser(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

// SeekUpdates seamlessly waits for user's data update notifications from the server.
func (s *Service) SeekUpdates(ctx context.Context) {
	lostConn := false
	md := metadata.Pairs("authorization", s.cfg.Token)
Reconnection:
	for {
		timeout := 10 * time.Second
		reqCtx := metadata.NewOutgoingContext(ctx, md)
		stream, err := s.client.SeekUpdates(reqCtx, &userspb.SeekUpdateRequest{})
		if err != nil {
			if !lostConn {
				lostConn = true
				s.state.IsOnline = false
				fmt.Println("You are being offline - failed to connect to the server.")
			}
			time.Sleep(timeout)
			continue Reconnection
		}
		if lostConn {
			lostConn = false
			s.state.IsOnline = true
			fmt.Println("You are being online - reconnected to the server.")
		}
		for {
			select {
			case <-ctx.Done():
				return
			default:
				update, err := stream.Recv()
				if err != nil {
					lostConn = true
					fmt.Println("You are being offline - lost connection to the server.")
					time.Sleep(timeout)
					continue Reconnection
				}
				fmt.Println("update received")

				notification := common.Notification{
					EntityType:   update.GetEntityType(),
					ActionType:   update.GetActionType(),
					EntityNumber: update.GetNumber(),
				}
				s.notificationChan <- notification
				s.isBeingUpdated = true
			}
		}
	}
}
