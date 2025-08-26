package repository

import (
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/db"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/entity"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrPersistingSessions = errors.New("err happened during the process of session creation response data save.")
	ErrPersistingLobbies  = errors.New("err happened during the process of lobby creation response data save.")
	ErrPersistingMessages = errors.New("err happened during the process of message creation response data save.")
	ErrPersistingUsers    = errors.New("err happened during the process of user creation response data save.")
)

var (
	// GetSessionsRepository retrieves instance of the sessions repository, performing initial creation if needed.
	GetSessionsRepository = sync.OnceValue[SessionsRepository](createSessionsRepository)

	// GetLobbiesRepository retrieves instance of the lobbies repository, performing initial creation if needed.
	GetLobbiesRepository = sync.OnceValue[LobbiesRepository](createLobbiesRepository)

	// GetMessagesRepository retrieves instance of the messages repository, performing initial creation if needed.
	GetMessagesRepository = sync.OnceValue[MessagesRepository](createMessagesRepository)

	// GetUsersRepository retrieves instance of the users repository, performing initial creation if needed.
	GetUsersRepository = sync.OnceValue[UsersRepository](createUsersRepository)
)

// SessionsRepository represents sessions entity repository.
type SessionsRepository interface {
	InsertOrUpdate(request dto.SessionsRepositoryInsertOrUpdateRequest) error
	DeleteByID(id int64) error
	GetByID(id int64) (*entity.SessionEntity, error)
	GetByIssuer(issuer int64) ([]*entity.SessionEntity, error)
}

// sessionsRepositoryImpl represents implementation of SessionsRepository.
type sessionsRepositoryImpl struct{}

// Insert inserts new sessions entity to the storage or updates existing ones.
func (w *sessionsRepositoryImpl) InsertOrUpdate(request dto.SessionsRepositoryInsertOrUpdateRequest) error {
	instance := db.GetInstance()

	err := instance.Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{
			"started",
		}),
	}).Create(&entity.SessionEntity{
		Name:    request.Name,
		Seed:    request.Seed,
		Issuer:  request.Issuer,
		Started: request.Started,
	}).Error

	return errors.Wrap(err, ErrPersistingSessions.Error())
}

// DeleteByID deletes session by the provided id.
func (w *sessionsRepositoryImpl) DeleteByID(id int64) error {
	instance := db.GetInstance()

	return instance.Table((&entity.SessionEntity{}).TableName()).
		Where("id = ?", id).
		Delete(&entity.SessionEntity{}).Error
}

// GetByID retrieves a session for the provided id.
func (w *sessionsRepositoryImpl) GetByID(id int64) (*entity.SessionEntity, error) {
	instance := db.GetInstance()

	var result *entity.SessionEntity

	err := instance.Table((&entity.SessionEntity{}).TableName()).
		Where("id = ?", id).
		Find(&result).Error

	return result, err
}

// GetByIssuer retrieves all available sessions for the provided issuer.
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

// LobbiesRepository represents lobbies entity repository.
type LobbiesRepository interface {
	InsertOrUpdate(request dto.LobbiesRepositoryInsertOrUpdateRequest) error
	DeleteByUserID(userID int64) error
	GetByUserID(userID int64) ([]*entity.LobbyEntity, bool, error)
	GetBySessionID(sessionID int64) ([]*entity.LobbyEntity, bool, error)
}

// lobbiesRepositoryImpl represents implementation of LobbiesRepository.
type lobbiesRepositoryImpl struct{}

// Insert inserts new lobbies entity to the storage or updates existing ones.
func (w *lobbiesRepositoryImpl) InsertOrUpdate(request dto.LobbiesRepositoryInsertOrUpdateRequest) error {
	instance := db.GetInstance()

	err := instance.Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{
			"health",
			"active",
			"eliminated",
			"position_x",
			"position_y",
		}),
	}).Create(&entity.LobbyEntity{
		UserID:     request.UserID,
		SessionID:  request.SessionID,
		Skin:       int64(request.Skin),
		Health:     int64(request.Health),
		Active:     request.Active,
		Eliminated: request.Eliminated,
		PositionX:  request.PositionX,
		PositionY:  request.PositionY,
	}).Error

	return errors.Wrap(err, ErrPersistingLobbies.Error())
}

// DeleteByUserID deletes lobby by the provided user id.
func (w *lobbiesRepositoryImpl) DeleteByUserID(userID int64) error {
	instance := db.GetInstance()

	return instance.Table((&entity.LobbyEntity{}).TableName()).
		Where("user_id = ?", userID).
		Delete(&entity.LobbyEntity{}).Error
}

// GetByUserID retrieves lobby by the provided user id.
func (w *lobbiesRepositoryImpl) GetByUserID(userID int64) ([]*entity.LobbyEntity, bool, error) {
	instance := db.GetInstance()

	var result []*entity.LobbyEntity

	err := instance.Table((&entity.LobbyEntity{}).TableName()).
		Where("user_id = ?", userID).
		Find(&result).Error

	if err != gorm.ErrRecordNotFound {
		return result, true, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}

	return nil, false, err
}

// GetBySessionID retrieves lobby by the provided session id.
func (w *lobbiesRepositoryImpl) GetBySessionID(sessionID int64) ([]*entity.LobbyEntity, bool, error) {
	instance := db.GetInstance()

	var result []*entity.LobbyEntity

	err := instance.Table((&entity.LobbyEntity{}).TableName()).
		Where("session_id = ?", sessionID).
		Find(&result).Error

	if err != gorm.ErrRecordNotFound {
		return result, true, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}

	return nil, false, err
}

// createLobbiesRepository initializes lobbiesRepositoryImpl.
func createLobbiesRepository() LobbiesRepository {
	return new(lobbiesRepositoryImpl)
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
	Insert(name string) error
	Exists(name string) (bool, error)
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

// Exists checks if user with the given name exists.
func (w *usersRepositoryImpl) Exists(name string) (bool, error) {
	instance := db.GetInstance()

	var exists bool

	err := instance.Model(&entity.UserEntity{}).
		Select("count(*) > 0").
		Where("name = ?", name).
		Find(&exists).Error

	if err != nil {
		return false, err
	}

	return exists, err
}

// GetByName retrieves user with the given name.
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
