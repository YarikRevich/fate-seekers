package cache

import (
	"errors"
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/logging"
	lru "github.com/hashicorp/golang-lru/v2"
)

var (
	// GetInstance retrieves instance of the networking cache, performing initilization if needed.
	GetInstance = sync.OnceValue[*NetworkingCache](newNetworkingCache)
)

var ErrFailedToEvictSessionCacheEntity = errors.New("err failed to evict session cache entity")

const (
	// Represents max amount of users per session.
	maxSessionUsers = 8
)

// NetworkingCache represents networking cache.
type NetworkingCache struct {
	// Represents sessions cache instance.
	sessions *lru.Cache[string, []string]

	// Represents expirable messages cache, which contains offset for the message table.
	// If user stops request messages, all the messages would be retrieved.
	messages *lru.Cache[string, int]

	// Represents users cache instance.
	users *lru.Cache[string, int64]
}

// AddSession adds session cache instance with the provided key and value.
func (nc *NetworkingCache) AddSessions(key string, value []string) {
	nc.sessions.Add(key, value)
}

// GetSession retrieves session cache instance by the provided key.
func (nc *NetworkingCache) GetSessions(key string) (any, bool) {
	return nc.sessions.Get(key)
}

// EvictSessions evicts sessions cache for the provided key.
func (nc *NetworkingCache) EvictSessions(key string) {
	if ok := nc.sessions.Remove(key); !ok {
		logging.GetInstance().Error(ErrFailedToEvictSessionCacheEntity.Error())
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
	sessions, err := lru.New[string, []string](config.GetOperationMaxSessionsAmount())
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
		sessions: sessions,
		messages: messages,
		users:    users,
	}
}
