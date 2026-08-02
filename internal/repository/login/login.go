package login

import (
	"context"
	"fmt"

	"github.com/artni96/GophKeeper/internal/model/login"
	"github.com/artni96/GophKeeper/internal/repository/common"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const tableName = "logins"

var ErrLoginNotFound = gorm.ErrRecordNotFound

type RepositoryI interface {
	Create(ctx context.Context, entity login.CreateLogin) error
	GetByNumber(ctx context.Context, number int64, userID uuid.UUID) (*login.Login, error)
}

// Repository implements the Login repository to manage login-related data through the database.
type Repository struct {
	db *gorm.DB
}

// NewRepository initializes and return the new Login repository instance.
func NewRepository(sqlxDB *sqlx.DB, logger *zap.Logger) (*Repository, error) {
	db, err := common.InitDBConnByGORM(sqlxDB, logger)
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

// Create creates a new Login entity in the database.
func (r *Repository) Create(ctx context.Context, entity login.CreateLogin) error {
	tx := r.db.WithContext(ctx).Begin()

	entityNumber, err := common.GetNextRecordNumber(tx, entity.UserID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get the next user login number: %w", err)
	}
	entity.Number = entityNumber

	err = tx.Table(tableName).Create(&entity).Error
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create login: %w", err)
	}
	return nil
}

// GetByNumber selects and returns a Login entity by its number from the database.
func (r *Repository) GetByNumber(ctx context.Context, number int64, userID uuid.UUID) (*login.Login, error) {
	var entity login.Login

	db := r.db.WithContext(ctx)
	if err := db.Table(tableName).Where("number = ? AND user_id = ?", number, userID).First(&entity).Error; err != nil {
		return nil, fmt.Errorf("%w: number: %d", ErrLoginNotFound, number)
	}
	return &entity, nil
}
