package cache

import (
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/logging"

	// LRU is utilized as a map with a limited size solution in this case.
	// Should be replaced with a native solution without external library usage.
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
	lobbySets *lru.Cache[int64, []dto.CacheLobbySetEntity]

	// Represents mutex used for lobby sets related transactions.
	lobbySetsMutex sync.Mutex

	// Represents user activity cache instance.
	userActivity *lru.Cache[string, time.Duration]

	// Represents metadata cache instance.
	metadata *lru.Cache[string, []*dto.CacheMetadataEntity]

	// Represents mutex used for metadata related transactions.
	metadataMutex sync.Mutex

	// Represents users cache instance.
	users *lru.Cache[string, int64]

	// Represents generated chests cache instance.
	generatedChests *lru.Cache[string, []*dto.CacheGeneratedChestEntity]

	// Represents mutex used for generated chests related transactions.
	generatedChestsMutex sync.Mutex

	// Represents generated health packs cache instance.
	generatedHealthPacks *lru.Cache[string, []*dto.CacheGeneratedHealthPacksEntity]

	// Represents mutex used for generated health packs related transactions.
	generatedHealthPacksMutex sync.Mutex
}

// BeginSessionsTransaction begins sessions cache instance transaction.
func (nc *NetworkingCache) BeginSessionsTransaction() {
	nc.sessionsMutex.Lock()
}

// CommitSessionsTransaction commits sessions cache instance transaction.
func (nc *NetworkingCache) CommitSessionsTransaction() {
	nc.sessionsMutex.Unlock()
}

// AddSessions adds session cache instance with the provided key and value.
func (nc *NetworkingCache) AddSessions(key int64, value dto.CacheSessionEntity) {
	nc.sessions.Add(key, value)
}

// GetSessions retrieves session cache instance by the provided key.
func (nc *NetworkingCache) GetSessions(key int64) (dto.CacheSessionEntity, bool) {
	return nc.sessions.Get(key)
}

// GetSessionsMappings retrieves all sessions mapping cache instances.
func (nc *NetworkingCache) GetSessionsMappings() map[int64]dto.CacheSessionEntity {
	result := make(map[int64]dto.CacheSessionEntity)

	for _, key := range nc.sessions.Keys() {
		value, _ := nc.GetSessions(key)

		result[key] = value
	}

	return result
}

// EvictSessions evicts sessions cache for the provided key.
func (nc *NetworkingCache) EvictSessions(key int64) {
	nc.sessions.Remove(key)
}

