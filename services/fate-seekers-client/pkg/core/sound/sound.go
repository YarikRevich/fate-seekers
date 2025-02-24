package sound

import "sync"

var (
	// GetCollectionsRepository retrieves instance of the collections repository, performing initial creation if needed.
	GetInstance = sync.OnceValue[CollectionsRepository](createCollectionsRepository)
)

type SoundManager struct {
}

func newSoundManager() *SoundManager {
	return nil
}
