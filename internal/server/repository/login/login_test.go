package login

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/artni96/GophKeeper/internal/server/constants"
	"github.com/artni96/GophKeeper/internal/server/model/common"
	"github.com/artni96/GophKeeper/internal/server/model/login"
	"github.com/artni96/GophKeeper/internal/server/repository/user"
	userfixture "github.com/artni96/GophKeeper/internal/server/repository/user/fixture"
	"github.com/artni96/GophKeeper/tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func initUserRepo(db *sql.DB) (*user.Repository, error) {
	repo, err := user.NewRepository(db, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize test user repository: %w", err)
	}
	return repo, nil
}

func initLoginRepo(db *sql.DB) (*Repository, error) {
	repo, err := NewRepository(db, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize test login repository: %w", err)
	}
	return repo, nil
}

func loginFixture(ctx context.Context, loginRepo *Repository, entityData login.CreateLogin) {
	notNullFields := []string{
		"title",
		"description",

		"login_value",
		"login_nonce",
		"login_key_id",

		"password_value",
		"password_nonce",
		"password_key_id",

		"url",
		"user_id",
		"created_at",
		"number",
	}
	_, err := loginRepo.Create(ctx, entityData, notNullFields)
	if err != nil {
		log.Fatal(err)
	}
}

type testConfig struct {
	ctx       context.Context
	depends   *tests.TestDependencies
	loginRepo *Repository
	userRepo  *user.Repository
}

func (c *testConfig) init(t *testing.T) {
	c.ctx = context.Background()
	newTC, err := tests.NewTestDependencies(c.ctx)
	if err != nil {
		log.Fatal(err)
	}
	c.depends = newTC
	loginRepo, err := initLoginRepo(c.depends.DB)
	if err != nil {
		t.Fatal(err)
	}
	c.loginRepo = loginRepo
	userRepo, err := initUserRepo(c.depends.DB)
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
	firstUser, err := userfixture.CreateFirstUser(tc.ctx, tc.depends.DB, tc.userRepo)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name          string
		data          login.CreateLogin
		notNullFields []string
		failure       bool
		error         error
	}{
		{
			name: "success",
			data: login.CreateLogin{
				Login: common.EncryptedField{
					Value: []byte("test login creation"),
					Nonce: []byte("test nonce"),
					KeyID: 1,
				},

				Password: common.EncryptedField{
					Value: []byte("test password"),
					Nonce: []byte("test nonce"),
					KeyID: 1,
				},

				UserID: firstUser,
				Title:  "test title creation",

				CreatedAt:   time.Now(),
				Description: "test description creation",
				URL:         "test url creation",
			},
			notNullFields: []string{
				"title",
				"description",

				"login_value",
				"login_nonce",
				"login_key_id",

				"password_value",
				"password_nonce",
				"password_key_id",

				"url",
				"user_id",
				"created_at",
				"number",
			},
			failure: false,
		},

		{
			name: "failure - no user id",
			data: login.CreateLogin{
				Login: common.EncryptedField{
					Value: []byte("test login creation"),
					Nonce: []byte("test nonce"),
					KeyID: 1,
				},

				Password: common.EncryptedField{
					Value: []byte("test password"),
					Nonce: []byte("test nonce"),
					KeyID: 1,
				},

				UserID: firstUser,
				Title:  "test title creation",

				CreatedAt:   time.Now(),
				Description: "test description creation",
				URL:         "test url creation",
			},
			notNullFields: []string{
				"title",
				"description",
				"login_value",
				"login_nonce",
				"login_key_id",

				"password_value",
				"password_nonce",
				"password_key_id",

				"url",
				"created_at",
				"number",
			},
			failure: true,
			error:   constants.ErrEntityAlreadyExists,
		},
		{
			name: "failure - title duplicate",
			data: login.CreateLogin{
				Login: common.EncryptedField{
					Value: []byte("test login creation"),
					Nonce: []byte("test nonce"),
					KeyID: 1,
				},

				Password: common.EncryptedField{
					Value: []byte("test password"),
					Nonce: []byte("test nonce"),
					KeyID: 1,
				},

				UserID: firstUser,
				Title:  "test title creation",

				CreatedAt:   time.Now(),
				Description: "test description creation",
				URL:         "test url creation",
			},
			notNullFields: []string{
				"title",
				"description",
				"login_value",
				"login_nonce",
				"login_key_id",

				"password_value",
				"password_nonce",
				"password_key_id",

				"url",
				"user_id",
				"created_at",
				"number",
			},
			failure: true,
			error:   constants.ErrEntityAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			number, err := tc.loginRepo.Create(tc.ctx, tt.data, tt.notNullFields)
			assert.Equal(t, tt.failure, err != nil)
			if !tt.failure {
				assert.Equal(t, uint64(1), number)
			} else {
				assert.ErrorIs(t, err, tt.error)
			}
		})
	}
}

