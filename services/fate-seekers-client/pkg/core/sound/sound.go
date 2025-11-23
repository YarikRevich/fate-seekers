package sound

import (
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/sound/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/sound/fx"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/sound/music"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

var (
	// GetInstance retrieves instance of the sound manager, performing initial creation if needed.
	GetInstance = sync.OnceValue[*SoundManager](newSoundManager)
)

// SoundManager represents global sound manager.
type SoundManager struct {
	// Represents instance of sound UI FX manager.
	soundUIFxManager *fx.SoundFXManager

	// Represents instance of sound events FX manager.
	soundEventsFxManager *fx.SoundFXManager

	// Represents instance of sound sounder main FX manager.
	soundSounderMainFxManager *fx.SoundFXManager

	// Represents instance of sound music manager.
	soundMusicManager *music.SoundMusicManager
}

// GetSoundUIFxManager retrieves instance of sound UI FX manager.
func (sm *SoundManager) GetSoundUIFxManager() *fx.SoundFXManager {
	return sm.soundUIFxManager
}

// GetSoundEventsFxManager retrieves instance of sound events FX manager.
func (sm *SoundManager) GetSoundEventsFxManager() *fx.SoundFXManager {
	return sm.soundEventsFxManager
}

// GetSoundSounderMainFxManager retrieves instance of sound sounder main FX manager.
func (sm *SoundManager) GetSoundSounderMainFxManager() *fx.SoundFXManager {
	return sm.soundEventsFxManager
}

// GetSoundMusicManager retrieves instance of sound music manager.
func (sm *SoundManager) GetSoundMusicManager() *music.SoundMusicManager {
	return sm.soundMusicManager
}

// InitSoundAmbientBatch performs ambient sound players batch initialization.
func (sm *SoundManager) InitSoundAmbientBatch() {
	sm.soundMusicManager.PushAmbient(loader.AmbientMusicSound)
}

// newSoundManager initializes SoundManager.
func newSoundManager() *SoundManager {
	audioContext := audio.NewContext(common.SampleRate)

	soundUIFxManager := fx.NewSoundFxManager(audioContext)
	soundUIFxManager.Init()

	soundEventsFxManager := fx.NewSoundFxManager(audioContext)
	soundEventsFxManager.Init()

	soundSounderMainFxManager := fx.NewSoundFxManager(audioContext)
	soundSounderMainFxManager.Init()

	soundMusicManager := music.NewSoundMusicManager(audioContext)
	soundMusicManager.Init()

	result := &SoundManager{
		soundUIFxManager:          soundUIFxManager,
		soundEventsFxManager:      soundEventsFxManager,
		soundSounderMainFxManager: soundSounderMainFxManager,
		soundMusicManager:         soundMusicManager,
	}

	return result
}
