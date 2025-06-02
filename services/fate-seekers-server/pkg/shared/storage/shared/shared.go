package shared

import (
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/loader"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/tools/animation/combiner"
)

var (
	// GetInstance retrieves instance of the shared storage, performing initilization if needed.
	GetInstance = sync.OnceValue[*SharedStorage](newSharedStorage)
)

// SharedStorage represents shared storage holder.
type SharedStorage struct {
	// Represents global background animation.
	backgroundAnimation *combiner.AnimationCombiner

	// Represents global shader start time.
	shaderStartTime time.Time
}

// GetBackgroundAnimation retrieves background animation.
func (ss *SharedStorage) GetBackgroundAnimation() *combiner.AnimationCombiner {
	return ss.backgroundAnimation
}

// GetShaderTime represents shader global time.
func (ss *SharedStorage) GetShaderTime() float32 {
	return float32(time.Since(ss.shaderStartTime).Seconds())
}

// newSharedStorage initializes shared storage.
func newSharedStorage() *SharedStorage {
	return &SharedStorage{
		backgroundAnimation: combiner.NewAnimationCombiner(
			loader.GetInstance().GetAnimation(loader.Background1Animation, false),
			loader.GetInstance().GetAnimation(loader.Background2Animation, false),
			loader.GetInstance().GetAnimation(loader.Background3Animation, false),
			loader.GetInstance().GetAnimation(loader.Background4Animation, false),
			loader.GetInstance().GetAnimation(loader.Background5Animation, false),
			loader.GetInstance().GetAnimation(loader.Background6Animation, false),
		),
		shaderStartTime: time.Now(),
	}
}
