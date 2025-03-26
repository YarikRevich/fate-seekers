package music

import (
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/sound/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/logging"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/pkg/errors"
)

var (
	ErrMusicIsNotPlaying      = errors.New("err happened music is not playing")
	ErrMusicIsAlreadyStopping = errors.New("err happened music is already stopping")
	ErrMusicIsAlreadyPlaying  = errors.New("err happened music is already playing")
)

const (
	// Represents processing ticker duration.
	processingTickerPeriod = time.Millisecond * 400

	// Represents start temp ticker period.
	startTempTickerPeriod = time.Millisecond * 100

	// Represents end temp ticker period.
	endTempTickerPeriod = time.Millisecond * 120

	// Represents volume ticker period.
	volumeTickerPeriod = time.Second

	// Represents fadeout duration.
	fadeoutDuration = time.Second * 3

	// Represents volume shift step.
	volumeShift = 20.0
)

// SoundMusicManager represents sound music manager, used for both ambient and music streams management.
type SoundMusicManager struct {
	// Represents audio context used for stream creation.
	audioContext *audio.Context

	// Represents time ticker used for player processing batch.
	ticker *time.Ticker

	// Represents processing ambient player infinite batch.
	ambientProcessingBatch []*dto.AmbientSoundUnit

	// Represents lock used to interrupt ambient processing operation.
	ambientProcessingInterruption atomic.Bool

	// Represents currently playing ambient player.
	currentAmbientPlayer atomic.Pointer[dto.AmbientSoundUnit]

	// Represents state if ambient volume is configured.
	ambientVolumeConfigured atomic.Bool

	// Represents currently playing music player.
	currentMusicPlayer atomic.Pointer[dto.MusicSoundUnit]

	// Represents state if music volume is configured.
	musicVolumeConfigured atomic.Bool

	// Rerpesents state if music end state is passively stopped processed.
	musicEndProcessed sync.Mutex

	// Represents music stopping state.
	musicStopping atomic.Bool
}

