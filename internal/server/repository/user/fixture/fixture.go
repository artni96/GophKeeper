package fixture

import (
	"context"
	"fmt"

	usermodel "github.com/artni96/GophKeeper/internal/server/model/user"
	"github.com/artni96/GophKeeper/internal/server/repository/user"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func CreateFirstUser(ctx context.Context, db *sqlx.DB, repo *user.Repository) (uuid.UUID, error) {
	if repo == nil {
		return uuid.Nil, fmt.Errorf("user repository is not initialized")
	}
	userData := usermodel.UserCreate{
		Username:       "user 1",
		HashedPassword: "user 1",
	}
	err := repo.Create(ctx, userData)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create first user: %w", err)
	}
	query := "SELECT id FROM users WHERE username = 'user 1'"
	var userID uuid.UUID
	err = db.DB.QueryRow(query).Scan(&userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get first user: %w", err)
	}
	return userID, nil
}

func CreateSecondUser(ctx context.Context, db *sqlx.DB, repo *user.Repository) (uuid.UUID, error) {
	if repo == nil {
		return uuid.Nil, fmt.Errorf("user repository is not initialized")
	}
	userData := usermodel.UserCreate{
		Username:       "user 2",
		HashedPassword: "user 2",
	}
	err := repo.Create(ctx, userData)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create second user: %w", err)
	}
	query := "SELECT id FROM users WHERE username = 'user 2'"
	var userID uuid.UUID
	err = db.DB.QueryRow(query).Scan(&userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get second user: %w", err)
	}
	return userID, nil
}
