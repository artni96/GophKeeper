package text

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/artni96/GophKeeper/internal/server/constants"
	"github.com/artni96/GophKeeper/internal/server/model/text"
	"github.com/artni96/GophKeeper/internal/server/repository/user"
	userfixture "github.com/artni96/GophKeeper/internal/server/repository/user/fixture"

	"github.com/artni96/GophKeeper/tests"
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

func initTextRepo(db *sqlx.DB) (*Repository, error) {
	repo, err := NewRepository(db, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize test login repository: %w", err)
	}
	return repo, nil
}

func textFixture(ctx context.Context, cardRepo *Repository, entityData text.CreateText) uint64 {
	notNullFields := []string{
		"title",
		"description",

		"hashed_text",
		"nonce",
		"key_id",

		"user_id",
		"created_at",
		"number",
	}
	number, err := cardRepo.Create(ctx, entityData, notNullFields)
	if err != nil {
		log.Fatal(err)
	}
	return number
}

type testConfig struct {
	ctx      context.Context
	depends  *tests.TestDependencies
	textRepo *Repository
	userRepo *user.Repository
}

func (c *testConfig) init(t *testing.T) {
	c.ctx = context.Background()
	newTC, err := tests.NewTestDependencies(c.ctx)
	if err != nil {
		log.Fatal(err)
	}
	c.depends = newTC
	textRepo, err := initTextRepo(c.depends.DB)
	if err != nil {
		t.Fatal(err)
	}
	c.textRepo = textRepo
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

func TestCreateText(t *testing.T) {
	tc := newTestConfig(t)
	firstUser, err := userfixture.CreateFirstUser(tc.ctx, tc.depends.DB, tc.userRepo)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name          string
		data          text.CreateText
		notNullFields []string
		failure       bool
		error         error
	}{
		{
			name: "success",
			data: text.CreateText{
				Title:       "test title",
				Description: "test description",

				Text:  []byte("test text"),
				Nonce: []byte("test nonce"),
				KeyID: 1,

				UserID:    firstUser,
				CreatedAt: time.Now(),
			},
			notNullFields: []string{
				"title",
				"description",

				"hashed_text",
				"nonce",
				"key_id",

				"user_id",
				"created_at",
				"number",
			},
			failure: false,
			error:   nil,
		},
		{
			name: "failure - title duplicate",
			data: text.CreateText{
				Title:       "test title",
				Description: "test description",

				Text:  []byte("test text"),
				Nonce: []byte("test nonce"),
				KeyID: 1,

				UserID:    firstUser,
				CreatedAt: time.Now(),
			},
			notNullFields: []string{
				"title",
				"description",

				"hashed_text",
				"nonce",
				"key_id",

				"user_id",
				"created_at",
				"number",
			},
			failure: true,
			error:   constants.ErrEntityAlreadyExists,
		},
		{
			name: "failure - no title",
			data: text.CreateText{
				Description: "test description",

				Text:  []byte("test text"),
				Nonce: []byte("test nonce"),
				KeyID: 1,

				UserID:    firstUser,
				CreatedAt: time.Now(),
			},
			notNullFields: []string{
				"description",

				"hashed_text",
				"nonce",
				"key_id",

				"user_id",
				"created_at",
				"number",
			},
			failure: true,
			error:   constants.ErrRequiredField,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			number, err := tc.textRepo.Create(tc.ctx, tt.data, tt.notNullFields)
			assert.Equal(t, tt.failure, err != nil)
			if tt.failure {
				assert.ErrorIs(t, err, tt.error)
			} else {
				assert.Equal(t, number, uint64(1))
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

	fixtureData := text.CreateText{
		Title:       "test title",
		Description: "test description",

		Text:  []byte("test text"),
		Nonce: []byte("test nonce"),
		KeyID: 1,

		UserID:    firstUser,
		CreatedAt: time.Now(),
	}
	entityNumber := textFixture(tc.ctx, tc.textRepo, fixtureData)

	tests := []struct {
		name    string
		number  uint64
		userID  uuid.UUID
		failure bool
	}{
		{
			name:    "success",
			number:  entityNumber,
			userID:  firstUser,
			failure: false,
		},
		{
			name:    "failure - the first user does not have a text record with number 2",
			number:  2,
			userID:  firstUser,
			failure: true,
		},
		{
			name:    fmt.Sprintf("failure - the second user does not have a text record with number %d", entityNumber),
			number:  entityNumber,
			userID:  secondUser,
			failure: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbEntity, err := tc.textRepo.GetByNumber(tc.ctx, tt.number, tt.userID)
			assert.Equal(t, tt.failure, err != nil)
			if err != nil {
				if !tt.failure {
					assert.Equal(t, fixtureData.Title, dbEntity.Title)
					assert.Equal(t, fixtureData.Description, dbEntity.Description)
					assert.Equal(t, fixtureData.Text, dbEntity.Text)
					assert.Equal(t, tt.userID, dbEntity.UserID)
					assert.Equal(t, tt.number, dbEntity.Number)
					assert.Equal(t, fixtureData.CreatedAt, dbEntity.CreatedAt)
				}
			}
		})
	}
}

func TestUpdateText(t *testing.T) {
	tc := newTestConfig(t)
	firstUser, err := userfixture.CreateFirstUser(tc.ctx, tc.depends.DB, tc.userRepo)
	if err != nil {
		t.Fatal(err)
	}

	secondUser, err := userfixture.CreateSecondUser(tc.ctx, tc.depends.DB, tc.userRepo)
	if err != nil {
		t.Fatal(err)
	}

	fixtureData := text.CreateText{
		Title:       "update test title",
		Description: "update test description",

		Text:  []byte("update test text"),
		Nonce: []byte("update test nonce"),
		KeyID: 1,

		UserID:    firstUser,
		CreatedAt: time.Now(),
	}
	entityNumber := textFixture(tc.ctx, tc.textRepo, fixtureData)

	tests := []struct {
		name          string
		userID        uuid.UUID
		data          text.UpdateText
		number        uint64
		notNullFields []string
		failure       bool
		error         error
	}{
		{
			name:   "success",
			userID: firstUser,
			data: text.UpdateText{
				Title:       "updated title",
				Description: "updated description",

				Text:  []byte("updated text"),
				Nonce: []byte("update test nonce"),
				KeyID: 1,

				UpdatedAt: time.Now(),
			},

			number: entityNumber,
			notNullFields: []string{
				"title",
				"description",

				"hashed_text",
				"nonce",
				"key_id",

				"updated_at",
			},
			failure: false,
			error:   nil,
		},
		{
			name:   "failure - the first user does not have a text record with number 2",
			userID: firstUser,
			data: text.UpdateText{
				Title:       "updated title",
				Description: "updated description",

				Text:  []byte("updated text"),
				Nonce: []byte("update test nonce"),
				KeyID: 1,

				UpdatedAt: time.Now(),
			},
			number: 2,
			notNullFields: []string{
				"title",
				"description",

				"hashed_text",
				"nonce",
				"key_id",

				"updated_at",
			},
			failure: true,
			error:   constants.ErrEntityNotFound,
		},
		{
			name:   "failure - the second user does not have a text record with number 1",
			userID: secondUser,
			data: text.UpdateText{
				Title:       "updated title",
				Description: "updated description",

				Text:  []byte("updated text"),
				Nonce: []byte("update test nonce"),
				KeyID: 1,

				UpdatedAt: time.Now(),
			},
			number: entityNumber,
			notNullFields: []string{
				"title",
				"description",

				"hashed_text",
				"nonce",
				"key_id",

				"updated_at",
			},
			failure: true,
			error:   constants.ErrEntityNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = tc.textRepo.Update(tc.ctx, tt.data, tt.number, tt.userID, tt.notNullFields)
			fmt.Println(tt.data)
			fmt.Println(tt.number)
			assert.Equal(t, tt.failure, err != nil)
			if !tt.failure {
				dbEntity, err := tc.textRepo.GetByNumber(tc.ctx, tt.number, tt.userID)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.data.Title, dbEntity.Title)
				assert.Equal(t, tt.data.Description, dbEntity.Description)
				assert.Equal(t, tt.data.Text, dbEntity.Text)
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

	userTextNumberMap := make(map[string]int)

	for i := 1; i < 8; i++ {
		var userID uuid.UUID
		if i%2 == 0 {
			userID = firstUser
		} else {
			userID = secondUser
		}

		var testUserID uuid.UUID
		if i%2 == 0 {
			testUserID = firstUser
		} else {
			testUserID = secondUser
		}

		fixtureData := text.CreateText{
			Title:       fmt.Sprintf("test title %d", i),
			Description: fmt.Sprintf("test description %d", i),

			Text:  []byte(fmt.Sprintf("test text %d", i)),
			Nonce: []byte(fmt.Sprintf("test nonce %d", i)),
			KeyID: uint64(i),

			UserID:    testUserID,
			CreatedAt: time.Now(),
		}
		_ = textFixture(tc.ctx, tc.textRepo, fixtureData)

		if _, ok := userTextNumberMap[userID.String()]; !ok {
			userTextNumberMap[userID.String()] = 1
		} else {
			userTextNumberMap[userID.String()] += 1
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
			correctNumber: userTextNumberMap[firstUser.String()],
		},
		{
			name:          "second user",
			userID:        secondUser,
			correctNumber: userTextNumberMap[secondUser.String()],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tc.textRepo.GetList(tc.ctx, tt.userID)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(len(result))
			assert.Equal(t, tt.correctNumber, len(result))
			assert.IsType(t, result[0], text.Text{})
			assert.Equal(t, result[0].Number, uint64(1))
		})
	}
}

func TestDeleteText(t *testing.T) {
	tc := newTestConfig(t)
	firstUser, err := userfixture.CreateFirstUser(tc.ctx, tc.depends.DB, tc.userRepo)
	if err != nil {
		t.Fatal(err)
	}

	secondUser, err := userfixture.CreateSecondUser(tc.ctx, tc.depends.DB, tc.userRepo)
	if err != nil {
		t.Fatal(err)
	}

	fixtureData := text.CreateText{
		Title:       "update test title",
		Description: "update test description",

		Text:  []byte("update test text"),
		Nonce: []byte("update test nonce"),
		KeyID: 1,

		UserID:    firstUser,
		CreatedAt: time.Now(),
	}
	entityNumber := textFixture(tc.ctx, tc.textRepo, fixtureData)

	tests := []struct {
		name    string
		userID  uuid.UUID
		number  uint64
		failure bool
		error   error
	}{
		{
			name:    "success",
			userID:  firstUser,
			number:  entityNumber,
			failure: false,
			error:   nil,
		},
		{
			name:    "failure - the first user does not have a card number with number 1 (already deleted)",
			userID:  firstUser,
			number:  1,
			failure: true,
			error:   constants.ErrEntityNotFound,
		},
		{
			name:    "failure - the second user does not have a card number with number 1",
			userID:  secondUser,
			number:  1,
			failure: true,
			error:   constants.ErrEntityNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tc.textRepo.Delete(tc.ctx, tt.number, tt.userID)
			assert.Equal(t, tt.failure, err != nil)
			if tt.failure {
				assert.ErrorIs(t, err, tt.error)
			}
		})
	}
}