// Init starts sound music manager processing worker. Additionally starts volume update worker.
func (smm *SoundMusicManager) Init() {
	go func() {
		for {
			select {
			case <-smm.ticker.C:
				if len(smm.ambientProcessingBatch) > 0 && smm.currentAmbientPlayer.Load() == nil {
					smm.ambientVolumeConfigured.Store(false)

					smm.currentAmbientPlayer.Store(smm.ambientProcessingBatch[0])

					smm.ambientProcessingBatch = append(smm.ambientProcessingBatch[:0], smm.ambientProcessingBatch[1:]...)

					sound := loader.GetInstance().GetSoundMusic(smm.currentAmbientPlayer.Load().Name)

					player, err := smm.audioContext.NewPlayerF32(sound)
					if err != nil {
						logging.GetInstance().Fatal(errors.Wrap(err, common.ErrSoundPlayerAccess.Error()).Error())
					}

					smm.ambientProcessingBatch = append(smm.ambientProcessingBatch, &dto.AmbientSoundUnit{
						Name:     smm.currentAmbientPlayer.Load().Name,
						Duration: time.Second * time.Duration(sound.Length()) / common.BytesPerSample / common.SampleRate,
						Player:   player,
					})

					smm.currentAmbientPlayer.Load().Player.SetVolume(0)

					interruptionChan := make(chan int, 1)

					go func() {
						tempTicker := time.NewTicker(startTempTickerPeriod)

						var volume float64

						for smm.currentAmbientPlayer.Load() != nil && smm.currentAmbientPlayer.Load().Player.Volume() < float64(config.GetSettingsSoundMusic())/100 {
							if smm.ambientProcessingInterruption.Load() {
								tempTicker.Stop()

								smm.ambientVolumeConfigured.Store(true)

								return
							}

							select {
							case <-interruptionChan:
								tempTicker.Stop()

								return
							case <-tempTicker.C:
								tempTicker.Stop()

								if volume+volumeShift <= float64(config.GetSettingsSoundMusic()) {
									volume += volumeShift
								} else {
									volume = float64(config.GetSettingsSoundMusic())
								}

								smm.currentAmbientPlayer.Load().Player.SetVolume(volume / 100)

								tempTicker.Reset(startTempTickerPeriod)
							}
						}

						smm.ambientVolumeConfigured.Store(true)
					}()

					smm.currentAmbientPlayer.Load().Player.Play()

					go func() {
						tempTicker := time.NewTicker(endTempTickerPeriod)

						for {
							select {
							case <-tempTicker.C:
								tempTicker.Stop()

								if smm.currentAmbientPlayer.Load().Duration-smm.currentAmbientPlayer.Load().Player.Position() <= fadeoutDuration &&
									!smm.ambientProcessingInterruption.Load() {
									if smm.ambientVolumeConfigured.Load() {
										smm.ambientVolumeConfigured.Store(false)
									}

									select {
									case interruptionChan <- 1:
									default:
									}

									if smm.currentAmbientPlayer.Load().Player.Volume() != 0 {
										normalized := smm.currentAmbientPlayer.Load().Player.Volume() / 2

										fadeFactor := math.Max(0, math.Sin((math.Pi/2)*normalized))

										newVolume := float64(smm.currentAmbientPlayer.Load().Player.Volume()) * fadeFactor

										if smm.currentAmbientPlayer.Load().Player.Volume() >= 0 {
											smm.currentAmbientPlayer.Load().Player.SetVolume(newVolume * float64(config.GetSettingsSoundMusic()) / 100)
										}
									}
								}

								if !smm.currentAmbientPlayer.Load().Player.IsPlaying() {
									smm.ambientVolumeConfigured.Store(true)

									if err := smm.currentAmbientPlayer.Load().Player.Close(); err != nil {
										logging.GetInstance().Fatal(errors.Wrap(err, common.ErrSoundPlayerAccess.Error()).Error())
									}

									tempTicker.Stop()

									smm.currentAmbientPlayer.Store(nil)

									return
								}

								tempTicker.Reset(endTempTickerPeriod)
							}
						}
					}()
				}
			}
		}
	}()

	go func() {
		volumeTicker := time.NewTicker(volumeTickerPeriod)

		for {
			select {
			case <-volumeTicker.C:
				if store.GetSoundMusicUpdated() == value.SOUND_MUSIC_UPDATED_FALSE_VALUE {
					if smm.ambientVolumeConfigured.Load() && !smm.ambientProcessingInterruption.Load() {
						if smm.IsMusicPlaying() && !smm.musicVolumeConfigured.Load() {
							continue
						}

						dispatcher.GetInstance().Dispatch(
							action.NewSetSoundMusicUpdated(value.SOUND_MUSIC_UPDATED_TRUE_VALUE))

						smm.currentAmbientPlayer.Load().Player.SetVolume(float64(config.GetSettingsSoundMusic()) / 100)

						if smm.IsMusicPlaying() {
							smm.currentMusicPlayer.Load().Player.SetVolume(float64(config.GetSettingsSoundMusic()) / 100)
						}
					}
				}
			}
		}
	}()
}

// PushAmbient pushes a new ambient track to the queue.
func (smm *SoundMusicManager) PushAmbient(name string) {
	sound := loader.GetInstance().GetSoundMusic(name)

	player, err := smm.audioContext.NewPlayerF32(sound)
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, common.ErrSoundPlayerAccess.Error()).Error())
	}

	smm.ambientProcessingBatch = append(smm.ambientProcessingBatch, &dto.AmbientSoundUnit{
		Name:     name,
		Duration: time.Second * time.Duration(sound.Length()) / common.BytesPerSample / common.SampleRate,
		Player:   player,
	})
}

