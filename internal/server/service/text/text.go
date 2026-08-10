package text

import (
	"context"
	"time"

	"github.com/artni96/GophKeeper/internal/server/config"
	textmodel "github.com/artni96/GophKeeper/internal/server/model/text"
	"github.com/artni96/GophKeeper/internal/server/repository/text"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ServiceI represents the methods of the Login service.
type ServiceI interface {
	Create(ctx context.Context, data textmodel.CreateTextRequest) (uint64, error)
	GetByNumber(ctx context.Context, number uint64, userID uuid.UUID) (*textmodel.Text, error)
	Update(ctx context.Context, data textmodel.UpdateTextRequest, number uint64, userID uuid.UUID) error
	Delete(ctx context.Context, number uint64, userID uuid.UUID) error
	GetList(ctx context.Context, userID uuid.UUID) ([]textmodel.Text, error)
}

// Service implements the Text service to manage text-related data business logic.
type Service struct {
	repo   text.RepositoryI
	cfg    *config.Config
	logger *zap.Logger
}

// NewService initializes and returns the new Text service instance.
func NewService(cfg *config.Config, logger *zap.Logger, repo text.RepositoryI) *Service {
	return &Service{
		repo:   repo,
		cfg:    cfg,
		logger: logger,
	}
}

// Create creates a new Text entity by the repository.
func (s *Service) Create(ctx context.Context, data textmodel.CreateTextRequest) (uint64, error) {
	notNullFields := make([]string, 0, 3)
	entityToCreate := textmodel.CreateText{}

	if data.Text != nil {
		entityToCreate.Text = data.Text.Value
		notNullFields = append(notNullFields, "hashed_text")

		entityToCreate.Nonce = data.Nonce.Value
		notNullFields = append(notNullFields, "nonce")

		entityToCreate.KeyID = data.KeyID.Value
		notNullFields = append(notNullFields, "key_id")
	}

	if data.Title != nil {
		entityToCreate.Title = data.Title.Value
		notNullFields = append(notNullFields, "title")
	}

	if data.Description != nil {
		entityToCreate.Description = data.Description.Value
		notNullFields = append(notNullFields, "description")
	}

	entityToCreate.CreatedAt = time.Now()
	notNullFields = append(notNullFields, "created_at")

	entityToCreate.UserID = data.UserID
	notNullFields = append(notNullFields, "user_id")

	notNullFields = append(notNullFields, "number")

	number, err := s.repo.Create(ctx, entityToCreate, notNullFields)
	if err != nil {
		return 0, err
	}
	return number, nil
}

// GetByNumber selects and returns a Text entity by its number by the repository.
func (s *Service) GetByNumber(ctx context.Context, number uint64, userID uuid.UUID) (*textmodel.Text, error) {
	entity, err := s.repo.GetByNumber(ctx, number, userID)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// Update updates a Text entity by its number by the repository.
func (s *Service) Update(ctx context.Context, data textmodel.UpdateTextRequest, number uint64, userID uuid.UUID) error {
	notNullFields := make([]string, 0, 3)
	dataToUpdate := textmodel.UpdateText{}

	if data.Text != nil {
		dataToUpdate.Text = data.Text.Value
		notNullFields = append(notNullFields, "hashed_text")

		dataToUpdate.Nonce = data.Nonce.Value
		notNullFields = append(notNullFields, "nonce")

		dataToUpdate.KeyID = data.KeyID.Value
		notNullFields = append(notNullFields, "key_id")
	}

	if data.Title != nil {
		dataToUpdate.Title = data.Title.Value
		notNullFields = append(notNullFields, "title")
	}

	if data.Description != nil {
		dataToUpdate.Description = data.Description.Value
		notNullFields = append(notNullFields, "description")
	}

	dataToUpdate.UpdatedAt = time.Now()
	notNullFields = append(notNullFields, "updated_at")

	notNullFields = append(notNullFields, "number")

	err := s.repo.Update(ctx, dataToUpdate, number, userID, notNullFields)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a Text entity by its number by the repository.
func (s *Service) Delete(ctx context.Context, number uint64, userID uuid.UUID) error {
	err := s.repo.Delete(ctx, number, userID)
	if err != nil {
		return err
	}
	return nil
}

// GetList returns the list of user's text-related entities by the repository.
func (s *Service) GetList(ctx context.Context, userID uuid.UUID) ([]textmodel.Text, error) {
	result, err := s.repo.GetList(ctx, userID)
	if err != nil {
		return nil, err
	}
	return result, nil
}
