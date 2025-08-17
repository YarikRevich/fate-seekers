package cache

import (
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/logging"
	lru "github.com/hashicorp/golang-lru/v2"
)

var (
	// GetInstance retrieves instance of the networking cache, performing initilization if needed.
	GetInstance = sync.OnceValue[*NetworkingCache](newNetworkingCache)
)

// NetworkingCache represents networking cache.
type NetworkingCache struct {
	// Represents sessions cache instance.
	sessions *lru.Cache[int64, dto.CacheSessionEntity]

	// Represents mutex sessions related transactions.
	sessionsMutex sync.Mutex

	// Represents user sessions cache instance.
	userSessions *lru.Cache[string, []dto.CacheSessionEntity]

	// Represents mutex used for user sessions related transactions.
	userSessionsMutex sync.Mutex

	// Represents lobby sets cache instance. Value contains issuer names only.
	lobbySets *lru.Cache[int64, []string]

	// Represents mutex used for lobby sets related transactions.
	lobbySetsMutex sync.Mutex

	// Represents user activity cache instance.
	userActivity *lru.Cache[string, time.Duration]

	// Represents metadata cache instance.
	metadata *lru.Cache[string, []dto.CacheMetadataEntity]

	// Represents mutex used for metadata related transactions.
	metadataMutex sync.Mutex

	// Represents expirable messages cache, which contains offset for the message table.
	// If user stops request messages, all the messages would be retrieved.
	messages *lru.Cache[string, int]

	// Represents mutex used for messages related transactions.
	messagesMutex sync.Mutex

	// Represents users cache instance.
	users *lru.Cache[string, int64]
}

// BeginUserSessionsTransaction begins user sessions cache instance transaction.
func (nc *NetworkingCache) BeginUserSessionsTransaction() {
	nc.userSessionsMutex.Lock()
}

// CommitUserSessionsTransaction commits user sessions cache instance transaction.
func (nc *NetworkingCache) CommitUserSessionsTransaction() {
	nc.userSessionsMutex.Unlock()
}

// AddUserSessions adds user session cache instance with the provided key and value.
func (nc *NetworkingCache) AddUserSessions(key string, value []dto.CacheSessionEntity) {
	nc.userSessions.Add(key, value)
}

// GetUserSessions retrieves user session cache instance by the provided key.
func (nc *NetworkingCache) GetUserSessions(key string) ([]dto.CacheSessionEntity, bool) {
	return nc.userSessions.Get(key)
}

// EvictUserSessions evicts user sessions cache for the provided key.
func (nc *NetworkingCache) EvictUserSessions(key string) {
	nc.userSessions.Remove(key)
}

// BeginLobbySetTransaction begins lobby set cache instance transaction.
func (nc *NetworkingCache) BeginLobbySetTransaction() {
	nc.lobbySetsMutex.Lock()
}

// CommitLobbySetTransaction commits lobby set cache instance transaction.
func (nc *NetworkingCache) CommitLobbySetTransaction() {
	nc.lobbySetsMutex.Unlock()
}

// AddLobbySet adds lobby set cache instance with the provided key and value.
func (nc *NetworkingCache) AddLobbySet(key int64, value []string) {
	nc.lobbySets.Add(key, value)
}

// GetLobbies retrieves lobby cache instance by the provided key.
func (nc *NetworkingCache) GetLobbySet(key int64) ([]string, bool) {
	return nc.lobbySets.Get(key)
}

// EvictLobbySet evicts lobby set cache for the provided key.
func (nc *NetworkingCache) EvictLobbySet(key int64) {
	nc.lobbySets.Remove(key)
}

// AddUserActivity adds user activity cache instance with the provided key and value.
func (nc *NetworkingCache) AddUserActivity(key string, value time.Duration) {
	nc.userActivity.Add(key, value)
}

// GetUserActivity retrieves user activity cache instance by the provided key.
func (nc *NetworkingCache) GetUserActivity(key string) (time.Duration, bool) {
	return nc.userActivity.Get(key)
}

// BeginMetadataTransaction begins metadata cache instance transaction.
func (nc *NetworkingCache) BeginMetadataTransaction() {
	nc.metadataMutex.Lock()
}