// IsMusicPlaying checks if music player is currently active and playing.
func (smm *SoundMusicManager) IsMusicPlaying() bool {
	return smm.currentMusicPlayer.Load() != nil
}

// IsMusicStopping checks if music player is currently stopping.
func (smm *SoundMusicManager) IsMusicStopping() bool {
	return smm.musicStopping.Load()
}

// StartMusic starts playing music stream of the provided name.
func (smm *SoundMusicManager) StartMusic(name string) {
	if smm.currentMusicPlayer.Load() != nil {
		logging.GetInstance().Fatal(ErrMusicIsAlreadyPlaying.Error())
	}

	go func() {
		smm.musicEndProcessed.Lock()

		sound := loader.GetInstance().GetSoundMusic(name)

		player, err := smm.audioContext.NewPlayerF32(sound)
		if err != nil {
			logging.GetInstance().Fatal(errors.Wrap(err, common.ErrSoundPlayerAccess.Error()).Error())
		}

		smm.currentMusicPlayer.Store(&dto.MusicSoundUnit{
			Duration: time.Second * time.Duration(sound.Length()) / common.BytesPerSample / common.SampleRate,
			Player:   player,
		})

		tempTicker := time.NewTicker(startTempTickerPeriod)

		if smm.currentAmbientPlayer.Load() != nil {
			smm.ambientProcessingInterruption.Store(true)

			for smm.currentAmbientPlayer.Load().Player.Volume() != 0 {
				select {
				case <-tempTicker.C:
					normalized := smm.currentAmbientPlayer.Load().Player.Volume() / 1.1

					fadeFactor := math.Max(0, math.Sin((math.Pi/2)*normalized))

					newVolume := float64(smm.currentAmbientPlayer.Load().Player.Volume()) * fadeFactor

					if smm.currentAmbientPlayer.Load().Player.Volume() >= 0 {
						smm.currentAmbientPlayer.Load().Player.SetVolume(newVolume * float64(config.GetSettingsSoundMusic()) / 100)
					}
				}
			}

			smm.ambientProcessingInterruption.Store(false)
		}

		smm.currentMusicPlayer.Load().Player.SetVolume(0)

		smm.currentMusicPlayer.Load().Player.Play()

		smm.musicVolumeConfigured.Store(false)

		interruptionChan := make(chan int, 1)

		var volume float64

		for smm.currentMusicPlayer.Load() != nil && smm.currentMusicPlayer.Load().Player.Volume() < float64(config.GetSettingsSoundMusic())/100 {
			select {
			case <-interruptionChan:
				tempTicker.Stop()

				smm.musicVolumeConfigured.Store(true)

				return
			case <-tempTicker.C:
				tempTicker.Stop()

				if volume+volumeShift <= float64(config.GetSettingsSoundMusic()) {
					volume += volumeShift
				} else {
					volume = float64(config.GetSettingsSoundMusic())
				}

				smm.currentMusicPlayer.Load().Player.SetVolume(volume / 100)

				tempTicker.Reset(startTempTickerPeriod)
			}
		}

		smm.musicVolumeConfigured.Store(true)

		tempTicker.Stop()

		smm.musicEndProcessed.Unlock()

		go func() {
			tempTicker := time.NewTicker(endTempTickerPeriod)

			for {
				select {
				case <-tempTicker.C:
					tempTicker.Stop()

					smm.musicEndProcessed.Lock()

					if smm.currentMusicPlayer.Load() != nil {
						if smm.currentMusicPlayer.Load().Duration-smm.currentMusicPlayer.Load().Player.Position() <= fadeoutDuration {
							if !smm.ambientProcessingInterruption.Load() {
								smm.ambientProcessingInterruption.Store(true)
							}

							select {
							case interruptionChan <- 1:
							default:
							}

							if smm.currentMusicPlayer.Load().Player.Volume() != 0 {
								normalized := smm.currentMusicPlayer.Load().Player.Volume() / 2

								fadeFactor := math.Max(0, math.Sin((math.Pi/2)*normalized))

								newVolume := float64(smm.currentMusicPlayer.Load().Player.Volume()) * fadeFactor

								if smm.currentMusicPlayer.Load().Player.Volume() >= 0 {
									smm.currentMusicPlayer.Load().Player.SetVolume(newVolume * float64(config.GetSettingsSoundMusic()) / 100)
								}
							}
						}

						if !smm.currentMusicPlayer.Load().Player.IsPlaying() {
							if err := smm.currentMusicPlayer.Load().Player.Close(); err != nil {
								logging.GetInstance().Fatal(errors.Wrap(err, common.ErrSoundPlayerAccess.Error()).Error())
							}

							tempTicker.Stop()

							smm.ambientProcessingInterruption.Store(false)

							smm.musicEndProcessed.Unlock()

							smm.currentMusicPlayer.Store(nil)

							return
						}
					}

					smm.musicEndProcessed.Unlock()

					tempTicker.Reset(endTempTickerPeriod)
				}
			}
		}()
	}()
}

