package login

import (
	"context"
	"time"

	loginspb "github.com/artni96/GophKeeper/api/proto/logins"
	"github.com/artni96/GophKeeper/internal/client/constants"
	"github.com/artni96/GophKeeper/internal/client/model/common"
	"github.com/artni96/GophKeeper/internal/client/utils"

	"github.com/artni96/GophKeeper/internal/client/config"

	"github.com/artni96/GophKeeper/internal/client/model/login"
	loginrepo "github.com/artni96/GophKeeper/internal/client/repository/login"

	"google.golang.org/grpc"
)

// Service implements the Login client service to manage login-related data business logic.
type Service struct {
	cfg    *config.Config
	Client loginspb.LoginServiceClient
	repo   loginrepo.RepositoryI
	state  *config.State
}

// NewService initializes and returns the new Login service instance.
func NewService(cfg *config.Config, conn *grpc.ClientConn, repo loginrepo.RepositoryI, state *config.State) *Service {
	return &Service{
		cfg:    cfg,
		Client: loginspb.NewLoginServiceClient(conn),
		repo:   repo,
		state:  state,
	}
}

// Add adds a new Login entity to the repository.
func (s *Service) Add(ctx context.Context, entityNumber uint64) error {
	loginGetReq := &loginspb.LoginGetRequest{}
	loginGetReq.SetNumber(entityNumber)

	pbEntity, err := s.Client.GetLogin(ctx, loginGetReq)
	if err != nil {
		return err
	}

	createdAt, err := time.Parse(time.RFC3339, pbEntity.GetCreatedAt())
	if err != nil {
		return err
	}

	var updatedAt time.Time
	updatedAt, err = time.Parse(time.RFC3339, pbEntity.GetUpdatedAt())
	if err != nil {
		updatedAt = time.Time{}
	}
	nonce := pbEntity.GetNonce()
	keyID := pbEntity.GetKeyId()
	aesKey := s.state.Keys[keyID]

	decryptedLogin, err := utils.DecryptField(pbEntity.GetLogin(), aesKey, nonce)
	if err != nil {
		return err
	}

	decryptedPassword, err := utils.DecryptField(pbEntity.GetPassword(), aesKey, nonce)
	if err != nil {
		return err
	}

	entity := login.Login{
		Login:       string(decryptedLogin),
		Password:    string(decryptedPassword),
		Title:       pbEntity.GetTitle(),
		Number:      pbEntity.GetNumber(),
		URL:         pbEntity.GetUrl(),
		Description: pbEntity.GetDescription(),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
	s.repo.Add(entity)
	return nil
}

// AddBatch adds a list of new Login entities to the repository.
func (s *Service) AddBatch(ctx context.Context) error {
	loginListReq := &loginspb.LoginGetListRequest{}

	loginListResp, err := s.Client.GetListLogin(ctx, loginListReq)
	if err != nil {
		return err
	}

	pbEntities := loginListResp.GetLogins()

	var entities []login.Login
	for _, pbEntity := range pbEntities {
		nonce := pbEntity.GetNonce()
		keyID := pbEntity.GetKeyId()
		aesKey := s.state.Keys[keyID]

		decryptedLogin, err := utils.DecryptField(pbEntity.GetLogin(), aesKey, nonce)
		if err != nil {
			return err
		}

		decryptedPassword, err := utils.DecryptField(pbEntity.GetPassword(), aesKey, nonce)
		if err != nil {
			return err
		}
		createdAt, err := time.Parse(time.RFC3339, pbEntity.GetCreatedAt())
		if err != nil {
			return err
		}

		var updatedAt time.Time
		updatedAt, err = time.Parse(time.RFC3339, pbEntity.GetUpdatedAt())
		if err != nil {
			updatedAt = time.Time{}
		}

		entity := login.Login{
			Login:       string(decryptedLogin),
			Password:    string(decryptedPassword),
			Title:       pbEntity.GetTitle(),
			Number:      pbEntity.GetNumber(),
			URL:         pbEntity.GetUrl(),
			Description: pbEntity.GetDescription(),
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}
		entities = append(entities, entity)

	}
	s.repo.AddBatch(entities)
	return nil
}

// Get returns the Login entity by its number from the repository.
func (s *Service) Get(entityNumber uint64) (login.Login, error) {
	entity, err := s.repo.Get(entityNumber)
	if err != nil {
		return login.Login{}, err
	}
	return entity, nil
}

// GetList returns the list of Login entities from the repository.
func (s *Service) GetList() []common.Entity {
	return s.repo.GetList()
}

// Update updates the Login entity in the repository.
func (s *Service) Update(ctx context.Context, entityNumber uint64) error {
	req := &loginspb.LoginGetRequest{}
	req.SetNumber(entityNumber)
	pbEntity, err := s.Client.GetLogin(ctx, req)
	if err != nil {
		return constants.ErrEntityNotFound
	}

	nonce := pbEntity.GetNonce()
	keyID := pbEntity.GetKeyId()
	aesKey := s.state.Keys[keyID]

	decryptedLogin, err := utils.DecryptField(pbEntity.GetLogin(), aesKey, nonce)
	if err != nil {
		return err
	}

	decryptedPassword, err := utils.DecryptField(pbEntity.GetPassword(), aesKey, nonce)
	if err != nil {
		return err
	}

	createdAt, err := time.Parse(time.RFC3339, pbEntity.GetCreatedAt())
	if err != nil {
		return err
	}
	updatedAt, err := time.Parse(time.RFC3339, pbEntity.GetUpdatedAt())
	if err != nil {
		return err
	}
	updatedEntity := login.Login{
		Login:       string(decryptedLogin),
		Password:    string(decryptedPassword),
		Title:       pbEntity.GetTitle(),
		Number:      pbEntity.GetNumber(),
		URL:         pbEntity.GetUrl(),
		Description: pbEntity.GetDescription(),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
	err = s.repo.Update(updatedEntity)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes the Login entity by its number from the repository.
func (s *Service) Delete(entityNumber uint64) error {
	err := s.repo.Delete(entityNumber)
	if err != nil {
		return err
	}
	return nil
}
