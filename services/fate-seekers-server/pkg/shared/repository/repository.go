package repository

import (
	"fmt"
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/db"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/entity"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	ErrPersistingCollection = errors.New("err happened during the process of collection creation response data save.")
	ErrPersistingFlags      = errors.New("err happened during the process of flag creation response data save.")
)

var (
	// GetCollectionsRepository retrieves instance of the collections repository, performing initial creation if needed.
	GetCollectionsRepository = sync.OnceValue[CollectionsRepository](createCollectionsRepository)

	// GetFlagsRepository retrieves instance of the collections repository, performing initial creation if needed.
	GetFlagsRepository = sync.OnceValue[FlagsRepository](createFlagsRepository)
)

// CollectionsRepository represents collections entity repository.
type CollectionsRepository interface {
	Insert(name string) error
	Exists(name string) (bool, error)
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

// Exists checks if any collection with the given name exists.
func (w *collectionsRepositoryImpl) Exists(name string) (bool, error) {
	instance := db.GetInstance()

	err := instance.Model(&entity.CollectionEntity{}).
		Where("name = ?", name).
		First(&entity.CollectionEntity{}).Error

	if err != nil {
		return false, err
	}

	return err != gorm.ErrRecordNotFound, nil
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

// FlagsRepository represents flags entity repository.
type FlagsRepository interface {
	InsertOrUpdate(name, value string) error
	GetByName(name string) (*entity.FlagsEntity, bool, error)
}

// flagsRepositoryImpl represents implementation of FlagsRepository.
type flagsRepositoryImpl struct{}

// InsertOrUpdate inserts or updates flags entity to the storage.
func (w *flagsRepositoryImpl) InsertOrUpdate(name, value string) error {
	instance := db.GetInstance()

	err := instance.Create(&entity.FlagsEntity{
		Name:  name,
		Value: value,
	}).Error

	return errors.Wrap(err, ErrPersistingFlags.Error())
}

// GetByName checks if any flag with the given name exists.
func (w *flagsRepositoryImpl) GetByName(name string) (*entity.FlagsEntity, bool, error) {
	instance := db.GetInstance()

	var result *entity.FlagsEntity

	err := instance.Model(&entity.FlagsEntity{}).
		Where("name = ?", name).
		First(&result).Error

	if err != nil {
		return nil, false, err
	}

	return result, err != gorm.ErrRecordNotFound, nil
}

// createFlagsRepository initializes flagsRepositoryImpl.
func createFlagsRepository() FlagsRepository {
	return new(flagsRepositoryImpl)
}
