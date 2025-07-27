package cache

import (
	"errors"
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/logging"
	lru "github.com/hashicorp/golang-lru/v2"
)

var (
	// GetInstance retrieves instance of the networking cache, performing initilization if needed.
	GetInstance = sync.OnceValue[*NetworkingCache](newNetworkingCache)
)

var (
	ErrFailedToEvictSessionCacheEntity  = errors.New("err failed to evict session cache entity")
	ErrFailedToEvictLobbySetCacheEntity = errors.New("err failed to evict lobby set cache entity")
	ErrFailedToEvictMetadataCacheEntity = errors.New("err failed to evict metadata cache entity")
)

const (
	// Represents max amount of users per session.
	maxSessionUsers = 8
)

// lobby -> lobby set -> metadata

// NetworkingCache represents networking cache.
type NetworkingCache struct {
	// Represents sessions cache instance.
	sessions *lru.Cache[string, []dto.CacheSessionEntity]

	// Represents lobby sets cache instance. Value contains issuer names only.
	lobbySets *lru.Cache[int64, []string]

	// Represents metadata cache instance.
	metadata *lru.Cache[string, dto.CacheMetadataEntity]

	// Represents expirable messages cache, which contains offset for the message table.
	// If user stops request messages, all the messages would be retrieved.
	messages *lru.Cache[string, int]

	// Represents users cache instance.
	users *lru.Cache[string, int64]
}

// AddSession adds session cache instance with the provided key and value.
func (nc *NetworkingCache) AddSessions(key string, value []dto.CacheSessionEntity) {
	nc.sessions.Add(key, value)
}

// GetSession retrieves session cache instance by the provided key.
func (nc *NetworkingCache) GetSessions(key string) ([]dto.CacheSessionEntity, bool) {
	return nc.sessions.Get(key)
}

// EvictSessions evicts sessions cache for the provided key.
func (nc *NetworkingCache) EvictSessions(key string) {
	if ok := nc.sessions.Remove(key); !ok {
		logging.GetInstance().Error(ErrFailedToEvictSessionCacheEntity.Error())
	}
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
	if ok := nc.lobbySets.Remove(key); !ok {
		logging.GetInstance().Error(ErrFailedToEvictLobbySetCacheEntity.Error())
	}
}

// AddMetadata adds metadata cache instance with the provided key and value.
func (nc *NetworkingCache) AddMetadata(key string, value dto.CacheMetadataEntity) {
	nc.metadata.Add(key, value)
}

// GetMetadata retrieves metadata cache instance by the provided key.
func (nc *NetworkingCache) GetMetadata(key string) (dto.CacheMetadataEntity, bool) {
	return nc.metadata.Get(key)
}

// GetMetadataMappings retrieves all metadata mapping cache instances.
func (nc *NetworkingCache) GetMetadataMappings() map[string]dto.CacheMetadataEntity {
	result := make(map[string]dto.CacheMetadataEntity)

	for _, key := range nc.metadata.Keys() {
		value, _ := nc.GetMetadata(key)

		result[key] = value
	}

	return result
}

// EvictMetadata evicts metadata cache for the provided key.
func (nc *NetworkingCache) EvictMetadata(key string) {
	if ok := nc.metadata.Remove(key); !ok {
		logging.GetInstance().Error(ErrFailedToEvictMetadataCacheEntity.Error())
	}
}

// AddMessage retrieves messages cache instance.
func (nc *NetworkingCache) AddMessages(key string, value int) {
	nc.messages.Add(key, value)
}

// GetMessage retrieves messages cache instance by the provided key.
func (nc *NetworkingCache) GetMessage(key string) (int, bool) {
	return nc.messages.Get(key)
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
	sessions, err := lru.New[string, []dto.CacheSessionEntity](config.GetOperationMaxSessionsAmount())
	if err != nil {
		logging.GetInstance().Fatal(err.Error())
	}

	lobbySets, err := lru.New[int64, []string](config.GetOperationMaxSessionsAmount())
	if err != nil {
		logging.GetInstance().Fatal(err.Error())
	}

	metadata, err := lru.New[string, dto.CacheMetadataEntity](
		config.GetOperationMaxSessionsAmount() * config.GetOperationMaxSessionsAmount())
	if err != nil {
		logging.GetInstance().Fatal(err.Error())
	}

	messages, err := lru.New[string, int](config.GetOperationMaxSessionsAmount())
	if err != nil {
		logging.GetInstance().Fatal(err.Error())
	}

	users, err := lru.New[string, int64](config.GetOperationMaxSessionsAmount() * maxSessionUsers)
	if err != nil {
		logging.GetInstance().Fatal(err.Error())
	}

	return &NetworkingCache{
		sessions:  sessions,
		lobbySets: lobbySets,
		metadata:  metadata,
		messages:  messages,
		users:     users,
	}
}
