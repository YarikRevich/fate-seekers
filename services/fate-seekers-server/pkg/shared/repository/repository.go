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
	ErrPersistingSessions    = errors.New("err happened during the process of session creation response data save.")
	ErrPersistingGenerations = errors.New("err happened during the process of generation creation response data save.")
	ErrPersistingLobbies     = errors.New("err happened during the process of lobby creation response data save.")
	ErrPersistingMessages    = errors.New("err happened during the process of message creation response data save.")
	ErrPersistingUsers       = errors.New("err happened during the process of user creation response data save.")
)

var (
	// GetSessionsRepository retrieves instance of the sessions repository, performing initial creation if needed.
	GetSessionsRepository = sync.OnceValue[SessionsRepository](createSessionsRepository)

	// GetGenerationRepository retrieves instance of the generations repository, performing initial creation if needed.
	GetGenerationRepository = sync.OnceValue[GenerationsRepository](createGenerationsRepository)

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
	GetByID(id int64) (*entity.SessionEntity, bool, error)
	GetByIssuer(issuer int64) ([]*entity.SessionEntity, error)
	GetByName(name string) (*entity.SessionEntity, bool, error)
	ExistsByName(name string) (bool, error)
}

// sessionsRepositoryImpl represents implementation of SessionsRepository.
type sessionsRepositoryImpl struct {
	// Represents mutex used for database session repository related operations.
	mu sync.RWMutex
}

// InsertOrUpdate inserts or updates new sessions entity to the storage or updates existing ones.
func (w *sessionsRepositoryImpl) InsertOrUpdate(request dto.SessionsRepositoryInsertOrUpdateRequest) error {
	w.mu.Lock()

	instance := db.GetInstance()

	err := instance.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "name"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"started",
		}),
	}).Create(&entity.SessionEntity{
		Name:    request.Name,
		Seed:    request.Seed,
		Issuer:  request.Issuer,
		Started: request.Started,
	}).Error

	if err != nil {
		w.mu.Unlock()

		return errors.Wrap(err, ErrPersistingSessions.Error())
	}

	w.mu.Unlock()

	return nil
}

// DeleteByID deletes session by the provided id.
func (w *sessionsRepositoryImpl) DeleteByID(id int64) error {
	w.mu.Lock()

	instance := db.GetInstance()

	err := instance.Table((&entity.SessionEntity{}).TableName()).
		Where("id = ?", id).
		Delete(&entity.SessionEntity{}).Error

	w.mu.Unlock()

	return err
}

// GetByID retrieves a session for the provided id.
func (w *sessionsRepositoryImpl) GetByID(id int64) (*entity.SessionEntity, bool, error) {
	w.mu.RLock()

	instance := db.GetInstance()

	var result *entity.SessionEntity

	err := instance.Table((&entity.SessionEntity{}).TableName()).
		Preload((&entity.UserEntity{}).TableView()).
		Where("id = ?", id).
		First(&result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.mu.RUnlock()

			return result, false, nil
		}

		w.mu.RUnlock()

		return result, false, err
	}

	w.mu.RUnlock()

	return result, true, nil
}

// GetByIssuer retrieves all available sessions for the provided issuer.
func (w *sessionsRepositoryImpl) GetByIssuer(issuer int64) ([]*entity.SessionEntity, error) {
	w.mu.RLock()

	instance := db.GetInstance()

	var result []*entity.SessionEntity

	err := instance.Table((&entity.SessionEntity{}).TableName()).
		Preload((&entity.UserEntity{}).TableView()).
		Where("issuer = ?", issuer).
		Find(&result).Error

	w.mu.RUnlock()

	return result, err
}

// GetByName retrieves available session for the provided name.
func (w *sessionsRepositoryImpl) GetByName(name string) (*entity.SessionEntity, bool, error) {
	w.mu.RLock()

	instance := db.GetInstance()

	var result *entity.SessionEntity

	err := instance.Table((&entity.SessionEntity{}).TableName()).
		Preload((&entity.UserEntity{}).TableView()).
		Where("name = ?", name).
		First(&result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.mu.RUnlock()

			return result, false, nil
		}

		w.mu.RUnlock()

		return result, false, err
	}

	w.mu.RUnlock()

	return result, true, nil
}

