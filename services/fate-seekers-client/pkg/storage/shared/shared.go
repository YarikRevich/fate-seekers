package shared

import (
	"sync"
	"time"

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

	// Represents global blinking screen animation.
	blinkingScreenAnimation *combiner.AnimationCombiner

	// Represents global shader start time.
	shaderStartTime time.Time
}

// GetBackgroundAnimation retrieves background animation.
func (ss *SharedStorage) GetBackgroundAnimation() *combiner.AnimationCombiner {
	return ss.backgroundAnimation
}

// GetBlinkingScreenAnimation retrieves blinking screen animation.
func (ss *SharedStorage) GetBlinkingScreenAnimation() *combiner.AnimationCombiner {
	return ss.blinkingScreenAnimation
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
		blinkingScreenAnimation: combiner.NewAnimationCombiner(
			loader.GetInstance().GetAnimation(loader.BlinkingScreen1Animation, false),
			loader.GetInstance().GetAnimation(loader.BlinkingScreen2Animation, false),
			loader.GetInstance().GetAnimation(loader.BlinkingScreen3Animation, false),
			loader.GetInstance().GetAnimation(loader.BlinkingScreen4Animation, false),
		),
		shaderStartTime: time.Now(),
	}
}
