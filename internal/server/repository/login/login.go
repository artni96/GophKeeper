package login

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/artni96/GophKeeper/internal/server/constants"
	"github.com/artni96/GophKeeper/internal/server/model/login"
	"github.com/artni96/GophKeeper/internal/server/repository/common"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const tableName = "logins"

// RepositoryI represents all method for the Login repository.
//
//go:generate mockgen -source=login.go -destination=mocks/login_repository_mock.go -package=mocks
type RepositoryI interface {
	Create(ctx context.Context, entity login.CreateLogin, notNullFields []string) (uint64, error)
	GetByNumber(ctx context.Context, number uint64, userID uuid.UUID) (*login.Login, error)
	Update(ctx context.Context, entity login.UpdateLogin, number uint64, userID uuid.UUID, fieldsToUpdate []string) error
	Delete(ctx context.Context, number uint64, userID uuid.UUID) error
	GetList(ctx context.Context, userID uuid.UUID) ([]login.Login, error)
}

// Repository implements the Login repository to manage login-related data through the database.
type Repository struct {
	db *gorm.DB
}

// NewRepository initializes and return the new Login repository instance.
func NewRepository(sqlDB *sql.DB, logger *zap.Logger) (*Repository, error) {
	db, err := common.InitDBConnByGORM(sqlDB, logger)
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

// Create creates a new Login entity in the database.
func (r *Repository) Create(ctx context.Context, entity login.CreateLogin, notNullFields []string) (uint64, error) {
	tx := r.db.WithContext(ctx).Begin()

	entityNumber, err := common.GetNextRecordNumber(tx, entity.UserID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to get the next user login number: %w", err)
	}
	entity.Number = entityNumber

	err = tx.Table(tableName).Select(notNullFields).Create(&entity).Error
	if err != nil {
		tx.Rollback()
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return 0, constants.ErrEntityAlreadyExists
		}
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to create login: %w", err)
	}
	return entityNumber, nil
}

// GetByNumber selects and returns the Login entity by its number from the database.
func (r *Repository) GetByNumber(ctx context.Context, number uint64, userID uuid.UUID) (*login.Login, error) {
	var entity login.Login

	db := r.db.WithContext(ctx)
	if err := db.Table(tableName).Where("number = ? AND user_id = ?", number, userID).First(&entity).Error; err != nil {
		return nil, fmt.Errorf("%w: number: %d", constants.ErrEntityNotFound, number)
	}
	return &entity, nil
}

// Update updates the Login entity in the database by its number (only for authors).
func (r *Repository) Update(ctx context.Context, entity login.UpdateLogin, number uint64, userID uuid.UUID, fieldsToUpdate []string) error {
	db := r.db.WithContext(ctx)
	result := db.Table(tableName).Select(fieldsToUpdate).Where("number = ? AND user_id = ?", number, userID).Updates(entity)
	if result.Error != nil {
		return fmt.Errorf("failed to update login with the number %d: %w", number, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: number: %d", constants.ErrEntityNotFound, number)
	}
	return nil
}

// Delete removes the Login entity from the database by its number (only for authors).
func (r *Repository) Delete(ctx context.Context, number uint64, userID uuid.UUID) error {
	db := r.db.WithContext(ctx)
	result := db.Table(tableName).Where("number = ? AND user_id = ?", number, userID).Delete(nil)
	if result.Error != nil {
		return fmt.Errorf("failed to delete login with the number %d: %w", number, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: number: %d", constants.ErrEntityNotFound, number)
	}
	return nil
}

// GetList returns the list of user's login-related entities from the database (only for authors).
func (r *Repository) GetList(ctx context.Context, userID uuid.UUID) ([]login.Login, error) {
	var dbEntities []login.Login
	result := r.db.WithContext(ctx).Table(tableName).Where("user_id = ?", userID).Find(&dbEntities)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get list of logins: %w", result.Error)
	}
	return dbEntities, nil
}