// ExistsByName checks if session exists for the provided name.
func (w *sessionsRepositoryImpl) ExistsByName(name string) (bool, error) {
	w.mu.RLock()

	instance := db.GetInstance()

	var result *entity.SessionEntity

	err := instance.Table((&entity.SessionEntity{}).TableName()).
		Preload((&entity.UserEntity{}).TableView()).
		Where("name = ?", name).
		First(&result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.mu.RUnlock()

			return false, nil
		}

		w.mu.RUnlock()

		return false, err
	}

	w.mu.RUnlock()

	return true, nil
}

// createSessionsRepository initializes sessionsRepositoryImpl.
func createSessionsRepository() SessionsRepository {
	return new(sessionsRepositoryImpl)
}

// GenerationsRepository represents generations entity repository.
type GenerationsRepository interface {
	InsertOrUpdate(request dto.GenerationsRepositoryInsertOrUpdateRequest) error
	GetBySessionID(sessionID int64) ([]*entity.GenerationsEntity, error)
}

// generationsRepositoryImpl represents implementation of GenerationsRepository.
type generationsRepositoryImpl struct {
	// Represents mutex used for database generation repository related operations.
	mu sync.RWMutex
}

// InsertOrUpdate inserts new generations entity to the storage or updates existing ones.
func (w *generationsRepositoryImpl) InsertOrUpdate(request dto.GenerationsRepositoryInsertOrUpdateRequest) error {
	w.mu.Lock()

	instance := db.GetInstance()

	err := instance.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"active",
		}),
	}).Create(&entity.GenerationsEntity{
		ID:        request.ID,
		SessionID: request.SessionID,
		Name:      request.Name,
		Type:      request.Type,
		Active:    request.Active,
	}).Error

	if err != nil {
		w.mu.Unlock()

		return errors.Wrap(err, ErrPersistingGenerations.Error())
	}

	w.mu.Unlock()

	return nil
}

// GetBySessionID retrieves all available generations for the provided session id.
func (w *generationsRepositoryImpl) GetBySessionID(sessionID int64) ([]*entity.GenerationsEntity, error) {
	w.mu.RLock()

	instance := db.GetInstance()

	var result []*entity.GenerationsEntity

	err := instance.Table((&entity.GenerationsEntity{}).TableName()).
		Where("session_id = ?", sessionID).
		Find(&result).Error

	w.mu.RUnlock()

	return result, err
}

// createGenerationsRepository initializes generationsRepositoryImpl.
func createGenerationsRepository() GenerationsRepository {
	return new(generationsRepositoryImpl)
}

// LobbiesRepository represents lobbies entity repository.
type LobbiesRepository interface {
	InsertOrUpdate(request dto.LobbiesRepositoryInsertOrUpdateRequest) error
	InsertOrUpdateWithTransaction(transaction *gorm.DB, request dto.LobbiesRepositoryInsertOrUpdateRequest) error
	DeleteByUserIDAndSessionID(userID, sessionID int64) error
	DeleteByUserIDAndSessionIDWithTransaction(transaction *gorm.DB, userID, sessionID int64) error
	GetByUserID(userID int64) ([]*entity.LobbyEntity, bool, error)
	GetBySessionID(sessionID int64) ([]*entity.LobbyEntity, bool, error)
	Lock()
	Unlock()
}

// lobbiesRepositoryImpl represents implementation of LobbiesRepository.
type lobbiesRepositoryImpl struct {
	// Represents internal mutex used for database lobbies repository related operations.
	mu sync.RWMutex

	// Represents external mutex used for database management.
	externalLock sync.Mutex
}

