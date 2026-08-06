package user

import (
	"context"
	"fmt"
	"log"
	"testing"

	usermodel "github.com/artni96/GophKeeper/internal/server/model/user"
	"github.com/artni96/GophKeeper/tests"
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
			err := tc.userRepo.Create(tc.ctx, tt.req)
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
	err := tc.userRepo.Create(tc.ctx, userData)
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
