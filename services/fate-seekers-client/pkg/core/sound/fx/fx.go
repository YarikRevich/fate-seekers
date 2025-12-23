package fx

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
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/pkg/errors"
)

const (
	// Represents processing ticker duration.
	processingTickerPeriod = time.Millisecond * 20

	// Represents start temp ticker period.
	startTempTickerPeriod = time.Millisecond * 20

	// Represents end temp ticker period.
	endTempTickerPeriod = time.Millisecond * 20

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
	// volumeConfigured atomic.Bool
	playerStarter sync.WaitGroup

	// Represents handbrake state channel.
	handbrake chan bool
}

// Init starts sound FX manager processing worker. Additionally starts volume update worker.
func (svm *SoundFXManager) Init() {
	go func() {
		for {
			select {
			case <-svm.ticker.C:
				if len(svm.processingBatch) > 0 && svm.currentPlayer.Load() == nil {
					svm.currentPlayer.Store(svm.processingBatch[0])

					svm.processingBatch = append(svm.processingBatch[:0], svm.processingBatch[1:]...)

					svm.currentPlayer.Load().Player.SetVolume(0)

					svm.playerStarter.Add(1)

					go func() {
						tempTicker := time.NewTicker(startTempTickerPeriod)

						var volume float64

						for svm.currentPlayer.Load() != nil && svm.currentPlayer.Load().Player.Volume() < float64(config.GetSettingsSoundFX())/100 {
							select {
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

						svm.playerStarter.Done()
					}()

					svm.currentPlayer.Load().Player.Play()

					go func() {
						svm.playerStarter.Wait()

						tempTicker := time.NewTicker(endTempTickerPeriod)

						for {
							select {
							case <-tempTicker.C:
								tempTicker.Stop()

								if svm.currentPlayer.Load().Duration-svm.currentPlayer.Load().Player.Position() <= fadeoutDuration {
									if svm.currentPlayer.Load().Player.Volume() != 0 {
										normalized := svm.currentPlayer.Load().Player.Volume() / 2

										fadeFactor := math.Max(0, math.Sin((math.Pi/2)*normalized))

										newVolume := float64(svm.currentPlayer.Load().Player.Volume()) * fadeFactor

										if svm.currentPlayer.Load().Player.Volume() >= 0 {
											svm.currentPlayer.Load().Player.SetVolume(newVolume * float64(config.GetSettingsSoundFX()) / 100)
										}
									}
								}

								var handbraked bool

								select {
								case handbraked = <-svm.handbrake:
								default:
								}

								if !svm.currentPlayer.Load().Player.IsPlaying() || handbraked {
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
}

// IsFXPlaying checks if fx player is currently active and playing.
func (svm *SoundFXManager) IsFXPlaying() bool {
	return svm.currentPlayer.Load() != nil
}

// StopFXPlaying makes fx player stop playing active sound.
func (svm *SoundFXManager) StopFXPlaying() {
	svm.handbrake <- true
}

// PushWithHandbrake pushes a new track immediately to the queue at the highest priority
// stopping previously playing track.
func (svm *SoundFXManager) PushWithHandbrake(name string) {
	go func() {
		if svm.currentPlayer.Load() != nil {
			svm.handbrake <- true
		}

		sound := loader.GetInstance().GetSoundFX(name)

		player, err := svm.audioContext.NewPlayerF32(sound)
		if err != nil {
			logging.GetInstance().Fatal(errors.Wrap(err, common.ErrSoundPlayerAccess.Error()).Error())
		}

		svm.processingBatch = svm.processingBatch[:0]

		svm.processingBatch = append([]*dto.FXSoundUnit{{
			Duration: time.Second * time.Duration(sound.Length()) / common.BytesPerSample / common.SampleRate,
			Player:   player,
		}}, svm.processingBatch...)
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
		handbrake:    make(chan bool),
	}
}