func TestGetByNumber(t *testing.T) {
	tc := newTestConfig(t)
	firstUser, err := userfixture.CreateFirstUser(tc.ctx, tc.depends.DB, tc.userRepo)
	if err != nil {
		t.Fatal(err)
	}
	secondUser, err := userfixture.CreateSecondUser(tc.ctx, tc.depends.DB, tc.userRepo)
	if err != nil {
		t.Fatal(err)
	}

	fixtureData := login.CreateLogin{
		Login: common.EncryptedField{
			Value: []byte("test login"),
			Nonce: []byte("test nonce"),
			KeyID: 1,
		},

		Password: common.EncryptedField{
			Value: []byte("test password"),
			Nonce: []byte("test nonce"),
			KeyID: 1,
		},

		UserID: firstUser,
		Title:  "test title",

		CreatedAt:   time.Now(),
		Description: "test description",
		URL:         "test url",
	}
	loginFixture(tc.ctx, tc.loginRepo, fixtureData)

	tests := []struct {
		name    string
		number  uint64
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
			name:    "failure - the first user does not have a login with number 2",
			number:  2,
			userID:  fixtureData.UserID,
			failure: true,
		},
		{
			name:    "failure - the second user does not have a card with number 1",
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
	firstUser, err := userfixture.CreateFirstUser(tc.ctx, tc.depends.DB, tc.userRepo)
	if err != nil {
		t.Fatal(err)
	}
	secondUser, err := userfixture.CreateSecondUser(tc.ctx, tc.depends.DB, tc.userRepo)
	if err != nil {
		t.Fatal(err)
	}

	fixtureData := login.CreateLogin{
		Login: common.EncryptedField{
			Value: []byte("test login"),
			Nonce: []byte("test nonce"),
			KeyID: 1,
		},

		Password: common.EncryptedField{
			Value: []byte("test password"),
			Nonce: []byte("test nonce"),
			KeyID: 1,
		},

		UserID: firstUser,
		Title:  "test title",

		CreatedAt:   time.Now(),
		Description: "test description",
		URL:         "test url",
	}
	loginFixture(tc.ctx, tc.loginRepo, fixtureData)
	tests := []struct {
		name           string
		userID         uuid.UUID
		data           login.UpdateLogin
		number         uint64
		fieldsToUpdate []string
		failure        bool
	}{
		{
			name: "success",
			data: login.UpdateLogin{
				Login: common.EncryptedField{
					Value: []byte("updated login"),
					Nonce: []byte("updated nonce"),
					KeyID: 2,
				},

				Password: common.EncryptedField{
					Value: []byte("updated password"),
					Nonce: []byte("updated nonce"),
					KeyID: 2,
				},

				Title:       "updated title",
				URL:         "updated url",
				Description: "updated description",
			},
			number: uint64(1),
			fieldsToUpdate: []string{
				"login_value",
				"login_nonce",
				"login_key_id",

				"password_value",
				"password_nonce",
				"password_key_id",

				"title",
				"url",
				"description",
				"nonce",
				"key_id",
			},
			failure: false,
			userID:  firstUser,
		},
		{
			name: "failure - login does not exist",
			data: login.UpdateLogin{
				Login: common.EncryptedField{
					Value: []byte("updated login"),
					Nonce: []byte("updated nonce"),
					KeyID: 2,
				},

				Password: common.EncryptedField{
					Value: []byte("updated password"),
					Nonce: []byte("updated nonce"),
					KeyID: 2,
				},

				Title:       "updated title",
				URL:         "updated url",
				Description: "updated description",
			},
			number: uint64(1),
			fieldsToUpdate: []string{
				"login_value",
				"login_nonce",
				"login_key_id",

				"password_value",
				"password_nonce",
				"password_key_id",

				"title",
				"url",
				"description",
				"nonce",
				"key_id",
			},
			failure: true,
			userID:  secondUser,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = tc.loginRepo.Update(tc.ctx, tt.data, tt.number, tt.userID, tt.fieldsToUpdate)
			assert.Equal(t, tt.failure, err != nil)
			if !tt.failure {
				dbEntity, err := tc.loginRepo.GetByNumber(tc.ctx, tt.number, tt.userID)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.userID, dbEntity.UserID)
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
	firstUser, err := userfixture.CreateFirstUser(tc.ctx, tc.depends.DB, tc.userRepo)
	if err != nil {
		t.Fatal(err)
	}
	secondUser, err := userfixture.CreateSecondUser(tc.ctx, tc.depends.DB, tc.userRepo)
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
			Login: common.EncryptedField{
				Value: []byte(fmt.Sprintf("test login %d", i)),
				Nonce: []byte(fmt.Sprintf("test nonce %d", i)),
				KeyID: 1,
			},

			Password: common.EncryptedField{
				Value: []byte(fmt.Sprintf("test password %d", i)),
				Nonce: []byte(fmt.Sprintf("test nonce %d", i)),
				KeyID: 1,
			},

			UserID:      userID,
			Title:       fmt.Sprintf("test title %d", i),
			Description: fmt.Sprintf("test description %d", i),
			URL:         fmt.Sprintf("test url %d", i),
		}
		loginFixture(tc.ctx, tc.loginRepo, fixtureData)

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
			assert.IsType(t, result[0], login.Login{})
			assert.Equal(t, result[0].Number, uint64(1))
		})
	}
}

func TestDelete(t *testing.T) {
	tc := newTestConfig(t)
	firstUser, err := userfixture.CreateFirstUser(tc.ctx, tc.depends.DB, tc.userRepo)
	if err != nil {
		t.Fatal(err)
	}
	secondUser, err := userfixture.CreateSecondUser(tc.ctx, tc.depends.DB, tc.userRepo)
	if err != nil {
		t.Fatal(err)
	}
	fixtureData := login.CreateLogin{
		Login: common.EncryptedField{
			Value: []byte("test login"),
			Nonce: []byte("test nonce"),
			KeyID: 1,
		},

		Password: common.EncryptedField{
			Value: []byte("test password"),
			Nonce: []byte("test nonce"),
			KeyID: 1,
		},

		UserID:      firstUser,
		Title:       "test title",
		Description: "test description",
		URL:         "test url",
	}
	loginFixture(tc.ctx, tc.loginRepo, fixtureData)

	tests := []struct {
		name    string
		userID  uuid.UUID
		number  uint64
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
