package repository

import (
	"fmt"
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/db"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/entity"
	"github.com/pkg/errors"
)

var (
	ErrPersistingCollection = errors.New("err happened during the process of collection creation response data save.")
)

var (
	// GetCollectionsRepository retrieves instance of the collections repository, performing initial creation if needed.
	GetCollectionsRepository = sync.OnceValue[CollectionsRepository](createCollectionsRepository)
)

// CollectionsRepository represents collections entity repository.
type CollectionsRepository interface {
	Insert(name string) error
	IsEmpty() (bool, error)
	GetAll() ([]entity.CollectionEntity, error)
}

// collectionsRepositoryImpl represents implementation of CollectionsRepository.
type collectionsRepositoryImpl struct{}

// Insert inserts new collection entity to the storage.
func (w *collectionsRepositoryImpl) Insert(name string) error {
	instance := db.GetInstance()

	err := instance.Create(&entity.CollectionEntity{Name: name}).Error

	return errors.Wrap(err, ErrPersistingCollection.Error())
}

// IsEmpty checks if any collection is set.
func (w *collectionsRepositoryImpl) IsEmpty() (bool, error) {
	instance := db.GetInstance()

	var result bool

	err := instance.Raw(
		fmt.Sprintf(
			"SELECT COUNT(*) = 0 FROM %s",
			(&entity.CollectionEntity{}).TableName()),
	).Scan(&result).Error

	return result, err
}

// GetAll retrieves all available collections.
func (w *collectionsRepositoryImpl) GetAll() ([]entity.CollectionEntity, error) {
	instance := db.GetInstance()

	var result []entity.CollectionEntity

	err := instance.Table((&entity.CollectionEntity{}).TableName()).
		Find(&result).Error

	return result, err
}

// createCollectionsRepository initializes collectionsRepositoryImpl.
func createCollectionsRepository() CollectionsRepository {
	return new(collectionsRepositoryImpl)
}
