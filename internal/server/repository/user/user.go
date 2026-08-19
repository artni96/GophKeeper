package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/artni96/GophKeeper/internal/server/model/user"
	"github.com/artni96/GophKeeper/internal/server/repository/common"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	tableName = "users"
	keysTable = "user_keys"
)

var (
	PgErr                *pgconn.PgError
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

// RepositoryI represents the User repository interface.
type RepositoryI interface {
	Create(ctx context.Context, userEntity user.UserCreate, keyEntity user.UserKeyCreate) error
	GetByUsername(ctx context.Context, username string) (user.User, error)
	GetUserKeysList(ctx context.Context, userID uuid.UUID) ([]user.UserKey, error)
}

// Repository implements the User repository to manage user-related data through the database.
type Repository struct {
	db *gorm.DB
}

// NewRepository initializes and return the new User repository instance.
func NewRepository(sqlDB *sql.DB, logger *zap.Logger) (*Repository, error) {
	db, err := common.InitGORMDBConn(sqlDB, logger)
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

// Create creates a new user entity with its encrypted keys.
func (r *Repository) Create(ctx context.Context, entity user.UserCreate, keyEntity user.UserKeyCreate) error {
	tx := r.db.Begin().WithContext(ctx)

	err := tx.Table(tableName).Create(&entity).Error
	if err != nil {
		if errors.As(err, &PgErr) && PgErr.Code == "23505" {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to create new user: %w", err)
	}

	if entity.ID != nil {
		keyEntity.UserID = *entity.ID
	}

	err = tx.Table(keysTable).Create(&keyEntity).Error
	if err != nil {
		return fmt.Errorf("failed to create key: %w", err)
	}
	if tx.Commit().Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create new user: %w", tx.Commit().Error)
	}
	return nil
}

// GetByUsername returns a User instance by its username.
func (r *Repository) GetByUsername(ctx context.Context, username string) (user.User, error) {
	var dbEntity user.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&dbEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user.User{}, ErrUserNotFound
		}
		return user.User{}, fmt.Errorf("failed to get user by username: %w", err)
	}
	return dbEntity, nil
}

// GetUserKeysList returns all user keys by its userID.
func (r *Repository) GetUserKeysList(ctx context.Context, userID uuid.UUID) ([]user.UserKey, error) {
	var userKeys []user.UserKey
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&userKeys).Error
	if err != nil {
		return userKeys, fmt.Errorf("failed to get user keys: %w", err)
	}
	return userKeys, nil
}
