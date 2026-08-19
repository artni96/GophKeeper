package user

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"sync/atomic"
	"time"

	userspb "github.com/artni96/GophKeeper/api/proto/users"
	"github.com/artni96/GophKeeper/internal/client/config"
	"github.com/artni96/GophKeeper/internal/client/model/common"
	"github.com/artni96/GophKeeper/internal/client/model/user"
	"golang.org/x/crypto/pbkdf2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Service implements the User client service to manage user-related business logic.
type Service struct {
	cfg              *config.Config
	client           userspb.UserServiceClient
	state            *config.State
	notificationChan chan common.Notification
	isBeingUpdated   *atomic.Bool
}

// NewService initializes and returns the new User service instance.
func NewService(
	cfg *config.Config,
	state *config.State,
	conn *grpc.ClientConn,
	notificationChan chan common.Notification,
	isBeingUpdated *atomic.Bool,
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

	userCredentials, err := s.client.Login(ctx, req)
	if err != nil {
		return err
	}
	s.state.Token = userCredentials.GetToken()
	s.state.IsOnline = true
	s.state.Password = userEntity.Password

	userKeys := userCredentials.GetUserKeys()
	for _, key := range userKeys {
		err = s.handleKey(key)
		if err != nil {
			return err
		}
	}
	return nil
}

// handleKey extracts aesKey from *userspb.UserKey and saves it to the State.Keys.
func (s *Service) handleKey(key *userspb.UserKey) error {
	keyID := key.GetKeyId()
	encryptedKey := key.GetEncryptedKey()
	salt := key.GetSalt()
	derivedKey := pbkdf2.Key([]byte(s.state.Password), salt, 10000, 32, sha256.New)

	aesblock, err := aes.NewCipher(derivedKey)
	if err != nil {
		return err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return err
	}

	nonce := encryptedKey[:aesgcm.NonceSize()]
	ciphertext := encryptedKey[aesgcm.NonceSize():]

	aesKey, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	s.state.Keys[keyID] = aesKey
	if key.GetIsActive() {
		s.state.ActiveKey = aesKey
		s.state.ActiveKeyID = key.GetKeyId()
	}
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
	md := metadata.Pairs("authorization", s.state.Token)
Reconnection:
	for {
		timeout := 10 * time.Second
		reqCtx := metadata.NewOutgoingContext(ctx, md)
		stream, err := s.client.SeekUpdates(reqCtx, &userspb.SeekUpdateRequest{})
		if err != nil {
			if !lostConn {
				lostConn = true
				s.state.IsOnline = false
			}
			time.Sleep(timeout)
			continue Reconnection
		}
		if lostConn {
			lostConn = false
			s.state.IsOnline = true
			fmt.Println()
			fmt.Println("You are being online - reconnected to the server")
		}
		for {
			select {
			case <-ctx.Done():
				return
			default:
				update, err := stream.Recv()
				if err != nil {
					lostConn = true
					fmt.Println()
					fmt.Println("You are being offline - lost connection to the server")
					time.Sleep(timeout)
					continue Reconnection
				}
				//fmt.Println("\n\nupdate received")

				notification := common.Notification{
					EntityType:   update.GetEntityType(),
					ActionType:   update.GetActionType(),
					EntityNumber: update.GetNumber(),
				}
				//fmt.Println(notification)
				s.notificationChan <- notification
				s.isBeingUpdated.Store(true)
			}
		}
	}
}