// EvictSessionsByName evicts sessions cache for the provided name value.
func (nc *NetworkingCache) EvictSessionsByName(name string) {
	for _, key := range nc.sessions.Keys() {
		value, _ := nc.GetSessions(key)

		if value.Name == name {
			nc.sessions.Remove(key)
		}
	}
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
func (nc *NetworkingCache) AddLobbySet(key int64, value []dto.CacheLobbySetEntity) {
	nc.lobbySets.Add(key, value)
}

// GetLobbies retrieves lobby cache instance by the provided key.
func (nc *NetworkingCache) GetLobbySet(key int64) ([]dto.CacheLobbySetEntity, bool) {
	return nc.lobbySets.Get(key)
}

// GetLobbySetMappings retrieves all lobby set mapping cache instances.
func (nc *NetworkingCache) GetLobbySetMappings() map[int64][]dto.CacheLobbySetEntity {
	result := make(map[int64][]dto.CacheLobbySetEntity)

	for _, key := range nc.lobbySets.Keys() {
		value, _ := nc.GetLobbySet(key)

		result[key] = value
	}

	return result
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
func (nc *NetworkingCache) AddMetadata(key string, value []*dto.CacheMetadataEntity) {
	nc.metadata.Add(key, value)
}

// GetMetadata retrieves metadata cache instance by the provided key.
func (nc *NetworkingCache) GetMetadata(key string) ([]*dto.CacheMetadataEntity, bool) {
	return nc.metadata.Get(key)
}

// GetMetadataMappings retrieves all metadata mapping cache instances.
func (nc *NetworkingCache) GetMetadataMappings() map[string][]*dto.CacheMetadataEntity {
	result := make(map[string][]*dto.CacheMetadataEntity)

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

// AddUser adds users cache instance with the provided key and value.
func (nc *NetworkingCache) AddUser(key string, value int64) {
	nc.users.Add(key, value)
}

// GetUsers retrieves users cache instance by the provided key.
func (nc *NetworkingCache) GetUsers(key string) (int64, bool) {
	return nc.users.Get(key)
}

// BeginGeneratedChestsTransaction begins generated chests cache instance transaction.
func (nc *NetworkingCache) BeginGeneratedChestsTransaction() {
	nc.generatedChestsMutex.Lock()
}

// CommitGeneratedChestsTransaction commits generated chests cache instance transaction.
func (nc *NetworkingCache) CommitGeneratedChestsTransaction() {
	nc.generatedChestsMutex.Unlock()
}

// AddGeneratedChests adds generated chests cache instance.
func (nc *NetworkingCache) AddGeneratedChests(key string, value []*dto.CacheGeneratedChestEntity) {
	nc.generatedChests.Add(key, value)
}

// GetGeneratedChests retrieves generated chests cache instance by the provided key.
func (nc *NetworkingCache) GetGeneratedChests(key string) ([]*dto.CacheGeneratedChestEntity, bool) {
	return nc.generatedChests.Get(key)
}

// EvictGeneratedChests evicts generated chests cache for the provided key.
func (nc *NetworkingCache) EvictGeneratedChests(key string) {
	nc.generatedChests.Remove(key)
}

// BeginGeneratedHealthPacksTransaction begins generated health packs cache instance transaction.
func (nc *NetworkingCache) BeginGeneratedHealthPacksTransaction() {
	nc.generatedHealthPacksMutex.Lock()
}

// CommitGeneratedHealthPacksTransaction commits generated health packs cache instance transaction.
func (nc *NetworkingCache) CommitGeneratedHealthPacksTransaction() {
	nc.generatedHealthPacksMutex.Unlock()
}

// AddGeneratedHealthPacks adds generated health packs cache instance.
func (nc *NetworkingCache) AddGeneratedHealthPacks(key string, value []*dto.CacheGeneratedHealthPacksEntity) {
	nc.generatedHealthPacks.Add(key, value)
}

// GetGeneratedHealthPacks retrieves generated health packs cache instance by the provided key.
func (nc *NetworkingCache) GetGeneratedHealthPacks(key string) ([]*dto.CacheGeneratedHealthPacksEntity, bool) {
	return nc.generatedHealthPacks.Get(key)
}

// EvictGeneratedHealthPacks evicts generated health packs cache for the provided key.
func (nc *NetworkingCache) EvictGeneratedHealthPacks(key string) {
	nc.generatedHealthPacks.Remove(key)
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

	lobbySets, err := lru.New[int64, []dto.CacheLobbySetEntity](config.GetOperationMaxSessionsAmount())
	if err != nil {
		logging.GetInstance().Fatal(err.Error())
	}

	userActivity, err := lru.New[string, time.Duration](
		config.GetOperationMaxSessionsAmount() * config.MAX_SESSION_USERS)
	if err != nil {
		logging.GetInstance().Fatal(err.Error())
	}

	metadata, err := lru.New[string, []*dto.CacheMetadataEntity](
		config.GetOperationMaxSessionsAmount() * config.MAX_SESSION_USERS)
	if err != nil {
		logging.GetInstance().Fatal(err.Error())
	}

	users, err := lru.New[string, int64](config.GetOperationMaxSessionsAmount() * config.MAX_SESSION_USERS)
	if err != nil {
		logging.GetInstance().Fatal(err.Error())
	}

	generatedChests, err := lru.New[string, []*dto.CacheGeneratedChestEntity](
		config.GetOperationMaxSessionsAmount() * config.GetOperationMaxChestsAmount())
	if err != nil {
		logging.GetInstance().Fatal(err.Error())
	}

	generatedHealthPacks, err := lru.New[string, []*dto.CacheGeneratedHealthPacksEntity](
		config.GetOperationMaxSessionsAmount() * config.GetOperationMaxHealthPacksAmount())
	if err != nil {
		logging.GetInstance().Fatal(err.Error())
	}

	return &NetworkingCache{
		sessions:             sessions,
		userSessions:         userSessions,
		lobbySets:            lobbySets,
		userActivity:         userActivity,
		metadata:             metadata,
		users:                users,
		generatedChests:      generatedChests,
		generatedHealthPacks: generatedHealthPacks,
	}
}
