package card

import (
	"sync"

	"github.com/artni96/GophKeeper/internal/client/constants"
	"github.com/artni96/GophKeeper/internal/client/model/card"
	"github.com/artni96/GophKeeper/internal/client/model/common"
)

// RepositoryI represents the methods of the Card client repository.
type RepositoryI interface {
	Add(entity card.Card)
	AddBatch(entities []card.Card)
	Get(entityNumber uint64) (card.Card, error)
	GetList() []common.Entity
	Update(updatedEntity card.Card) error
	Delete(entityNumber uint64) error
}

// Repository implements the Card client repository to manage card-related data through the in-memory storages -
// shortDataStorage and extendedDataStorage.
type Repository struct {
	numberMap           map[uint64]string
	shortDataStorage    []common.Entity
	extendedDataStorage map[uint64]card.Card
	mu                  sync.Mutex
}

// NewRepository initializes and return the new Card client repository instance.
func NewRepository(numberMap map[uint64]string) *Repository {
	return &Repository{
		extendedDataStorage: make(map[uint64]card.Card),
		shortDataStorage:    []common.Entity{},
		numberMap:           numberMap,
	}
}

// Add adds a new Card entity into the in-memory storages.
func (repo *Repository) Add(entity card.Card) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.extendedDataStorage[entity.Number] = entity

	shortDataEntity := common.Entity{
		Number:      entity.Number,
		Title:       entity.Title,
		Description: entity.Description,
	}
	repo.shortDataStorage = append(repo.shortDataStorage, shortDataEntity)
	repo.numberMap[entity.Number] = "card"
}

// AddBatch adds a list of new Card entity into the in-memory storages.
func (repo *Repository) AddBatch(entities []card.Card) {
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
		repo.numberMap[entity.Number] = "card"
	}
}

// Get returns the Card entity from the extendedDataStorage.
func (repo *Repository) Get(entityNumber uint64) (card.Card, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	entity, ok := repo.extendedDataStorage[entityNumber]
	if !ok {
		return card.Card{}, constants.ErrEntityNotFound
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

// Update replaces the old Card entity with the new one into the in-memory storages.
func (repo *Repository) Update(updatedEntity card.Card) error {
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

// Delete removes the Card entity from the in-memory storages.
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
