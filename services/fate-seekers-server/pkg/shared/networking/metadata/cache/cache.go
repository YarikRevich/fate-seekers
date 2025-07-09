package cache

import (
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/logging"
	lru "github.com/hashicorp/golang-lru/v2"
)

var (
	// GetInstance retrieves instance of the networking metadata cache, performing initilization if needed.
	GetInstance = sync.OnceValue[*NetworkingMetadataCache](newNetworkingMetadataCache)
)

const (
	// Represents sessions cache instance size.
	sessionsSize = 128
)

// NetworkingMetadataCache represents networking metadata cache.
type NetworkingMetadataCache struct {
	// Represents sessions cache instance.
	sessions *lru.Cache[string, any]
}

func (nmc *NetworkingMetadataCache) GetSessions() *lru.Cache[string, any] {
	return nmc.sessions
}

// newNetworkingMetadataCache initializes NetworkingMetadataCache.
func newNetworkingMetadataCache() *NetworkingMetadataCache {
	sessions, err := lru.New[string, any](sessionsSize)
	if err != nil {
		logging.GetInstance().Fatal(err.Error())
	}

	return &NetworkingMetadataCache{
		sessions: sessions,
	}
}
