package sound

import (
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/sound/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/sound/fx"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/sound/music"
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

	// Represents instance of sound music manager.
	soundMusicManager *music.SoundMusicManager
}

// GetSoundFxManager retrieves instance of sound FX manager.
func (sm *SoundManager) GetSoundFxManager() *fx.SoundFXManager {
	return sm.soundFxManager
}

// GetSoundMusicManager retrieves instance of sound music manager.
func (sm *SoundManager) GetSoundMusicManager() *music.SoundMusicManager {
	return sm.soundMusicManager
}

// initSoundAmbientBatch performs ambient sound players batch initialization.
func (sm *SoundManager) initSoundAmbientBatch() {
	// sm.soundMusicManager.PushAmbient(loader.TestVFXSound)
	// sm.soundMusicManager.PushAmbient(loader.TestVFXSound)
	// sm.soundMusicManager.PushAmbient(loader.TestVFXSound)
}

// newSoundManager initializes SoundManager.
func newSoundManager() *SoundManager {
	audioContext := audio.NewContext(common.SampleRate)

	soundFxManager := fx.NewSoundFxManager(audioContext)
	soundFxManager.Init()

	soundMusicManager := music.NewSoundMusicManager(audioContext)
	soundMusicManager.Init()

	result := &SoundManager{
		soundFxManager:    soundFxManager,
		soundMusicManager: soundMusicManager,
	}

	result.initSoundAmbientBatch()

	return result
}
