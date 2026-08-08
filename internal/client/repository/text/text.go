package text

import (
	"sync"

	"github.com/artni96/GophKeeper/internal/client/constants"
	"github.com/artni96/GophKeeper/internal/client/model/common"
	"github.com/artni96/GophKeeper/internal/client/model/text"
)

// RepositoryI represents the methods of the Text client repository.
type RepositoryI interface {
	Add(entity text.Text)
	AddBatch(entities []text.Text)
	Get(entityNumber uint64) (text.Text, error)
	GetList() []common.Entity
	Update(updatedEntity text.Text) error
	Delete(entityNumber uint64) error
}

// Repository implements the Text client repository to manage text-related data through the in-memory storages -
// shortDataStorage and extendedDataStorage.
type Repository struct {
	numberMap           map[uint64]string
	shortDataStorage    []common.Entity
	extendedDataStorage map[uint64]text.Text
	mu                  sync.Mutex
}

// NewRepository initializes and return the new Text client repository instance.
func NewRepository(numberMap map[uint64]string) *Repository {
	return &Repository{
		extendedDataStorage: make(map[uint64]text.Text),
		shortDataStorage:    []common.Entity{},
		numberMap:           numberMap,
	}
}

// Add adds a new Text entity into the in-memory storages.
func (repo *Repository) Add(entity text.Text) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.extendedDataStorage[entity.Number] = entity

	shortDataEntity := common.Entity{
		Number:      entity.Number,
		Title:       entity.Title,
		Description: entity.Description,
	}
	repo.shortDataStorage = append(repo.shortDataStorage, shortDataEntity)
	repo.numberMap[entity.Number] = "text"
}

// AddBatch adds a list of new Text entity into the in-memory storages.
func (repo *Repository) AddBatch(entities []text.Text) {
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
		repo.numberMap[entity.Number] = "text"
	}
}

// Get returns the Text entity from the extendedDataStorage.
func (repo *Repository) Get(entityNumber uint64) (text.Text, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	entity, ok := repo.extendedDataStorage[entityNumber]
	if !ok {
		return text.Text{}, constants.ErrEntityNotFound
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

// Update replaces the old Text entity with the new one into the in-memory storages.
func (repo *Repository) Update(updatedEntity text.Text) error {
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

// Delete removes the Text entity from the in-memory storages.
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
	delete(repo.extendedDataStorage, entityNumber)
	return nil
}
