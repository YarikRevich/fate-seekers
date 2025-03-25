package fx

import (
	"math"
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
	processingTickerPeriod = time.Millisecond * 400

	// Represents temp ticker period.
	startTempTickerPeriod = time.Millisecond * 10

	// Represents temp ticker period.
	endTempTickerPeriod = time.Millisecond * 120

	// Represents fadeout duration.
	fadeoutDuration = time.Millisecond * 1300

	// Represents volume decremention step.
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
	currentPlayer *dto.FXSoundUnit
}

// Init starts sound FX manager processing worker.
func (svm *SoundFXManager) Init() {
	go func() {
		for {
			select {
			case <-svm.ticker.C:
				if len(svm.processingBatch) > 0 && svm.currentPlayer == nil {
					svm.currentPlayer = svm.processingBatch[0]

					svm.processingBatch = append(svm.processingBatch[:0], svm.processingBatch[1:]...)

					svm.currentPlayer.Player.SetVolume(0)

					interruptionChan := make(chan int, 1)

					go func() {
						tempTicker := time.NewTicker(startTempTickerPeriod)

						var volume float64

						for svm.currentPlayer != nil && svm.currentPlayer.Player.Volume()*float64(config.GetSettingsSoundFX()) < float64(config.GetSettingsSoundFX()) {
							select {
							case <-interruptionChan:
								tempTicker.Stop()

								return
							case <-tempTicker.C:
								tempTicker.Stop()

								if volume+volumeShift <= float64(config.GetSettingsSoundFX()) {
									volume += volumeShift
								} else {
									volume = float64(config.GetSettingsSoundFX())
								}

								svm.currentPlayer.Player.SetVolume(volume / float64(config.GetSettingsSoundFX()))

								tempTicker.Reset(startTempTickerPeriod)
							}
						}
					}()

					svm.currentPlayer.Player.Play()

					go func() {
						tempTicker := time.NewTicker(endTempTickerPeriod)

						for {
							select {
							case <-tempTicker.C:
								tempTicker.Stop()

								if svm.currentPlayer.Duration-svm.currentPlayer.Player.Position() <= fadeoutDuration {
									select {
									case interruptionChan <- 1:
									default:
									}

									if svm.currentPlayer.Player.Volume() != 0 {
										remainingTime := svm.currentPlayer.Duration - svm.currentPlayer.Player.Position()

										normalized := float64(remainingTime.Milliseconds()) / float64(fadeoutDuration.Milliseconds())

										fadeFactor := math.Max(0, math.Sin((math.Pi/2)*normalized))

										newVolume := float64(svm.currentPlayer.Player.Volume()) * fadeFactor

										if svm.currentPlayer.Player.Volume() > 0 {
											svm.currentPlayer.Player.SetVolume(newVolume)
										}
									}
								}

								if !svm.currentPlayer.Player.IsPlaying() {
									if err := svm.currentPlayer.Player.Close(); err != nil {
										logging.GetInstance().Fatal(errors.Wrap(err, common.ErrSoundPlayerAccess.Error()).Error())
									}

									tempTicker.Stop()

									svm.currentPlayer = nil

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