// insertOrUpdate inserts new lobbies entity to the storage or updates existing ones with the provided db instance.
func (w *lobbiesRepositoryImpl) insertOrUpdate(instance *gorm.DB, request dto.LobbiesRepositoryInsertOrUpdateRequest) error {
	w.mu.Lock()

	err := instance.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "session_id"},
			{Name: "skin"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"health",
			"active",
			"eliminated",
			"position_x",
			"position_y",
			"host",
		}),
	}).Create(&entity.LobbyEntity{
		UserID:     request.UserID,
		SessionID:  request.SessionID,
		Skin:       int64(request.Skin),
		Health:     int64(request.Health),
		Active:     request.Active,
		Host:       request.Host,
		Eliminated: request.Eliminated,
		PositionX:  request.PositionX,
		PositionY:  request.PositionY,
	}).Error

	if err != nil {
		w.mu.Unlock()

		return errors.Wrap(err, ErrPersistingLobbies.Error())
	}

	w.mu.Unlock()

	return nil
}

// InsertOrUpdate inserts new lobbies entity to the storage or updates existing ones.
func (w *lobbiesRepositoryImpl) InsertOrUpdate(request dto.LobbiesRepositoryInsertOrUpdateRequest) error {
	return w.insertOrUpdate(db.GetInstance(), request)
}

// InsertOrUpdateWithTransaction inserts new lobbies entity to the storage or updates existing ones with provided transaction.
func (w *lobbiesRepositoryImpl) InsertOrUpdateWithTransaction(transaction *gorm.DB, request dto.LobbiesRepositoryInsertOrUpdateRequest) error {
	return w.insertOrUpdate(transaction, request)
}

// deleteByUserIDAndSessionID deletes lobby by the provided user id with provided db instance.
func (w *lobbiesRepositoryImpl) deleteByUserIDAndSessionID(instance *gorm.DB, userID, sessionID int64) error {
	w.mu.Lock()

	err := instance.Table((&entity.LobbyEntity{}).TableName()).
		Where("user_id = ? AND session_id = ?", userID, sessionID).
		Delete(&entity.LobbyEntity{}).Error

	w.mu.Unlock()

	return err
}

// DeleteByUserIDAndSessionID deletes lobby by the provided user id.
func (w *lobbiesRepositoryImpl) DeleteByUserIDAndSessionID(userID, sessionID int64) error {
	return w.deleteByUserIDAndSessionID(db.GetInstance(), userID, sessionID)
}

// DeleteByUserIDAndSessionIDWithTransaction deletes lobby by the provided user id with provided transaction.
func (w *lobbiesRepositoryImpl) DeleteByUserIDAndSessionIDWithTransaction(transaction *gorm.DB, userID, sessionID int64) error {
	return w.deleteByUserIDAndSessionID(transaction, userID, sessionID)
}

// GetByUserID retrieves lobby by the provided user id.
func (w *lobbiesRepositoryImpl) GetByUserID(userID int64) ([]*entity.LobbyEntity, bool, error) {
	w.mu.RLock()

	instance := db.GetInstance()

	var result []*entity.LobbyEntity

	err := instance.Table((&entity.LobbyEntity{}).TableName()).
		Preload((&entity.UserEntity{}).TableView()).
		Preload((&entity.SessionEntity{}).TableView()).
		Where("user_id = ?", userID).
		Find(&result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.mu.RUnlock()

			return nil, false, nil
		}

		w.mu.RUnlock()

		return nil, false, err
	}

	if len(result) == 0 {
		w.mu.RUnlock()

		return nil, false, nil
	}

	w.mu.RUnlock()

	return result, true, nil
}

