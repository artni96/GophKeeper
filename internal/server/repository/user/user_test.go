package user

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	usermodel "github.com/artni96/GophKeeper/internal/server/model/user"
	"github.com/artni96/GophKeeper/tests"
	"github.com/stretchr/testify/assert"
)

func initUserRepo(db *sql.DB) (*Repository, error) {
	repo, err := NewRepository(db, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize test user repository: %w", err)
	}
	return repo, nil
}

type testConfig struct {
	ctx        context.Context
	testConfig *tests.TestDependencies
	userRepo   *Repository
}

func (c *testConfig) init(t *testing.T) {
	c.ctx = context.Background()
	newTC, err := tests.NewTestDependencies(c.ctx)
	if err != nil {
		log.Fatal(err)
	}
	c.testConfig = newTC

	userRepo, err := initUserRepo(c.testConfig.DB)
	if err != nil {
		t.Fatal(err)
	}
	c.userRepo = userRepo
}

func newTestConfig(t *testing.T) *testConfig {
	newTC := &testConfig{}
	newTC.init(t)
	return newTC
}

func TestCreate(t *testing.T) {
	tc := newTestConfig(t)
	tests := []struct {
		name     string
		userData usermodel.UserCreate
		keyData  usermodel.UserKeyCreate
		failure  bool
	}{
		{
			name: "success",
			userData: usermodel.UserCreate{
				Username:       "test",
				HashedPassword: "test",
			},
			keyData: usermodel.UserKeyCreate{
				EncryptedKey: []byte("test"),
				Salt:         []byte("test"),
				IsActive:     true,
				CreatedAt:    time.Now(),
			},
			failure: false,
		},
		{
			name: "failure - user already exists",
			userData: usermodel.UserCreate{
				Username:       "test",
				HashedPassword: "test",
			},
			keyData: usermodel.UserKeyCreate{
				EncryptedKey: []byte("test"),
				Salt:         []byte("test"),
				IsActive:     true,
				CreatedAt:    time.Now(),
			},
			failure: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tc.userRepo.Create(tc.ctx, tt.userData, tt.keyData)
			assert.Equal(t, tt.failure, err != nil)
			if tt.failure {
				assert.ErrorIs(t, err, ErrUserAlreadyExists)
			}
		})
	}
}

func TestGetByUsername(t *testing.T) {
	tc := newTestConfig(t)
	userData := usermodel.UserCreate{
		Username:       "test",
		HashedPassword: "test",
	}
	keyData := usermodel.UserKeyCreate{
		EncryptedKey: []byte("test"),
		Salt:         []byte("test"),
		IsActive:     true,
		CreatedAt:    time.Now(),
	}
	err := tc.userRepo.Create(tc.ctx, userData, keyData)
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
			dbEntity, err := tc.userRepo.GetByUsername(tc.ctx, tt.username)
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
