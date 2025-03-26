package fx

import (
	"math"
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

const (
	// Represents processing ticker duration.
	processingTickerPeriod = time.Millisecond * 400

	// Represents start temp ticker period.
	startTempTickerPeriod = time.Millisecond * 10

	// Represents end temp ticker period.
	endTempTickerPeriod = time.Millisecond * 120

	// Represents volume ticker period.
	volumeTickerPeriod = time.Second

	// Represents fadeout duration.
	fadeoutDuration = time.Millisecond * 200

	// Represents volume shift step.
	volumeShift = 30.0
)

// SoundFXManager represents sound FX manager.
type SoundFXManager struct {
	// Represents audio context used for stream creation.
	audioContext *audio.Context

	// Represents time ticker used for player processing batch.
	ticker *time.Ticker

	// Represents processing player batch.
	processingBatch []*dto.FXSoundUnit

	// Represents currently playing player.
	currentPlayer atomic.Pointer[dto.FXSoundUnit]

	// Represents state if volume is configured.
	volumeConfigured atomic.Bool
}

// Init starts sound FX manager processing worker. Additionally starts volume update worker.
func (svm *SoundFXManager) Init() {
	go func() {
		for {
			select {
			case <-svm.ticker.C:
				if len(svm.processingBatch) > 0 && svm.currentPlayer.Load() == nil {
					svm.volumeConfigured.Store(false)

					svm.currentPlayer.Store(svm.processingBatch[0])

					svm.processingBatch = append(svm.processingBatch[:0], svm.processingBatch[1:]...)

					svm.currentPlayer.Load().Player.SetVolume(0)

					interruptionChan := make(chan int, 1)

					go func() {
						tempTicker := time.NewTicker(startTempTickerPeriod)

						var volume float64

						for svm.currentPlayer.Load() != nil && svm.currentPlayer.Load().Player.Volume() < float64(config.GetSettingsSoundFX())/100 {
							select {
							case <-interruptionChan:
								tempTicker.Stop()

								svm.volumeConfigured.Store(true)

								return
							case <-tempTicker.C:
								tempTicker.Stop()

								if volume+volumeShift <= float64(config.GetSettingsSoundFX()) {
									volume += volumeShift
								} else {
									volume = float64(config.GetSettingsSoundFX())
								}

								svm.currentPlayer.Load().Player.SetVolume(volume / 100)

								tempTicker.Reset(startTempTickerPeriod)
							}
						}

						svm.volumeConfigured.Store(true)
					}()

					svm.currentPlayer.Load().Player.Play()

					go func() {
						tempTicker := time.NewTicker(endTempTickerPeriod)

						for {
							select {
							case <-tempTicker.C:
								tempTicker.Stop()

								if svm.currentPlayer.Load().Duration-svm.currentPlayer.Load().Player.Position() <= fadeoutDuration {
									if svm.volumeConfigured.Load() {
										svm.volumeConfigured.Store(false)
									}

									select {
									case interruptionChan <- 1:
									default:
									}

									if svm.currentPlayer.Load().Player.Volume() != 0 {
										normalized := svm.currentPlayer.Load().Player.Volume() / 2

										fadeFactor := math.Max(0, math.Sin((math.Pi/2)*normalized))

										newVolume := float64(svm.currentPlayer.Load().Player.Volume()) * fadeFactor

										if svm.currentPlayer.Load().Player.Volume() >= 0 {
											svm.currentPlayer.Load().Player.SetVolume(newVolume * float64(config.GetSettingsSoundFX()) / 100)
										}
									}
								}

								if !svm.currentPlayer.Load().Player.IsPlaying() {
									svm.volumeConfigured.Store(true)

									if err := svm.currentPlayer.Load().Player.Close(); err != nil {
										logging.GetInstance().Fatal(errors.Wrap(err, common.ErrSoundPlayerAccess.Error()).Error())
									}

									tempTicker.Stop()

									svm.currentPlayer.Store(nil)

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
				if store.GetSoundFXUpdated() == value.SOUND_FX_UPDATED_FALSE_VALUE {
					if svm.volumeConfigured.Load() {
						dispatcher.GetInstance().Dispatch(
							action.NewSetSoundFXUpdated(value.SOUND_FX_UPDATED_TRUE_VALUE))

						if svm.currentPlayer.Load() != nil {
							svm.currentPlayer.Load().Player.SetVolume(float64(config.GetSettingsSoundFX()) / 100)
						}
					}
				}
			}
		}
	}()
}

// PushImmediately pushes a new track immediately to the queue at the highest priority.
func (svm *SoundFXManager) PushImmediately(name string) {
	sound := loader.GetInstance().GetSoundFX(name)

	player, err := svm.audioContext.NewPlayerF32(sound)
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, common.ErrSoundPlayerAccess.Error()).Error())
	}

	svm.processingBatch = append([]*dto.FXSoundUnit{{
		Duration: time.Second * time.Duration(sound.Length()) / common.BytesPerSample / common.SampleRate,
		Player:   player,
	}}, svm.processingBatch...)
}

// Push pushes a new track to the queue at the end.
func (svm *SoundFXManager) Push(name string) {
	sound := loader.GetInstance().GetSoundFX(name)

	player, err := svm.audioContext.NewPlayerF32(sound)
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, common.ErrSoundPlayerAccess.Error()).Error())
	}

	svm.processingBatch = append(svm.processingBatch, &dto.FXSoundUnit{
		Duration: time.Second * time.Duration(sound.Length()) / common.BytesPerSample / common.SampleRate,
		Player:   player,
	})
}

// NewSoundFxManager initializes SoundFXManager.
func NewSoundFxManager(audioContext *audio.Context) *SoundFXManager {
	return &SoundFXManager{
		audioContext: audioContext,
		ticker:       time.NewTicker(processingTickerPeriod),
	}
}