// GetBySessionID retrieves lobby by the provided session id.
func (w *lobbiesRepositoryImpl) GetBySessionID(sessionID int64) ([]*entity.LobbyEntity, bool, error) {
	w.mu.RLock()

	instance := db.GetInstance()

	var result []*entity.LobbyEntity

	err := instance.Table((&entity.LobbyEntity{}).TableName()).
		Preload((&entity.UserEntity{}).TableView()).
		Preload((&entity.SessionEntity{}).TableView()).
		Where("session_id = ?", sessionID).
		Find(&result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.mu.RUnlock()

			return nil, false, nil
		}

		w.mu.RUnlock()

		return nil, false, err
	}

	if len(result) == 0 {
		w.mu.RUnlock()

		return nil, false, nil
	}

	w.mu.RUnlock()

	return result, true, nil
}

// Lock locks external lock mutex
func (w *lobbiesRepositoryImpl) Lock() {
	w.externalLock.Lock()
}

// Unlock unlocks external lock mutex
func (w *lobbiesRepositoryImpl) Unlock() {
	w.externalLock.Unlock()
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
type messagesRepositoryImpl struct {
	// Represents mutex used for database messages repository related operations.
	mu sync.RWMutex
}

// Insert inserts new messages entity to the storage.
func (w *messagesRepositoryImpl) Insert(issuer int64, content string) error {
	w.mu.Lock()

	instance := db.GetInstance()

	err := instance.Create(
		&entity.MessageEntity{
			Issuer:  issuer,
			Content: content}).Error

	if err != nil {
		w.mu.Unlock()

		return errors.Wrap(err, ErrPersistingMessages.Error())
	}

	w.mu.Unlock()

	return nil
}

// GetAll retrieves all available sessions.
func (w *messagesRepositoryImpl) GetByIssuer(issuer int64) ([]*entity.MessageEntity, error) {
	w.mu.RLock()

	instance := db.GetInstance()

	var result []*entity.MessageEntity

	err := instance.Table((&entity.MessageEntity{}).TableName()).
		Preload((&entity.UserEntity{}).TableView()).
		Where("issuer = ?", issuer).
		Find(&result).Error

	w.mu.RUnlock()

	return result, err
}

// createMessagesRepository initializes messagesRepositoryImpl.
func createMessagesRepository() MessagesRepository {
	return new(messagesRepositoryImpl)
}

// UsersRepository represents users entity repository.
type UsersRepository interface {
	Insert(name string) error
	ExistsByName(name string) (bool, error)
	GetByName(name string) (*entity.UserEntity, bool, error)
}

// usersRepositoryImpl represents implementation of UsersRepository.
type usersRepositoryImpl struct {
	// Represents mutex used for database users repository related operations.
	mu sync.RWMutex
}

// Insert inserts users entity to the storage.
func (w *usersRepositoryImpl) Insert(name string) error {
	w.mu.Lock()

	instance := db.GetInstance()

	err := instance.Create(&entity.UserEntity{
		Name: name,
	}).Error

	if err != nil {
		w.mu.Unlock()

		return errors.Wrap(err, ErrPersistingUsers.Error())
	}

	w.mu.Unlock()

	return nil
}

// ExistsByName checks if user with the given name exists.
func (w *usersRepositoryImpl) ExistsByName(name string) (bool, error) {
	w.mu.RLock()

	instance := db.GetInstance()

	err := instance.Table((&entity.UserEntity{}).TableName()).
		Where("name = ?", name).
		First(&entity.UserEntity{}).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.mu.RUnlock()

			return false, nil
		}

		w.mu.RUnlock()

		return false, err
	}

	w.mu.RUnlock()

	return true, nil
}

// GetByName retrieves user with the given name.
func (w *usersRepositoryImpl) GetByName(name string) (*entity.UserEntity, bool, error) {
	w.mu.RLock()

	instance := db.GetInstance()

	var result *entity.UserEntity

	err := instance.Model(&entity.UserEntity{}).
		Where("name = ?", name).
		First(&result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.mu.RUnlock()

			return nil, false, nil
		}

		w.mu.RUnlock()

		return nil, false, err
	}

	w.mu.RUnlock()

	return result, true, nil
}

// createUsersRepository initializes usersRepositoryImpl.
func createUsersRepository() UsersRepository {
	return new(usersRepositoryImpl)
}
