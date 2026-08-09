package login

import (
	"sync"

	"github.com/artni96/GophKeeper/internal/client/constants"
	"github.com/artni96/GophKeeper/internal/client/model/common"
	"github.com/artni96/GophKeeper/internal/client/model/login"
)

// RepositoryI represents the methods of the Login client repository.
type RepositoryI interface {
	Add(entity login.Login)
	AddBatch(entities []login.Login)
	Get(entityNumber uint64) (login.Login, error)
	GetList() []common.Entity
	Update(updatedEntity login.Login) error
	Delete(entityNumber uint64) error
}

// Repository implements the Login client repository to manage login-related data through the in-memory storages -
// shortDataStorage and extendedDataStorage.
type Repository struct {
	numberMap           map[uint64]string
	shortDataStorage    []common.Entity
	extendedDataStorage map[uint64]login.Login
	mu                  sync.Mutex
}

// NewRepository initializes and return the new Login client repository instance.
func NewRepository(numberMap map[uint64]string) *Repository {
	return &Repository{
		extendedDataStorage: make(map[uint64]login.Login),
		shortDataStorage:    []common.Entity{},
		numberMap:           numberMap,
	}
}

// Add adds a new Login entity into the in-memory storages.
func (repo *Repository) Add(entity login.Login) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.extendedDataStorage[entity.Number] = entity

	shortDataEntity := common.Entity{
		Number:      entity.Number,
		Title:       entity.Title,
		Description: entity.Description,
	}
	repo.shortDataStorage = append(repo.shortDataStorage, shortDataEntity)
	repo.numberMap[entity.Number] = "login"
}

// AddBatch adds a list of new Login entity into the in-memory storages.
func (repo *Repository) AddBatch(entities []login.Login) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, entity := range entities {
		repo.extendedDataStorage[entity.Number] = entity

		shortDataEntity := common.Entity{
			Title:       entity.Title,
			Number:      entity.Number,
			Description: entity.Description,
		}
		repo.shortDataStorage = append(repo.shortDataStorage, shortDataEntity)
		repo.numberMap[entity.Number] = "login"
	}
}

// Get returns the Login entity from the extendedDataStorage.
func (repo *Repository) Get(entityNumber uint64) (login.Login, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	entity, ok := repo.extendedDataStorage[entityNumber]
	if !ok {
		return login.Login{}, constants.ErrEntityNotFound
	}

	return entity, nil
}

// GetList returns the list of Card entities with the basic data.
func (repo *Repository) GetList() []common.Entity {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	entities := append([]common.Entity{}, repo.shortDataStorage...)
	return entities
}

// Update replaces the old Login entity with the new one into the in-memory storages.
func (repo *Repository) Update(updatedEntity login.Login) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	entityNumber := updatedEntity.Number

	_, ok := repo.extendedDataStorage[entityNumber]
	if !ok {
		return constants.ErrEntityNotFound
	}
	repo.extendedDataStorage[entityNumber] = updatedEntity

	for i := range repo.shortDataStorage {
		if repo.shortDataStorage[i].Number == entityNumber {
			repo.shortDataStorage[i].Description = updatedEntity.Description
			repo.shortDataStorage[i].Title = updatedEntity.Title
		}
	}
	return nil
}

// Delete removes the Login entity from the in-memory storages.
func (repo *Repository) Delete(entityNumber uint64) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	_, ok := repo.extendedDataStorage[entityNumber]
	if !ok {
		return constants.ErrEntityNotFound
	}

	delete(repo.extendedDataStorage, entityNumber)

	for i := range repo.shortDataStorage {
		if repo.shortDataStorage[i].Number == entityNumber {
			repo.shortDataStorage = append(repo.shortDataStorage[:i], repo.shortDataStorage[i+1:]...)
		}
	}
	delete(repo.numberMap, entityNumber)
	return nil
}
