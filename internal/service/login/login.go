package login

import (
	"context"
	"fmt"
	"time"

	"github.com/artni96/GophKeeper/internal/config"
	"github.com/artni96/GophKeeper/internal/constants"
	loginmodel "github.com/artni96/GophKeeper/internal/model/login"
	"github.com/artni96/GophKeeper/internal/repository/login"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ServiceI interface {
	Create(ctx context.Context, data loginmodel.CreateLoginRequest) error
	GetByNumber(ctx context.Context, number int64, userID uuid.UUID) (*loginmodel.Login, error)
	Update(ctx context.Context, data loginmodel.UpdateLoginRequest, number int64, userID uuid.UUID) error
	Delete(ctx context.Context, number int64, userID uuid.UUID) error
	GetList(ctx context.Context, userID uuid.UUID) ([]loginmodel.GetListLoginResponse, error)
}
type Service struct {
	repo   login.RepositoryI
	cfg    *config.Config
	logger *zap.Logger
}

// NewService initializes and returns the new Login service instance.
func NewService(cfg *config.Config, logger *zap.Logger, repo login.RepositoryI) *Service {
	return &Service{
		repo:   repo,
		cfg:    cfg,
		logger: logger,
	}
}

// Create creates a new Login entity be the repository.
func (s *Service) Create(ctx context.Context, data loginmodel.CreateLoginRequest) error {
	if err := data.Validate(); err != nil {
		return fmt.Errorf("%w: %w", constants.ErrInvalidRequest, err)
	}
	creationData := loginmodel.CreateLogin{
		Login:       data.Login,
		Password:    data.Password,
		UserID:      data.UserID,
		Title:       data.Title,
		URL:         data.URL,
		Description: data.Description,
		CreatedAt:   time.Now(),
	}
	err := s.repo.Create(ctx, creationData)
	if err != nil {
		return err
	}
	return nil
}

// GetByNumber selects and returns a Login entity by its number by the repository.
func (s *Service) GetByNumber(ctx context.Context, number int64, userID uuid.UUID) (*loginmodel.Login, error) {
	entity, err := s.repo.GetByNumber(ctx, number, userID)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// Update updates a Login entity by its number by the repository.
func (s *Service) Update(ctx context.Context, data loginmodel.UpdateLoginRequest, number int64, userID uuid.UUID) error {
	dataToUpdate := loginmodel.UpdateLogin{
		Login:       data.Login,
		Password:    data.Password,
		Title:       data.Title,
		URL:         data.URL,
		Description: data.Description,
		UpdatedAt:   time.Now(),
	}

	err := s.repo.Update(ctx, dataToUpdate, number, userID)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a Login entity by its number by the repository.
func (s *Service) Delete(ctx context.Context, number int64, userID uuid.UUID) error {
	err := s.repo.Delete(ctx, number, userID)
	if err != nil {
		return err
	}
	return nil
}

// GetList returns the list of user's logins by the repository.
func (s *Service) GetList(ctx context.Context, userID uuid.UUID) ([]loginmodel.GetListLoginResponse, error) {
	result, err := s.repo.GetList(ctx, userID)
	if err != nil {
		return nil, err
	}
	return result, nil
}
