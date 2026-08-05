package text

import (
	"context"
	"errors"
	"fmt"

	"github.com/artni96/GophKeeper/internal/constants"
	commonmodel "github.com/artni96/GophKeeper/internal/model/common"
	"github.com/artni96/GophKeeper/internal/model/text"
	"github.com/artni96/GophKeeper/internal/repository/common"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const tableName = "texts"

// RepositoryI represents the methods of the Text repository.
type RepositoryI interface {
	Create(ctx context.Context, entity text.CreateText, notNullFields []string) (uint64, error)
	GetByNumber(ctx context.Context, number uint64, userID uuid.UUID) (*text.Text, error)
	Update(ctx context.Context, entity text.UpdateText, number uint64, userID uuid.UUID, notNullFields []string) error
	Delete(ctx context.Context, number uint64, userID uuid.UUID) error
	GetList(ctx context.Context, userID uuid.UUID) ([]commonmodel.GetListEntityResponse, error)
}

// Repository implements the Text repository to manage text-related data through the database.
type Repository struct {
	db *gorm.DB
}

// NewRepository initializes and return the new Text repository instance.
func NewRepository(sqlxDB *sqlx.DB, logger *zap.Logger) (*Repository, error) {
	db, err := common.InitDBConnByGORM(sqlxDB, logger)
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

// Create creates a new Text entity in the database.
func (r *Repository) Create(ctx context.Context, entity text.CreateText, notNullFields []string) (uint64, error) {
	tx := r.db.WithContext(ctx).Begin()

	entityNumber, err := common.GetNextRecordNumber(tx, entity.UserID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to get the next user entity number: %w", err)
	}
	entity.Number = entityNumber

	err = tx.Table(tableName).Select(notNullFields).Create(&entity).Error
	if err != nil {
		tx.Rollback()
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23502" {
				return 0, fmt.Errorf("%w: %s", constants.ErrRequiredField, pgErr.ColumnName)
			} else if pgErr.Code == "23505" {
				return 0, constants.ErrEntityAlreadyExists
			}
		}
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to create the text record: %w", err)
	}
	return entityNumber, nil
}

// GetByNumber selects and returns the Text entity by its number from the database.
func (r *Repository) GetByNumber(ctx context.Context, number uint64, userID uuid.UUID) (*text.Text, error) {
	var entity text.Text

	db := r.db.WithContext(ctx)
	if err := db.Table(tableName).Where(
		"number = ? AND user_id = ?", number, userID).First(&entity).Error; err != nil {
		return nil, fmt.Errorf("%w: number: %d", constants.ErrEntityNotFound, number)
	}
	return &entity, nil
}

// Update updates the Text entity in the database by its number (only for authors).
func (r *Repository) Update(
	ctx context.Context, entity text.UpdateText, number uint64, userID uuid.UUID, notNullFields []string) error {
	db := r.db.WithContext(ctx)
	result := db.Table(tableName).Select(notNullFields).Where(
		"number = ? AND user_id = ?", number, userID).Updates(entity)
	if result.Error != nil {
		return fmt.Errorf("failed to update the text record with the number %d: %w", number, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: number: %d", constants.ErrEntityNotFound, number)
	}
	return nil
}

// Delete removes the Text entity from the database by its number (only for authors).
func (r *Repository) Delete(ctx context.Context, number uint64, userID uuid.UUID) error {
	db := r.db.WithContext(ctx)
	result := db.Table(tableName).Where("number = ? AND user_id = ?", number, userID).Delete(nil)
	if result.Error != nil {
		return fmt.Errorf("failed to delete the text record with the number %d: %w", number, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: number: %d", constants.ErrEntityNotFound, number)
	}
	return nil
}

// GetList returns the list of user's text-related entities from the database (only for authors).
func (r *Repository) GetList(ctx context.Context, userID uuid.UUID) ([]commonmodel.GetListEntityResponse, error) {
	var dbEntities []commonmodel.GetListEntityResponse
	result := r.db.WithContext(ctx).Table(tableName).Where("user_id = ?", userID).Find(&dbEntities)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get list of logins: %w", result.Error)
	}
	return dbEntities, nil
}
