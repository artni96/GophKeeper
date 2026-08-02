package login

import (
	"context"
	"time"

	"github.com/artni96/GophKeeper/internal/config"
	loginmodel "github.com/artni96/GophKeeper/internal/model/login"
	"github.com/artni96/GophKeeper/internal/repository/login"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ServiceI interface {
	Create(ctx context.Context, entity loginmodel.CreateLoginRequest) error
	GetByNumber(ctx context.Context, number int64, userID uuid.UUID) (*loginmodel.Login, error)
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
func (s *Service) Create(ctx context.Context, entity loginmodel.CreateLoginRequest) error {
	entityToCreate := loginmodel.CreateLogin{
		Login:       entity.Login,
		Password:    entity.Password,
		UserID:      entity.UserID,
		Title:       entity.Title,
		URL:         entity.URL,
		Description: entity.Description,
		CreatedAt:   time.Now(),
	}
	err := s.repo.Create(ctx, entityToCreate)
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
