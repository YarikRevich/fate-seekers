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
	// Represents audio context.
	audioContext *audio.Context

	// Represents instance of sound UI FX manager.
	soundUIFxManager *fx.SoundFXManager

	// Represents instance of sound events FX manager.
	soundEventsFxManager *fx.SoundFXManager

	// Represents instance of sound events announcement FX manager.
	soundEventsAnnouncementFxManager *fx.SoundFXManager

	// Represents instance of sound sounder steps FX manager.
	soundSounderStepsFxManager *fx.SoundFXManager

	// Represents instance of sound sounder melee FX manager.
	soundSounderMeleeFxManager *fx.SoundFXManager

	// Represents instance of sound sounder chest FX manager.
	soundSounderChestFxManager *fx.SoundFXManager

	// Represents instance of sound sounder letter scroll activation FX manager.
	soundSounderLetterScrollActivationFxManager *fx.SoundFXManager

	// Represents instance of sound sounder health pack activation FX manager.
	soundSounderHealthPackActivationFxManager *fx.SoundFXManager

	// Represents map of sound sounder external FX manager instances.
	soundSounderExternalFxManagers map[string]*fx.SoundFXManager

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

// GetSoundEventsAnnouncementFxManager retrieves instance of sound events announcement FX manager.
func (sm *SoundManager) GetSoundEventsAnnouncementFxManager() *fx.SoundFXManager {
	return sm.soundEventsAnnouncementFxManager
}

// GetSoundSounderStepsFxManager retrieves instance of sound sounder steps FX manager.
func (sm *SoundManager) GetSoundSounderStepsFxManager() *fx.SoundFXManager {
	return sm.soundSounderStepsFxManager
}

// GetSoundSounderMeleeFxManager retrieves instance of sound sounder melee FX manager.
func (sm *SoundManager) GetSoundSounderMeleeFxManager() *fx.SoundFXManager {
	return sm.soundSounderMeleeFxManager
}

// GetSoundSounderChestFxManager retrieves instance of sound sounder chest FX manager.
func (sm *SoundManager) GetSoundSounderChestFxManager() *fx.SoundFXManager {
	return sm.soundSounderChestFxManager
}

// GetSoundSounderLetterScrollActivationFxManager retrieves instance of sound sounder letter scroll activation FX manager.
func (sm *SoundManager) GetSoundSounderLetterScrollActivationFxManager() *fx.SoundFXManager {
	return sm.soundSounderLetterScrollActivationFxManager
}

// GetSoundSounderHealthPackActivationFxManager retrieves instance of sound sounder health pack activation FX manager.
func (sm *SoundManager) GetSoundSounderHealthPackActivationFxManager() *fx.SoundFXManager {
	return sm.soundSounderLetterScrollActivationFxManager
}

// AddSoundSounderExternalFxManager adds an instance of sound sounder external FX manager.
func (sm *SoundManager) AddSoundSounderExternalFxManager(issuer string) {
	soundSounderExternalFxManager := fx.NewSoundFxManager(sm.audioContext)
	soundSounderExternalFxManager.Init()

	sm.soundSounderExternalFxManagers[issuer] = soundSounderExternalFxManager
}

// SoundSounderExternalFxManagerExists checks if instance of sound sounder external FX manager
// exists for the provided issuer.
func (sm *SoundManager) SoundSounderExternalFxManagerExists(issuer string) bool {
	_, ok := sm.soundSounderExternalFxManagers[issuer]

	return ok
}

// GetSoundSounderExternalFxManager retrieves instance of sound sounder external FX manager.
func (sm *SoundManager) GetSoundSounderExternalFxManager(issuer string) *fx.SoundFXManager {
	return sm.soundSounderExternalFxManagers[issuer]
}

// RemoveSoundSounderExternalFxManager performs a removal of the instance of sound sounder external FX manager
// for the provided issuer.
func (sm *SoundManager) RemoveSoundSounderExternalFxManager(issuer string) {
	delete(sm.soundSounderExternalFxManagers, issuer)
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

	soundEventsAnnouncementFxManager := fx.NewSoundFxManager(audioContext)
	soundEventsAnnouncementFxManager.Init()

	soundSounderStepsFxManager := fx.NewSoundFxManager(audioContext)
	soundSounderStepsFxManager.Init()

	soundSounderMeleeFxManager := fx.NewSoundFxManager(audioContext)
	soundSounderMeleeFxManager.Init()

	soundSounderChestFxManager := fx.NewSoundFxManager(audioContext)
	soundSounderChestFxManager.Init()

	soundSounderLetterScrollActivationFxManager := fx.NewSoundFxManager(audioContext)
	soundSounderLetterScrollActivationFxManager.Init()

	soundSounderHealthPackActivationFxManager := fx.NewSoundFxManager(audioContext)
	soundSounderHealthPackActivationFxManager.Init()

	soundMusicManager := music.NewSoundMusicManager(audioContext)
	soundMusicManager.Init()

	result := &SoundManager{
		audioContext:                                audioContext,
		soundUIFxManager:                            soundUIFxManager,
		soundEventsFxManager:                        soundEventsFxManager,
		soundEventsAnnouncementFxManager:            soundEventsAnnouncementFxManager,
		soundSounderStepsFxManager:                  soundSounderStepsFxManager,
		soundSounderMeleeFxManager:                  soundSounderMeleeFxManager,
		soundSounderChestFxManager:                  soundSounderChestFxManager,
		soundSounderLetterScrollActivationFxManager: soundSounderLetterScrollActivationFxManager,
		soundSounderHealthPackActivationFxManager:   soundSounderHealthPackActivationFxManager,
		soundMusicManager:                           soundMusicManager,
		soundSounderExternalFxManagers:              make(map[string]*fx.SoundFXManager),
	}

	return result
}
