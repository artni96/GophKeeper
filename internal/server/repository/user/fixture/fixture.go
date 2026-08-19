package fixture

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	usermodel "github.com/artni96/GophKeeper/internal/server/model/user"
	"github.com/artni96/GophKeeper/internal/server/repository/user"
	"github.com/google/uuid"
)

func CreateFirstUser(ctx context.Context, db *sql.DB, repo *user.Repository) (uuid.UUID, error) {
	if repo == nil {
		return uuid.Nil, fmt.Errorf("user repository is not initialized")
	}
	userData := usermodel.UserCreate{
		Username:       "user 1",
		HashedPassword: "user 1",
	}
	userKeys := usermodel.UserKeyCreate{
		EncryptedKey: []byte("user 1"),
		Salt:         []byte("user 1"),
		IsActive:     true,
		CreatedAt:    time.Now(),
	}
	err := repo.Create(ctx, userData, userKeys)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create first user: %w", err)
	}
	query := "SELECT id FROM users WHERE username = 'user 1'"
	var userID uuid.UUID
	err = db.QueryRow(query).Scan(&userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get first user: %w", err)
	}
	return userID, nil
}

func CreateSecondUser(ctx context.Context, db *sql.DB, repo *user.Repository) (uuid.UUID, error) {
	if repo == nil {
		return uuid.Nil, fmt.Errorf("user repository is not initialized")
	}
	userData := usermodel.UserCreate{
		Username:       "user 2",
		HashedPassword: "user 2",
	}
	userKeys := usermodel.UserKeyCreate{
		EncryptedKey: []byte("user 2"),
		Salt:         []byte("user 2"),
		IsActive:     true,
		CreatedAt:    time.Now(),
	}
	err := repo.Create(ctx, userData, userKeys)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create second user: %w", err)
	}
	query := "SELECT id FROM users WHERE username = 'user 2'"
	var userID uuid.UUID
	err = db.QueryRow(query).Scan(&userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get second user: %w", err)
	}
	return userID, nil
}
