package card

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/artni96/GophKeeper/internal/server/constants"
	"github.com/artni96/GophKeeper/internal/server/model/card"
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

func initCardRepo(db *sqlx.DB) (*Repository, error) {
	repo, err := NewRepository(db, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize test login repository: %w", err)
	}
	return repo, nil
}

func cardFixture(ctx context.Context, cardRepo *Repository, entityData card.CreateCard) {
	notNullFields := []string{
		"hashed_pan",
		"hashed_holder",
		"hashed_expiry_date",
		"hashed_cvv",
		"hashed_pin",
		"bank",
		"brand",
		"user_id",
		"title",
		"description",
		"created_at",
		"number",
	}
	_, err := cardRepo.Create(ctx, entityData, notNullFields)
	if err != nil {
		log.Fatal(err)
	}
}

type testConfig struct {
	ctx      context.Context
	depends  *tests.TestDependencies
	cardRepo *Repository
	userRepo *user.Repository
}

func (c *testConfig) init(t *testing.T) {
	c.ctx = context.Background()
	newTC, err := tests.NewTestDependencies(c.ctx)
	if err != nil {
		log.Fatal(err)
	}
	c.depends = newTC
	cardRepo, err := initCardRepo(c.depends.DB)
	if err != nil {
		t.Fatal(err)
	}
	c.cardRepo = cardRepo
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
		data          card.CreateCard
		notNullFields []string
		failure       bool
		error         error
	}{
		{
			name: "success",
			data: card.CreateCard{
				PAN:         1234567812345678,
				Holder:      "Test Holder",
				ExpiryDate:  "0535",
				CVV:         "123",
				PIN:         "1234",
				Bank:        "test",
				Brand:       "test",
				UserID:      firstUser,
				Title:       "Test Title",
				Description: "Test Description",
				CreatedAt:   time.Now(),
			},
			notNullFields: []string{
				"hashed_pan",
				"hashed_holder",
				"hashed_expiry_date",
				"hashed_cvv",
				"hashed_pin",
				"bank",
				"brand",
				"user_id",
				"title",
				"description",
				"created_at",
				"number",
			},
			failure: false,
		},
		{
			name: "failure - pan duplicate",
			data: card.CreateCard{
				PAN:         1234567812345678,
				Holder:      "Test Holder",
				ExpiryDate:  "0535",
				CVV:         "123",
				PIN:         "1234",
				Bank:        "test",
				Brand:       "test",
				UserID:      firstUser,
				Title:       "Test Title 1",
				Description: "Test Description",
				CreatedAt:   time.Now(),
			},
			failure: true,
			notNullFields: []string{
				"hashed_pan",
				"hashed_holder",
				"hashed_expiry_date",
				"hashed_cvv",
				"hashed_pin",
				"bank",
				"brand",
				"user_id",
				"title",
				"description",
				"created_at",
				"number",
			},
			error: constants.ErrEntityAlreadyExists,
		},
		{
			name: "failure - title duplicate",
			data: card.CreateCard{
				PAN:         1234567812345677,
				Holder:      "Test Holder",
				ExpiryDate:  "0535",
				CVV:         "123",
				PIN:         "1234",
				Bank:        "test",
				Brand:       "test",
				UserID:      firstUser,
				Title:       "Test Title",
				Description: "Test Description",
				CreatedAt:   time.Now(),
			},
			failure: true,
			notNullFields: []string{
				"hashed_pan",
				"hashed_holder",
				"hashed_expiry_date",
				"hashed_cvv",
				"hashed_pin",
				"bank",
				"brand",
				"user_id",
				"title",
				"description",
				"created_at",
				"number",
			},
			error: constants.ErrEntityAlreadyExists,
		},
		{
			name: "failure - no title",
			data: card.CreateCard{
				PAN:         1234567812345677,
				Holder:      "Test Holder",
				ExpiryDate:  "0535",
				CVV:         "123",
				PIN:         "1234",
				Bank:        "test",
				Brand:       "test",
				UserID:      firstUser,
				Description: "Test Description",
				CreatedAt:   time.Now(),
			},
			failure: true,
			notNullFields: []string{
				"hashed_pan",
				"hashed_holder",
				"hashed_expiry_date",
				"hashed_cvv",
				"hashed_pin",
				"bank",
				"brand",
				"user_id",
				"description",
				"created_at",
				"number",
			},
			error: constants.ErrRequiredField,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			number, err := tc.cardRepo.Create(tc.ctx, tt.data, tt.notNullFields)
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

	fixtureData := card.CreateCard{
		PAN:         1234567812345677,
		Holder:      "Test Holder",
		ExpiryDate:  "0535",
		CVV:         "123",
		PIN:         "1234",
		Bank:        "test",
		Brand:       "test",
		UserID:      firstUser,
		Title:       "Test Title",
		Description: "Test Description",
		CreatedAt:   time.Now(),
	}
	cardFixture(tc.ctx, tc.cardRepo, fixtureData)

	tests := []struct {
		name    string
		number  uint64
		userID  uuid.UUID
		failure bool
	}{
		{
			name:    "success",
			number:  1,
			userID:  firstUser,
			failure: false,
		},
		{
			name:    "failure - the first user does not have a card with number 2",
			number:  2,
			userID:  firstUser,
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
			dbEntity, err := tc.cardRepo.GetByNumber(tc.ctx, tt.number, tt.userID)
			assert.Equal(t, tt.failure, err != nil)
			if err != nil {
				if !tt.failure {
					assert.Equal(t, fixtureData.PAN, dbEntity.PAN)
					assert.Equal(t, fixtureData.Holder, dbEntity.Holder)
					assert.Equal(t, fixtureData.ExpiryDate, dbEntity.ExpiryDate)
					assert.Equal(t, fixtureData.CVV, dbEntity.CVV)
					assert.Equal(t, fixtureData.PIN, dbEntity.PIN)
					assert.Equal(t, fixtureData.Bank, dbEntity.Bank)
					assert.Equal(t, fixtureData.Brand, dbEntity.Brand)
					assert.Equal(t, fixtureData.UserID, dbEntity.UserID)
					assert.Equal(t, fixtureData.Title, dbEntity.Title)
					assert.Equal(t, fixtureData.Description, dbEntity.Description)
					assert.Equal(t, fixtureData.CreatedAt, dbEntity.CreatedAt)
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

	fixtureData := card.CreateCard{
		PAN:         1234567812345677,
		Holder:      "Test Holder",
		ExpiryDate:  "0535",
		CVV:         "123",
		PIN:         "1234",
		Bank:        "test",
		Brand:       "test",
		UserID:      firstUser,
		Title:       "Test Title",
		Description: "Test Description",
		CreatedAt:   time.Now(),
	}
	cardFixture(tc.ctx, tc.cardRepo, fixtureData)

	tests := []struct {
		name          string
		userID        uuid.UUID
		data          card.UpdateCard
		number        uint64
		notNullFields []string
		failure       bool
	}{
		{
			name:   "success",
			userID: firstUser,
			data: card.UpdateCard{
				PAN:         1234567812345678,
				Holder:      "Test Holder",
				ExpiryDate:  "0535",
				CVV:         "123",
				PIN:         "1234",
				Bank:        "test",
				Brand:       "test",
				Title:       "Test Title",
				Description: "Test Description",
				UpdatedAt:   time.Now(),
			},
			number: 1,
			notNullFields: []string{
				"hashed_pan",
				"hashed_holder",
				"hashed_expiry_date",
				"hashed_cvv",
				"hashed_pin",
				"bank",
				"brand",
				"title",
				"description",
				"updated_at",
			},
			failure: false,
		},
		{
			name:   "failure - the first user does not have a card with number 2",
			userID: firstUser,
			data: card.UpdateCard{
				PAN:         1234567812345678,
				Holder:      "Test Holder",
				ExpiryDate:  "0535",
				CVV:         "123",
				PIN:         "1234",
				Bank:        "test",
				Brand:       "test",
				Title:       "Test Title",
				Description: "Test Description",
				UpdatedAt:   time.Now(),
			},
			number: 2,
			notNullFields: []string{
				"hashed_pan",
				"hashed_holder",
				"hashed_expiry_date",
				"hashed_cvv",
				"hashed_pin",
				"bank",
				"brand",
				"title",
				"description",
				"updated_at",
			},
			failure: true,
		},
		{
			name:   "failure - the second user does not have a card with number 1",
			userID: secondUser,
			data: card.UpdateCard{
				PAN:         1234567812345678,
				Holder:      "Test Holder",
				ExpiryDate:  "0535",
				CVV:         "123",
				PIN:         "1234",
				Bank:        "test",
				Brand:       "test",
				Title:       "Test Title",
				Description: "Test Description",
				UpdatedAt:   time.Now(),
			},
			number: 1,
			notNullFields: []string{
				"hashed_pan",
				"hashed_holder",
				"hashed_expiry_date",
				"hashed_cvv",
				"hashed_pin",
				"bank",
				"brand",
				"title",
				"description",
				"updated_at",
			},
			failure: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = tc.cardRepo.Update(tc.ctx, tt.data, tt.number, tt.userID, tt.notNullFields)
			assert.Equal(t, tt.failure, err != nil)
			if !tt.failure {
				dbEntity, err := tc.cardRepo.GetByNumber(tc.ctx, tt.number, tt.userID)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.data.PAN, dbEntity.PAN)
				assert.Equal(t, tt.data.Holder, dbEntity.Holder)
				assert.Equal(t, tt.data.ExpiryDate, dbEntity.ExpiryDate)
				assert.Equal(t, tt.data.CVV, dbEntity.CVV)
				assert.Equal(t, tt.data.PIN, dbEntity.PIN)
				assert.Equal(t, tt.data.Bank, dbEntity.Bank)
				assert.Equal(t, tt.data.Brand, dbEntity.Brand)
				assert.Equal(t, tt.data.Title, dbEntity.Title)
				assert.Equal(t, tt.data.Description, dbEntity.Description)
				assert.Equal(t, tt.userID, dbEntity.UserID)
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

	userCardNumberMap := make(map[string]int)

	for i := 1; i < 8; i++ {
		var userID uuid.UUID
		if i%2 == 0 {
			userID = firstUser
		} else {
			userID = secondUser
		}
		strPan := fmt.Sprintf("123456781234567%d", i)
		testPan, err := strconv.Atoi(strPan)
		if err != nil {
			t.Fatal(err)
		}

		var testUserID uuid.UUID
		if i%2 == 0 {
			testUserID = firstUser
		} else {
			testUserID = secondUser
		}

		fixtureData := card.CreateCard{
			PAN:         uint64(testPan),
			Holder:      fmt.Sprintf("Test Holder %d", i),
			ExpiryDate:  fmt.Sprintf("053%d", i),
			CVV:         fmt.Sprintf("12%d", i),
			PIN:         fmt.Sprintf("123%d", i),
			Bank:        fmt.Sprintf("Test Bank %d", i),
			Brand:       fmt.Sprintf("Test Brand %d", i),
			UserID:      testUserID,
			Title:       fmt.Sprintf("Test Title %d", i),
			Description: fmt.Sprintf("Test Description %d", i),
			CreatedAt:   time.Now(),
		}
		cardFixture(tc.ctx, tc.cardRepo, fixtureData)

		if _, ok := userCardNumberMap[userID.String()]; !ok {
			userCardNumberMap[userID.String()] = 1
		} else {
			userCardNumberMap[userID.String()] += 1
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
			correctNumber: userCardNumberMap[firstUser.String()],
		},
		{
			name:          "second user",
			userID:        secondUser,
			correctNumber: userCardNumberMap[secondUser.String()],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tc.cardRepo.GetList(tc.ctx, tt.userID)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(len(result))
			assert.Equal(t, tt.correctNumber, len(result))
			assert.IsType(t, result[0], card.Card{})
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

	fixtureData := card.CreateCard{
		PAN:         1234567812345677,
		Holder:      "Test Holder",
		ExpiryDate:  "0535",
		CVV:         "123",
		PIN:         "1234",
		Bank:        "test",
		Brand:       "test",
		UserID:      firstUser,
		Title:       "Test Title",
		Description: "Test Description",
		CreatedAt:   time.Now(),
	}
	cardFixture(tc.ctx, tc.cardRepo, fixtureData)

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
			name:    "failure - the first user does not have a card number with number 1 (already deleted)",
			userID:  firstUser,
			number:  1,
			failure: true,
		},
		{
			name:    "failure - the second user does not have a card number with number 1",
			userID:  secondUser,
			number:  1,
			failure: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = tc.cardRepo.Delete(tc.ctx, tt.number, tt.userID)
			assert.Equal(t, tt.failure, err != nil)
		})
	}
}