// CommitMetadataTransaction commits metadata cache instance transaction.
func (nc *NetworkingCache) CommitMetadataTransaction() {
	nc.metadataMutex.Unlock()
}

// AddMetadata adds metadata cache instance with the provided key and value.
func (nc *NetworkingCache) AddMetadata(key string, value []dto.CacheMetadataEntity) {
	nc.metadata.Add(key, value)
}

// GetMetadata retrieves metadata cache instance by the provided key.
func (nc *NetworkingCache) GetMetadata(key string) ([]dto.CacheMetadataEntity, bool) {
	return nc.metadata.Get(key)
}

// GetMetadataMappings retrieves all metadata mapping cache instances.
func (nc *NetworkingCache) GetMetadataMappings() map[string][]dto.CacheMetadataEntity {
	result := make(map[string][]dto.CacheMetadataEntity)

	for _, key := range nc.metadata.Keys() {
		value, _ := nc.GetMetadata(key)

		result[key] = value
	}

	return result
}

// EvictMetadata evicts metadata cache for the provided key.
func (nc *NetworkingCache) EvictMetadata(key string) {
	nc.metadata.Remove(key)
}

// BeginMessagesTransaction begins messages cache instance transaction.
func (nc *NetworkingCache) BeginMessagesTransaction() {
	nc.messagesMutex.Lock()
}

// CommitMessagesTransaction commits messages cache instance transaction.
func (nc *NetworkingCache) CommitMessagesTransaction() {
	nc.messagesMutex.Unlock()
}

// AddMessage retrieves messages cache instance.
func (nc *NetworkingCache) AddMessages(key string, value int) {
	nc.messages.Add(key, value)
}

// GetMessage retrieves messages cache instance by the provided key.
func (nc *NetworkingCache) GetMessage(key string) (int, bool) {
	return nc.messages.Get(key)
}

// EvictMessages evicts messages cache for the provided key.
func (nc *NetworkingCache) EvictMessages(key string) {
	nc.messages.Remove(key)
}

// AddUser adds users cache instance with the provided key and value.
func (nc *NetworkingCache) AddUser(key string, value int64) {
	nc.users.Add(key, value)
}

// GetUsers retrieves users cache instance by the provided key.
func (nc *NetworkingCache) GetUsers(key string) (int64, bool) {
	return nc.users.Get(key)
}

// newNetworkingCache initializes NetworkingCache.
func newNetworkingCache() *NetworkingCache {
	sessions, err := lru.New[int64, dto.CacheSessionEntity](config.GetOperationMaxSessionsAmount())
	if err != nil {
		logging.GetInstance().Fatal(err.Error())
	}

	userSessions, err := lru.New[string, []dto.CacheSessionEntity](config.GetOperationMaxSessionsAmount())
	if err != nil {
		logging.GetInstance().Fatal(err.Error())
	}

	lobbySets, err := lru.New[int64, []string](config.GetOperationMaxSessionsAmount())
	if err != nil {
		logging.GetInstance().Fatal(err.Error())
	}

	userActivity, err := lru.New[string, time.Duration](
		config.GetOperationMaxSessionsAmount() * config.MAX_SESSION_USERS)
	if err != nil {
		logging.GetInstance().Fatal(err.Error())
	}

	metadata, err := lru.New[string, []dto.CacheMetadataEntity](
		config.GetOperationMaxSessionsAmount() * config.MAX_SESSION_USERS)
	if err != nil {
		logging.GetInstance().Fatal(err.Error())
	}

	messages, err := lru.New[string, int](config.GetOperationMaxSessionsAmount())
	if err != nil {
		logging.GetInstance().Fatal(err.Error())
	}

	users, err := lru.New[string, int64](config.GetOperationMaxSessionsAmount() * config.MAX_SESSION_USERS)
	if err != nil {
		logging.GetInstance().Fatal(err.Error())
	}

	return &NetworkingCache{
		sessions:     sessions,
		userSessions: userSessions,
		lobbySets:    lobbySets,
		userActivity: userActivity,
		metadata:     metadata,
		messages:     messages,
		users:        users,
	}
}
