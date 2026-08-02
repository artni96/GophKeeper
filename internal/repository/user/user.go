package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/artni96/GophKeeper/internal/model/user"
	"github.com/artni96/GophKeeper/internal/repository/common"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const tableName = "users"

var (
	PgErr                *pgconn.PgError
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

// RepositoryI represents the User repository interface.
type RepositoryI interface {
	Create(ctx context.Context, entity user.UserCreate) error
	GetByUsername(ctx context.Context, username string) (user.User, error)
}

// Repository implements the User repository to manage user-related data through the database.
type Repository struct {
	db *gorm.DB
}

// NewRepository initializes and return the new User repository instance.
func NewRepository(sqlxDB *sqlx.DB, logger *zap.Logger) (*Repository, error) {
	db, err := common.InitDBConnByGORM(sqlxDB, logger)
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

// Create creates a new user entity.
func (r *Repository) Create(ctx context.Context, entity user.UserCreate) error {
	err := r.db.WithContext(ctx).Table(tableName).Create(&entity).Error
	if err != nil {
		if errors.As(err, &PgErr) && PgErr.Code == "23505" {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to create user: %w", err)
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
	}
	return dbEntity, nil
}
