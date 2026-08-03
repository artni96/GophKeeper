package login

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/artni96/GophKeeper/internal/constants"
	"github.com/artni96/GophKeeper/internal/model/login"
	"github.com/artni96/GophKeeper/internal/repository/test_config"
	"github.com/artni96/GophKeeper/internal/repository/user"
	userfixture "github.com/artni96/GophKeeper/internal/repository/user/fixture"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func initUserRepo(db *sqlx.DB) (*user.Repository, error) {
	repo, err := user.NewRepository(db, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize test user repository: %w", err)
	}
	return repo, nil
}

func initLoginRepo(db *sqlx.DB) (*Repository, error) {
	repo, err := NewRepository(db, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize test login repository: %w", err)
	}
	return repo, nil
}

func loginFixture(ctx context.Context, loginRepo *Repository, entityData login.CreateLogin, userID uuid.UUID) {
	testEntity := login.CreateLogin{
		UserID:    entityData.UserID,
		Login:     entityData.Login,
		Password:  entityData.Password,
		Title:     entityData.Title,
		CreatedAt: entityData.CreatedAt,
	}
	err := loginRepo.Create(ctx, testEntity)
	if err != nil {
		log.Fatal(err)
	}
}

type testConfig struct {
	ctx        context.Context
	testConfig *test_config.TestConfig
	loginRepo  *Repository
	userRepo   *user.Repository
}

func (c *testConfig) init(t *testing.T) {
	c.ctx = context.Background()
	newTC, err := test_config.NewTestConfig(c.ctx)
	if err != nil {
		log.Fatal(err)
	}
	c.testConfig = newTC
	loginRepo, err := initLoginRepo(c.testConfig.DB)
	if err != nil {
		t.Fatal(err)
	}
	c.loginRepo = loginRepo
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
	firstUser, err := userfixture.CreateFirstUser(tc.ctx, tc.testConfig.DB, tc.userRepo)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		data    login.CreateLogin
		failure bool
	}{
		{
			name: "success",
			data: login.CreateLogin{
				Login:     "test login",
				Password:  "test password",
				UserID:    firstUser,
				Title:     "test title",
				Number:    1,
				CreatedAt: time.Now(),
			},
			failure: false,
		},
		{
			name: "failure - number unique constraint violation",
			data: login.CreateLogin{
				Login:     "test login",
				Password:  "test password",
				UserID:    firstUser,
				Title:     "test title",
				Number:    1,
				CreatedAt: time.Now(),
			},
			failure: true,
		},
		{
			name: "failure - no login",
			data: login.CreateLogin{
				Password:  "test password",
				UserID:    uuid.New(),
				Title:     "test title",
				Number:    2,
				CreatedAt: time.Now(),
			},
			failure: true,
		},
		{
			name: "failure - no password",
			data: login.CreateLogin{
				Login:     "test login",
				UserID:    firstUser,
				Title:     "test title",
				Number:    3,
				CreatedAt: time.Now(),
			},
			failure: true,
		},
		{
			name: "failure - no user id",
			data: login.CreateLogin{
				Login:     "test login",
				Password:  "test password",
				Title:     "test title",
				Number:    4,
				CreatedAt: time.Now(),
			},
			failure: true,
		},
		{
			name: "failure - title duplicate",
			data: login.CreateLogin{
				Title:     "test title",
				Login:     "test login",
				Password:  "test password",
				UserID:    firstUser,
				Number:    5,
				CreatedAt: time.Now(),
			},
			failure: true,
		},
		{
			name: "failure - no created at time",
			data: login.CreateLogin{
				Login:    "test login",
				Password: "test password",
				UserID:   firstUser,
				Title:    "test title",
				Number:   6,
			},
			failure: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = tc.loginRepo.Create(tc.ctx, tt.data)
			fmt.Println(tt.data)
			assert.Equal(t, tt.failure, err != nil)
		})
	}

}

