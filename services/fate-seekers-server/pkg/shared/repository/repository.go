package repository

import (
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/db"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/entity"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	ErrPersistingSessions = errors.New("err happened during the process of session creation response data save.")
	ErrPersistingMessages = errors.New("err happened during the process of message creation response data save.")
	ErrPersistingUsers    = errors.New("err happened during the process of user creation response data save.")
)

var (
	// GetSessionsRepository retrieves instance of the sessions repository, performing initial creation if needed.
	GetSessionsRepository = sync.OnceValue[SessionsRepository](createSessionsRepository)

	// GetMessagesRepository retrieves instance of the messages repository, performing initial creation if needed.
	GetMessagesRepository = sync.OnceValue[MessagesRepository](createMessagesRepository)

	// GetUsersRepository retrieves instance of the users repository, performing initial creation if needed.
	GetUsersRepository = sync.OnceValue[UsersRepository](createUsersRepository)
)

// SessionsRepository represents sessions entity repository.
type SessionsRepository interface {
	Insert(name string, issuer int64) error
	GetByIssuer(issuer int64) ([]*entity.SessionEntity, error)
}

// sessionsRepositoryImpl represents implementation of SessionsRepository.
type sessionsRepositoryImpl struct{}

// Insert inserts new sessions entity to the storage.
func (w *sessionsRepositoryImpl) Insert(name string, issuer int64) error {
	instance := db.GetInstance()

	err := instance.Create(
		&entity.SessionEntity{
			Name:   name,
			Issuer: issuer}).Error

	return errors.Wrap(err, ErrPersistingSessions.Error())
}

// GetByIssuer retrieves all available sessions.
func (w *sessionsRepositoryImpl) GetByIssuer(issuer int64) ([]*entity.SessionEntity, error) {
	instance := db.GetInstance()

	var result []*entity.SessionEntity

	err := instance.Table((&entity.SessionEntity{}).TableName()).
		Where("issuer = ?", issuer).
		Find(&result).Error

	return result, err
}

// createSessionsRepository initializes sessionsRepositoryImpl.
func createSessionsRepository() SessionsRepository {
	return new(sessionsRepositoryImpl)
}

// MessagesRepository represents messages entity repository.
type MessagesRepository interface {
	Insert(issuer int64, content string) error
	GetByIssuer(issuer int64) ([]*entity.MessageEntity, error)
}

// messagesRepositoryImpl represents implementation of MessagesRepository.
type messagesRepositoryImpl struct{}

// Insert inserts new messages entity to the storage.
func (w *messagesRepositoryImpl) Insert(issuer int64, content string) error {
	instance := db.GetInstance()

	err := instance.Create(
		&entity.MessageEntity{
			Issuer:  issuer,
			Content: content}).Error

	return errors.Wrap(err, ErrPersistingMessages.Error())
}

// GetAll retrieves all available sessions.
func (w *messagesRepositoryImpl) GetByIssuer(issuer int64) ([]*entity.MessageEntity, error) {
	instance := db.GetInstance()

	var result []*entity.MessageEntity

	err := instance.Table((&entity.MessageEntity{}).TableName()).
		Where("issuer = ?", issuer).
		Find(&result).Error

	return result, err
}

// createMessagesRepository initializes messagesRepositoryImpl.
func createMessagesRepository() MessagesRepository {
	return new(messagesRepositoryImpl)
}

// UsersRepository represents users entity repository.
type UsersRepository interface {
	Insert(nam string) error
	GetByName(name string) (*entity.UserEntity, bool, error)
}

// usersRepositoryImpl represents implementation of UsersRepository.
type usersRepositoryImpl struct{}

// Insert inserts users entity to the storage.
func (w *usersRepositoryImpl) Insert(name string) error {
	instance := db.GetInstance()

	err := instance.Create(&entity.UserEntity{
		Name: name,
	}).Error

	return errors.Wrap(err, ErrPersistingUsers.Error())
}

// GetByName checks if any user with the given name exists.
func (w *usersRepositoryImpl) GetByName(name string) (*entity.UserEntity, bool, error) {
	instance := db.GetInstance()

	var result *entity.UserEntity

	err := instance.Model(&entity.UserEntity{}).
		Where("name = ?", name).
		First(&result).Error

	if err != gorm.ErrRecordNotFound {
		return result, true, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}

	return nil, false, err
}

// createUsersRepository initializes usersRepositoryImpl.
func createUsersRepository() UsersRepository {
	return new(usersRepositoryImpl)
}
