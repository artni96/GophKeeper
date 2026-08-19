package login

import (
	"context"
	"time"

	"github.com/artni96/GophKeeper/internal/server/config"
	commonmodel "github.com/artni96/GophKeeper/internal/server/model/common"
	loginmodel "github.com/artni96/GophKeeper/internal/server/model/login"
	commonrepo "github.com/artni96/GophKeeper/internal/server/repository/common"
	"github.com/artni96/GophKeeper/internal/server/service/common"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ServiceI represents
type ServiceI interface {
	Create(ctx context.Context, data loginmodel.CreateLoginRequest) (uint64, error)
	GetByNumber(ctx context.Context, number uint64, userID uuid.UUID) (*loginmodel.Login, error)
	Update(ctx context.Context, data loginmodel.UpdateLoginRequest, number uint64, userID uuid.UUID) error
	Delete(ctx context.Context, number uint64, userID uuid.UUID) error
	GetList(ctx context.Context, userID uuid.UUID) ([]loginmodel.Login, error)
}
type Service struct {
	repo   commonrepo.CRepositoryI[commonmodel.CreateEntityI, loginmodel.UpdateLogin, loginmodel.Login]
	cfg    *config.Config
	logger *zap.Logger
}

// NewService initializes and returns the new Login service instance.
func NewService(cfg *config.Config, logger *zap.Logger, repo commonrepo.CRepositoryI[commonmodel.CreateEntityI, loginmodel.UpdateLogin, loginmodel.Login]) *Service {
	return &Service{
		repo:   repo,
		cfg:    cfg,
		logger: logger,
	}
}

// Create creates a new Login entity by the repository.
func (s *Service) Create(ctx context.Context, data loginmodel.CreateLoginRequest) (uint64, error) {
	notNullFields := make([]string, 0, 9)
	entityToCreate := loginmodel.CreateLogin{}

	if data.Login.Value != nil {
		common.EncryptedFieldSetter(&data.Login, &entityToCreate.Login, "login", &notNullFields)
	}

	if data.Password.Value != nil {
		common.EncryptedFieldSetter(&data.Password, &entityToCreate.Password, "password", &notNullFields)
	}

	if data.URL != nil {
		entityToCreate.URL = data.URL.Value
		notNullFields = append(notNullFields, "url")
	}

	if data.Title != nil {
		entityToCreate.Title = data.Title.Value
		notNullFields = append(notNullFields, "title")
	}

	if data.Description != nil {
		entityToCreate.Description = data.Description.Value
		notNullFields = append(notNullFields, "description")
	}

	entityToCreate.UserID = data.UserID
	notNullFields = append(notNullFields, "user_id")
	entityToCreate.CreatedAt = time.Now()
	notNullFields = append(notNullFields, "created_at")
	notNullFields = append(notNullFields, "number")

	number, err := s.repo.Create(ctx, &entityToCreate, notNullFields)
	if err != nil {
		return number, err
	}
	return number, nil
}

// GetByNumber selects and returns a Login entity by its number by the repository.
func (s *Service) GetByNumber(ctx context.Context, number uint64, userID uuid.UUID) (*loginmodel.Login, error) {
	entity, err := s.repo.GetByNumber(ctx, number, userID)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// Update updates a Login entity by its number by the repository.
func (s *Service) Update(ctx context.Context, data loginmodel.UpdateLoginRequest, number uint64, userID uuid.UUID) error {
	notNullFields := make([]string, 0, 8)
	dataToUpdate := loginmodel.UpdateLogin{}

	if data.Login.Value != nil {
		common.EncryptedFieldSetter(&data.Login, &dataToUpdate.Login, "login", &notNullFields)
	}

	if data.Password.Value != nil {
		common.EncryptedFieldSetter(&data.Password, &dataToUpdate.Password, "password", &notNullFields)
	}

	if data.URL != nil {
		dataToUpdate.URL = data.URL.Value
		notNullFields = append(notNullFields, "url")
	}

	if data.Description != nil {
		dataToUpdate.Description = data.Description.Value
		notNullFields = append(notNullFields, "description")
	}

	if data.Title != nil {
		dataToUpdate.Title = data.Title.Value
		notNullFields = append(notNullFields, "title")
	}

	dataToUpdate.UpdatedAt = time.Now()
	notNullFields = append(notNullFields, "updated_at")

	err := s.repo.Update(ctx, dataToUpdate, number, userID, notNullFields)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a Login entity by its number by the repository.
func (s *Service) Delete(ctx context.Context, number uint64, userID uuid.UUID) error {
	err := s.repo.Delete(ctx, number, userID)
	if err != nil {
		return err
	}
	return nil
}

// GetList returns the list of user's logins by the repository.
func (s *Service) GetList(ctx context.Context, userID uuid.UUID) ([]loginmodel.Login, error) {
	result, err := s.repo.GetList(ctx, userID)
	if err != nil {
		return nil, err
	}
	return result, nil
}
