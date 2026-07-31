package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/artni96/GophKeeper/internal/model"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

const TableName = "users"

var (
	PgErr                *pgconn.PgError
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

// RepositoryI represents the User repository interface.
type RepositoryI interface {
	Create(ctx context.Context, entity model.UserCreate) error
	GetByUsername(ctx context.Context, username string) (model.User, error)
}

// Repository implements the User repository to manage user-related data through the database.
type Repository struct {
	db *gorm.DB
}

// NewRepository initializes and return the new User repository instance.
func NewRepository(sqlxDB *sqlx.DB, logger *zap.Logger) (*Repository, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlxDB,
	}), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		logger.Info("Failed to connect to database", zap.Error(err))
		return nil, err
	}
	return &Repository{db: db}, nil
}

// Create creates a new user entity.
func (r *Repository) Create(ctx context.Context, entity model.UserCreate) error {
	err := r.db.WithContext(ctx).Table(TableName).Create(&entity).Error
	if err != nil {
		if errors.As(err, &PgErr) && PgErr.Code == "23505" {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetByUsername returns a User instance by its username.
func (r *Repository) GetByUsername(ctx context.Context, username string) (model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, ErrUserNotFound
		}
	}
	return user, nil
}
