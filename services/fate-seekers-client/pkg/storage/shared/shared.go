package shared

import (
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/animation/combiner"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
)

var (
	// GetInstance retrieves instance of the shared storage, performing initilization if needed.
	GetInstance = sync.OnceValue[*SharedStorage](newSharedStorage)
)

// SharedStorage represents shared storage holder.
type SharedStorage struct {
	// Represents global background animation.
	backgroundAnimation *combiner.AnimationCombiner
}

// GetBackgroundAnimation retrieves background animation.
func (ss *SharedStorage) GetBackgroundAnimation() *combiner.AnimationCombiner {
	return ss.backgroundAnimation
}

// newSharedStorage initializes shared storage.
func newSharedStorage() *SharedStorage {
	return &SharedStorage{
		backgroundAnimation: combiner.NewAnimationCombiner(
			loader.GetInstance().GetAnimation(loader.BackgroundAnimation, false),
		),
	}
}