// StopMusic stops currently playing music stream.
func (smm *SoundMusicManager) StopMusic() {
	if smm.currentMusicPlayer.Load() == nil {
		logging.GetInstance().Fatal(ErrMusicIsNotPlaying.Error())
	}

	if smm.musicStopping.Load() {
		logging.GetInstance().Fatal(ErrMusicIsAlreadyStopping.Error())
	}

	smm.musicStopping.Store(true)

	go func() {
		smm.musicEndProcessed.Lock()

		tempTicker := time.NewTicker(endTempTickerPeriod)

		smm.musicVolumeConfigured.Store(false)

		for smm.currentMusicPlayer.Load().Player.Volume() != 0 {
			select {
			case <-tempTicker.C:
				normalized := smm.currentMusicPlayer.Load().Player.Volume() / 2

				fadeFactor := math.Max(0, math.Sin((math.Pi/2)*normalized))

				newVolume := float64(smm.currentMusicPlayer.Load().Player.Volume()) * fadeFactor

				if smm.currentMusicPlayer.Load().Player.Volume() >= 0 {
					smm.currentMusicPlayer.Load().Player.SetVolume(newVolume * float64(config.GetSettingsSoundMusic()) / 100)
				}
			}
		}

		smm.musicVolumeConfigured.Store(true)

		smm.currentMusicPlayer.Load().Player.Pause()

		if smm.currentAmbientPlayer.Load() != nil {
			volume := smm.currentAmbientPlayer.Load().Player.Volume() * float64(config.GetSettingsSoundMusic())

			smm.ambientProcessingInterruption.Store(true)

			for smm.currentAmbientPlayer.Load().Player.Volume() < float64(config.GetSettingsSoundMusic())/100 {
				select {
				case <-tempTicker.C:
					if volume+volumeShift <= float64(config.GetSettingsSoundMusic()) {
						volume += volumeShift
					} else {
						volume = float64(config.GetSettingsSoundMusic())
					}

					smm.currentAmbientPlayer.Load().Player.SetVolume(volume / 100)
				}
			}

			smm.ambientProcessingInterruption.Store(false)
		}

		if err := smm.currentMusicPlayer.Load().Player.Close(); err != nil {
			logging.GetInstance().Fatal(errors.Wrap(err, common.ErrSoundPlayerAccess.Error()).Error())
		}

		tempTicker.Stop()

		smm.currentMusicPlayer.Store(nil)

		smm.musicEndProcessed.Unlock()

		smm.musicStopping.Store(false)
	}()
}

// NewSoundMusicManager initializes SoundMusicManager.
func NewSoundMusicManager(audioContext *audio.Context) *SoundMusicManager {
	return &SoundMusicManager{
		audioContext: audioContext,
		ticker:       time.NewTicker(processingTickerPeriod),
	}
}