func TestGetByNumber(t *testing.T) {
	tc := newTestConfig(t)
	firstUser, err := userfixture.CreateFirstUser(tc.ctx, tc.testConfig.DB, tc.userRepo)
	if err != nil {
		t.Fatal(err)
	}
	secondUser, err := userfixture.CreateSecondUser(tc.ctx, tc.testConfig.DB, tc.userRepo)
	if err != nil {
		t.Fatal(err)
	}

	fixtureData := login.CreateLogin{
		Login:    "test login",
		Password: "test password",
		UserID:   firstUser,
		Title:    "test title",
	}
	loginFixture(tc.ctx, tc.loginRepo, fixtureData, fixtureData.UserID)

	tests := []struct {
		name    string
		number  int64
		userID  uuid.UUID
		failure bool
	}{
		{
			name:    "success",
			number:  1,
			userID:  fixtureData.UserID,
			failure: false,
		},
		{
			name:    "failure - login does not exist",
			number:  2,
			userID:  firstUser,
			failure: true,
		},
		{
			name:    "failure - login does not exist",
			number:  1,
			userID:  secondUser,
			failure: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbEntity, err := tc.loginRepo.GetByNumber(tc.ctx, tt.number, tt.userID)
			if err != nil {
				if !tt.failure {
					fmt.Println(dbEntity)
					assert.Equal(t, fixtureData.UserID.String(), dbEntity.UserID.String())
					assert.Equal(t, fixtureData.Number, dbEntity.Number)
					assert.Equal(t, fixtureData.CreatedAt, dbEntity.CreatedAt)
					assert.Equal(t, fixtureData.UserID, dbEntity.Title)
					assert.Equal(t, fixtureData.Login, dbEntity.Login)
					assert.Equal(t, fixtureData.Password, dbEntity.Password)
				}
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	tc := newTestConfig(t)
	firstUser, err := userfixture.CreateFirstUser(tc.ctx, tc.testConfig.DB, tc.userRepo)
	if err != nil {
		t.Fatal(err)
	}
	secondUser, err := userfixture.CreateSecondUser(tc.ctx, tc.testConfig.DB, tc.userRepo)
	if err != nil {
		t.Fatal(err)
	}

	fixtureData := login.CreateLogin{
		Login:    "test login",
		Password: "test password",
		UserID:   firstUser,
		Title:    "test title",
	}
	loginFixture(tc.ctx, tc.loginRepo, fixtureData, fixtureData.UserID)
	tests := []struct {
		name    string
		userID  uuid.UUID
		data    login.UpdateLogin
		failure bool
	}{
		{
			name: "success",
			data: login.UpdateLogin{
				Login:       "updated login",
				Password:    "updated password",
				Title:       "updated title",
				URL:         "updated url",
				Description: "updated description",
			},
			failure: false,
			userID:  firstUser,
		},
		{
			name: "failure - login does not exist",
			data: login.UpdateLogin{
				Login:       "updated login",
				Password:    "updated password",
				Title:       "updated title",
				URL:         "updated url",
				Description: "updated description",
			},
			failure: true,
			userID:  secondUser,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = tc.loginRepo.Update(tc.ctx, tt.data, 1, tt.userID)
			assert.Equal(t, tt.failure, err != nil)
			if !tt.failure {
				dbEntity, err := tc.loginRepo.GetByNumber(tc.ctx, 1, tt.userID)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, firstUser, dbEntity.UserID)
				assert.Equal(t, tt.data.Title, dbEntity.Title)
				assert.Equal(t, tt.data.Description, dbEntity.Description)
				assert.Equal(t, tt.data.URL, dbEntity.URL)
				assert.Equal(t, tt.data.Description, dbEntity.Description)
				assert.Equal(t, tt.data.Login, dbEntity.Login)
				assert.Equal(t, tt.data.Password, dbEntity.Password)
			} else {
				assert.ErrorIs(t, err, constants.ErrEntityNotFound)
			}
		})
	}
}

func TestGetList(t *testing.T) {
	tc := newTestConfig(t)
	firstUser, err := userfixture.CreateFirstUser(tc.ctx, tc.testConfig.DB, tc.userRepo)
	if err != nil {
		t.Fatal(err)
	}
	secondUser, err := userfixture.CreateSecondUser(tc.ctx, tc.testConfig.DB, tc.userRepo)
	if err != nil {
		t.Fatal(err)
	}

	userLoginNumberMap := make(map[string]int)

	for i := 1; i < 8; i++ {
		var userID uuid.UUID
		if i%2 == 0 {
			userID = firstUser
		} else {
			userID = secondUser
		}
		fixtureData := login.CreateLogin{
			Login:    fmt.Sprintf("test login %d", i),
			Password: fmt.Sprintf("test password %d", i),
			UserID:   userID,
			Title:    fmt.Sprintf("test title %d", i),
		}
		loginFixture(tc.ctx, tc.loginRepo, fixtureData, firstUser)

		if _, ok := userLoginNumberMap[userID.String()]; !ok {
			userLoginNumberMap[userID.String()] = 1
		} else {
			userLoginNumberMap[userID.String()] += 1
		}
	}

	tests := []struct {
		name          string
		userID        uuid.UUID
		correctNumber int
	}{
		{
			name:          "first user",
			userID:        firstUser,
			correctNumber: userLoginNumberMap[firstUser.String()],
		},
		{
			name:          "second user",
			userID:        secondUser,
			correctNumber: userLoginNumberMap[secondUser.String()],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tc.loginRepo.GetList(tc.ctx, tt.userID)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.correctNumber, len(result))
			assert.IsType(t, result[0], login.GetListLoginResponse{})
			if tt.name == "first user" {
				assert.Equal(t, result[0].Number, int64(1))
			} else if tt.name == "second user" {
				assert.Equal(t, result[0].Number, int64(1))
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tc := newTestConfig(t)
	firstUser, err := userfixture.CreateFirstUser(tc.ctx, tc.testConfig.DB, tc.userRepo)
	if err != nil {
		t.Fatal(err)
	}
	secondUser, err := userfixture.CreateSecondUser(tc.ctx, tc.testConfig.DB, tc.userRepo)
	if err != nil {
		t.Fatal(err)
	}
	fixtureData := login.CreateLogin{
		Login:    "test login",
		Password: "test password",
		UserID:   firstUser,
		Title:    "test title",
	}
	loginFixture(tc.ctx, tc.loginRepo, fixtureData, fixtureData.UserID)

	tests := []struct {
		name    string
		userID  uuid.UUID
		number  int64
		failure bool
	}{
		{
			name:    "success",
			userID:  firstUser,
			number:  1,
			failure: false,
		},
		{
			name:    "failure - login does not exist",
			userID:  secondUser,
			number:  2,
			failure: true,
		},
		{
			name:    "failure - login does not exist",
			userID:  secondUser,
			number:  1,
			failure: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = tc.loginRepo.Delete(tc.ctx, tt.number, tt.userID)
			assert.Equal(t, tt.failure, err != nil)
		})
	}
}
