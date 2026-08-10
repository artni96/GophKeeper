package card

import (
	"context"
	"time"

	"github.com/artni96/GophKeeper/internal/server/config"
	cardmodel "github.com/artni96/GophKeeper/internal/server/model/card"
	"github.com/artni96/GophKeeper/internal/server/repository/card"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ServiceI represents the methods of the Login service.
type ServiceI interface {
	Create(ctx context.Context, data cardmodel.CreateCardRequest) (uint64, error)
	GetByNumber(ctx context.Context, number uint64, userID uuid.UUID) (*cardmodel.Card, error)
	Update(ctx context.Context, data cardmodel.UpdateCardRequest, number uint64, userID uuid.UUID) error
	Delete(ctx context.Context, number uint64, userID uuid.UUID) error
	GetList(ctx context.Context, userID uuid.UUID) ([]cardmodel.Card, error)
}

// Service implements the Card service to manage card-related data business logic.
type Service struct {
	repo   card.RepositoryI
	cfg    *config.Config
	logger *zap.Logger
}

// NewService initializes and returns the new Card service instance.
func NewService(cfg *config.Config, logger *zap.Logger, repo card.RepositoryI) *Service {
	return &Service{
		repo:   repo,
		cfg:    cfg,
		logger: logger,
	}
}

// Create creates a new Card entity by the repository.
func (s *Service) Create(ctx context.Context, data cardmodel.CreateCardRequest) (uint64, error) {
	notNullFields := make([]string, 0, 22)
	entityToCreate := cardmodel.CreateCard{}

	if data.PAN != nil {
		entityToCreate.PAN = data.PAN.Value
		notNullFields = append(notNullFields, "hashed_pan")

		entityToCreate.PanNonce = data.PANNonce.Value
		notNullFields = append(notNullFields, "pan_nonce")

		entityToCreate.PanKeyID = data.PANKeyID.Value
		notNullFields = append(notNullFields, "pan_key_id")
	}

	if data.Holder != nil {
		entityToCreate.Holder = data.Holder.Value
		notNullFields = append(notNullFields, "hashed_holder")

		entityToCreate.HolderNonce = data.HolderNonce.Value
		notNullFields = append(notNullFields, "holder_nonce")

		entityToCreate.HolderKeyID = data.HolderKeyID.Value
		notNullFields = append(notNullFields, "holder_key_id")
	}

	if data.ExpiryDate != nil {
		entityToCreate.ExpiryDate = data.ExpiryDate.Value
		notNullFields = append(notNullFields, "hashed_expiry_date")

		entityToCreate.ExpiryDateNonce = data.ExpiryDateNonce.Value
		notNullFields = append(notNullFields, "expiry_date_nonce")

		entityToCreate.ExpiryDateKeyID = data.ExpiryDateKeyID.Value
		notNullFields = append(notNullFields, "expiry_date_key_id")
	}

	if data.CVV != nil {
		entityToCreate.CVV = data.CVV.Value
		notNullFields = append(notNullFields, "hashed_cvv")

		entityToCreate.CVVNonce = data.CVVNonce.Value
		notNullFields = append(notNullFields, "cvv_nonce")

		entityToCreate.CVVKeyID = data.CVVKeyID.Value
		notNullFields = append(notNullFields, "cvv_key_id")
	}

	if data.PIN != nil {
		entityToCreate.PIN = data.PIN.Value
		notNullFields = append(notNullFields, "hashed_pin")

		entityToCreate.PINNonce = data.PINNonce.Value
		notNullFields = append(notNullFields, "pin_nonce")

		entityToCreate.PINKeyID = data.PINKeyID.Value
		notNullFields = append(notNullFields, "pin_key_id")
	}

	if data.Bank != nil {
		entityToCreate.Bank = data.Bank.Value
		notNullFields = append(notNullFields, "bank")
	}

	if data.Brand != nil {
		entityToCreate.Brand = data.Brand.Value
		notNullFields = append(notNullFields, "brand")
	}

	entityToCreate.UserID = data.UserID
	notNullFields = append(notNullFields, "user_id")

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
	notNullFields = append(notNullFields, "number")

	number, err := s.repo.Create(ctx, entityToCreate, notNullFields)
	if err != nil {
		return 0, err
	}
	return number, nil
}

// GetByNumber selects and returns a Card entity by its number by the repository.
func (s *Service) GetByNumber(ctx context.Context, number uint64, userID uuid.UUID) (*cardmodel.Card, error) {
	entity, err := s.repo.GetByNumber(ctx, number, userID)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// Update updates a Card entity by its number by the repository.
func (s *Service) Update(ctx context.Context, data cardmodel.UpdateCardRequest, number uint64, userID uuid.UUID) error {
	notNullFields := make([]string, 0, 20)
	dataToUpdate := cardmodel.UpdateCard{}

	if data.PAN != nil {
		dataToUpdate.PAN = data.PAN.Value
		notNullFields = append(notNullFields, "hashed_pan")

		dataToUpdate.PanNonce = data.PANNonce.Value
		notNullFields = append(notNullFields, "pan_nonce")

		dataToUpdate.PanKeyID = data.PANKeyID.Value
		notNullFields = append(notNullFields, "pan_key_id")
	}

	if data.Holder != nil {
		dataToUpdate.Holder = data.Holder.Value
		notNullFields = append(notNullFields, "hashed_holder")

		dataToUpdate.HolderNonce = data.HolderNonce.Value
		notNullFields = append(notNullFields, "holder_nonce")

		dataToUpdate.HolderKeyID = data.HolderKeyID.Value
		notNullFields = append(notNullFields, "holder_key_id")
	}

	if data.ExpiryDate != nil {
		dataToUpdate.ExpiryDate = data.ExpiryDate.Value
		notNullFields = append(notNullFields, "hashed_expiry_date")

		dataToUpdate.ExpiryDateNonce = data.ExpiryDateNonce.Value
		notNullFields = append(notNullFields, "expiry_date_nonce")

		dataToUpdate.ExpiryDateKeyID = data.ExpiryDateKeyID.Value
		notNullFields = append(notNullFields, "expiry_date_key_id")
	}

	if data.CVV != nil {
		dataToUpdate.CVV = data.CVV.Value
		notNullFields = append(notNullFields, "hashed_cvv")

		dataToUpdate.CVVNonce = data.CVVNonce.Value
		notNullFields = append(notNullFields, "cvv_nonce")

		dataToUpdate.CVVKeyID = data.CVVKeyID.Value
		notNullFields = append(notNullFields, "cvv_key_id")
	}

	if data.PIN != nil {
		dataToUpdate.PIN = data.PIN.Value
		notNullFields = append(notNullFields, "hashed_pin")

		dataToUpdate.PINNonce = data.PINNonce.Value
		notNullFields = append(notNullFields, "pin_nonce")

		dataToUpdate.PINKeyID = data.PINKeyID.Value
		notNullFields = append(notNullFields, "pin_key_id")
	}

	if data.Bank != nil {
		dataToUpdate.Bank = data.Bank.Value
		notNullFields = append(notNullFields, "bank")
	}

	if data.Brand != nil {
		dataToUpdate.Brand = data.Brand.Value
		notNullFields = append(notNullFields, "brand")
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

	err := s.repo.Update(ctx, dataToUpdate, number, userID, notNullFields)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a Card entity by its number by the repository.
func (s *Service) Delete(ctx context.Context, number uint64, userID uuid.UUID) error {
	err := s.repo.Delete(ctx, number, userID)
	if err != nil {
		return err
	}
	return nil
}

// GetList returns the list of user's card-related entities by the repository.
func (s *Service) GetList(ctx context.Context, userID uuid.UUID) ([]cardmodel.Card, error) {
	result, err := s.repo.GetList(ctx, userID)
	if err != nil {
		return nil, err
	}
	return result, nil
}
