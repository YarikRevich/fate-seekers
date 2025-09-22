package sound

import (
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/sound/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/sound/fx"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

var (
	// GetInstance retrieves instance of the sound manager, performing initial creation if needed.
	GetInstance = sync.OnceValue[*SoundManager](newSoundManager)
)

// SoundManager represents global sound manager.
type SoundManager struct {
	// Represents instance of sound FX manager.
	soundFxManager *fx.SoundFXManager
}

// GetSoundFxManager retrieves instance of sound FX manager.
func (sm *SoundManager) GetSoundFxManager() *fx.SoundFXManager {
	return sm.soundFxManager
}

// newSoundManager initializes SoundManager.
func newSoundManager() *SoundManager {
	audioContext := audio.NewContext(common.SampleRate)

	soundFxManager := fx.NewSoundFxManager(audioContext)
	soundFxManager.Init()

	result := &SoundManager{
		soundFxManager: soundFxManager,
	}

	return result
}
