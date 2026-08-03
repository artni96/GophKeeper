package user

import (
	"context"
	"fmt"
	"testing"

	usermodel "github.com/artni96/GophKeeper/internal/model/user"
	"github.com/artni96/GophKeeper/internal/repository/test_config"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func initUserRepo(db *sqlx.DB) (*Repository, error) {
	repo, err := NewRepository(db, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize test user repository: %w", err)
	}
	return repo, nil
}

func TestCreate(t *testing.T) {
	ctx := context.Background()
	tc, err := test_config.NewTestConfig(ctx)
	if err != nil {
		t.Fatal(err)
	}

	userRepo, err := initUserRepo(tc.DB)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		req     usermodel.UserCreate
		failure bool
	}{
		{
			name: "success",
			req: usermodel.UserCreate{
				Username:       "test",
				HashedPassword: "test",
			},
			failure: false,
		},
		{
			name: "failure - user already exists",
			req: usermodel.UserCreate{
				Username:       "test",
				HashedPassword: "test",
			},
			failure: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = userRepo.Create(ctx, tt.req)
			assert.Equal(t, tt.failure, err != nil)
			if tt.failure {
				assert.ErrorIs(t, err, ErrUserAlreadyExists)
			}
		})
	}
}

func TestGetByUsername(t *testing.T) {
	ctx := context.Background()
	tc, err := test_config.NewTestConfig(ctx)
	if err != nil {
		t.Fatal(err)
	}

	userRepo, err := initUserRepo(tc.DB)
	if err != nil {
		t.Fatal(err)
	}
	userData := usermodel.UserCreate{
		Username:       "test",
		HashedPassword: "test",
	}
	err = userRepo.Create(ctx, userData)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name     string
		failure  bool
		username string
	}{
		{
			name:     "success",
			failure:  false,
			username: userData.Username,
		},
		{
			name:     "failure - user does not exist",
			failure:  true,
			username: "test 2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbEntity, err := userRepo.GetByUsername(ctx, tt.username)
			assert.Equal(t, tt.failure, err != nil)
			if tt.failure {
				assert.ErrorIs(t, err, ErrUserNotFound)
			} else {
				assert.Equal(t, tt.username, dbEntity.Username)
				assert.Equal(t, tt.username, dbEntity.HashedPassword)
			}
		})
	}
}
