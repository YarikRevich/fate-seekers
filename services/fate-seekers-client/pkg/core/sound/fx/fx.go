package fx

import (
	"fmt"
	"time"

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
	tempTickerPeriod = time.Millisecond * 200

	// Represents fadeout duration.
	fadeoutDuration = time.Second * 2

	// Represents volume decremention step.
	volumeDecrementor = 20.0
)

// SoundFXManager represents sound FX manager.
type SoundFXManager struct {
	// Represents audio context used for stream creation.
	audioContext *audio.Context

	// Represents time ticker used for player processing batch.
	ticker *time.Ticker

	// Represents processing player batch.
	processingBatch []*dto.FXSoundUnit

	// Represents check if initial current player volume setup was performed.
	currentPlayerVolumeConfigured bool

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

					go func() {
						tempTicker := time.NewTicker(tempTickerPeriod)

						var volume float64

						for svm.currentPlayer.Player.Volume()*100 != 100 {
							select {
							case <-tempTicker.C:
								volume += volumeDecrementor

								fmt.Println(volume/100, svm.currentPlayer)

								svm.currentPlayer.Player.SetVolume(volume / 100)
							}
						}

						svm.currentPlayerVolumeConfigured = true

						tempTicker.Stop()
					}()

					svm.currentPlayer.Player.Play()

					go func() {
						tempTicker := time.NewTicker(tempTickerPeriod)

						for {
							select {
							case <-tempTicker.C:
								if svm.currentPlayer.Duration-svm.currentPlayer.Player.Position() <= fadeoutDuration && svm.currentPlayerVolumeConfigured {
									if svm.currentPlayer.Player.Volume() != 0 {
										svm.currentPlayer.Player.SetVolume(
											((svm.currentPlayer.Player.Volume() * 100) - volumeDecrementor) / 100)
									}

									if !svm.currentPlayer.Player.IsPlaying() {
										if err := svm.currentPlayer.Player.Close(); err != nil {
											logging.GetInstance().Fatal(errors.Wrap(err, common.ErrSoundPlayerAccess.Error()).Error())
										}
									}

									tempTicker.Stop()

									svm.currentPlayer = nil

									svm.currentPlayerVolumeConfigured = false

									break
								}
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
