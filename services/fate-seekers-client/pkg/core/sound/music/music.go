package music

import (
	"math"
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

var (
	ErrMusicIsNotPlaying     = errors.New("err happened music is not playing")
	ErrMusicIsAlreadyPlaying = errors.New("err happened music is already playing")
)

const (
	// Represents processing ticker duration.
	processingTickerPeriod = time.Millisecond * 400

	// Represents start temp ticker period.
	startTempTickerPeriod = time.Millisecond * 10

	// Represents end temp ticker period.
	endTempTickerPeriod = time.Millisecond * 120

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
	currentAmbientPlayer *dto.AmbientSoundUnit

	// Represents currently playing music player.
	currentMusicPlayer *dto.MusicSoundUnit
}

// Init starts sound music manager processing worker. Additionally starts volume update worker.
func (smm *SoundMusicManager) Init() {
	go func() {
		for {
			select {
			case <-smm.ticker.C:
				if len(smm.ambientProcessingBatch) > 0 && smm.currentAmbientPlayer == nil {
					smm.currentAmbientPlayer = smm.ambientProcessingBatch[0]

					smm.ambientProcessingBatch = append(smm.ambientProcessingBatch[:0], smm.ambientProcessingBatch[1:]...)

					sound := loader.GetInstance().GetSoundMusic(smm.currentAmbientPlayer.Name)

					player, err := smm.audioContext.NewPlayerF32(sound)
					if err != nil {
						logging.GetInstance().Fatal(errors.Wrap(err, common.ErrSoundPlayerAccess.Error()).Error())
					}

					smm.ambientProcessingBatch = append(smm.ambientProcessingBatch, &dto.AmbientSoundUnit{
						Name:     smm.currentAmbientPlayer.Name,
						Duration: time.Second * time.Duration(sound.Length()) / common.BytesPerSample / common.SampleRate,
						Player:   player,
					})

					smm.currentAmbientPlayer.Player.SetVolume(0)

					interruptionChan := make(chan int, 1)

					go func() {
						tempTicker := time.NewTicker(startTempTickerPeriod)

						var volume float64

						for smm.currentAmbientPlayer != nil && smm.currentAmbientPlayer.Player.Volume()*float64(config.GetSettingsSoundMusic()) < float64(config.GetSettingsSoundMusic()) {
							if smm.ambientProcessingInterruption.Load() {
								tempTicker.Stop()

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

								smm.currentAmbientPlayer.Player.SetVolume(volume / float64(config.GetSettingsSoundMusic()))

								tempTicker.Reset(startTempTickerPeriod)
							}
						}
					}()

					smm.currentAmbientPlayer.Player.Play()

					go func() {
						tempTicker := time.NewTicker(endTempTickerPeriod)

						for {
							select {
							case <-tempTicker.C:
								tempTicker.Stop()

								if smm.currentAmbientPlayer.Duration-smm.currentAmbientPlayer.Player.Position() <= fadeoutDuration && !smm.ambientProcessingInterruption.Load() {
									select {
									case interruptionChan <- 1:
									default:
									}

									if smm.currentAmbientPlayer.Player.Volume() != 0 {
										remainingTime := smm.currentAmbientPlayer.Duration - smm.currentAmbientPlayer.Player.Position()

										normalized := float64(remainingTime.Milliseconds()) / float64(fadeoutDuration.Milliseconds())

										fadeFactor := math.Max(0, math.Sin((math.Pi/2)*normalized))

										newVolume := float64(smm.currentAmbientPlayer.Player.Volume()) * fadeFactor

										if smm.currentAmbientPlayer.Player.Volume() >= 0 {
											smm.currentAmbientPlayer.Player.SetVolume(newVolume)
										}
									}
								}

								if !smm.currentAmbientPlayer.Player.IsPlaying() {
									if err := smm.currentAmbientPlayer.Player.Close(); err != nil {
										logging.GetInstance().Fatal(errors.Wrap(err, common.ErrSoundPlayerAccess.Error()).Error())
									}

									tempTicker.Stop()

									smm.currentAmbientPlayer = nil

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
	return smm.currentMusicPlayer != nil
}

// StartMusic starts playing music stream of the provided name.
func (smm *SoundMusicManager) StartMusic(name string) {
	if smm.currentMusicPlayer != nil {
		logging.GetInstance().Fatal(ErrMusicIsAlreadyPlaying.Error())
	}

	sound := loader.GetInstance().GetSoundMusic(name)

	player, err := smm.audioContext.NewPlayerF32(sound)
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, common.ErrSoundPlayerAccess.Error()).Error())
	}

	smm.currentMusicPlayer = &dto.MusicSoundUnit{
		Duration: time.Second * time.Duration(sound.Length()) / common.BytesPerSample / common.SampleRate,
		Player:   player,
	}

	go func() {
		tempTicker := time.NewTicker(startTempTickerPeriod)

		if smm.currentAmbientPlayer != nil {
			smm.ambientProcessingInterruption.Store(true)

			for smm.currentAmbientPlayer.Player.Volume() != 0 {
				select {
				case <-tempTicker.C:
					remainingTime := smm.currentAmbientPlayer.Duration - smm.currentAmbientPlayer.Player.Position()

					normalized := float64(remainingTime.Milliseconds()) / float64(fadeoutDuration.Milliseconds())

					fadeFactor := math.Max(0, math.Sin((math.Pi/2)*normalized))

					newVolume := float64(smm.currentAmbientPlayer.Player.Volume()) * fadeFactor

					if smm.currentAmbientPlayer.Player.Volume() >= 0 {
						smm.currentAmbientPlayer.Player.SetVolume(newVolume)
					}
				}
			}

			smm.ambientProcessingInterruption.Store(false)
		}

		smm.currentMusicPlayer.Player.SetVolume(0)

		smm.currentMusicPlayer.Player.Play()

		interruptionChan := make(chan int, 1)

		var volume float64

		for smm.currentMusicPlayer != nil && smm.currentMusicPlayer.Player.Volume()*float64(config.GetSettingsSoundMusic()) < float64(config.GetSettingsSoundMusic()) {
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

				smm.currentMusicPlayer.Player.SetVolume(volume / float64(config.GetSettingsSoundMusic()))

				tempTicker.Reset(startTempTickerPeriod)
			}
		}

		tempTicker.Stop()

		go func() {
			tempTicker := time.NewTicker(endTempTickerPeriod)

			for {
				select {
				case <-tempTicker.C:
					tempTicker.Stop()

					if smm.currentAmbientPlayer.Duration-smm.currentAmbientPlayer.Player.Position() <= fadeoutDuration {
						if !smm.ambientProcessingInterruption.Load() {
							smm.ambientProcessingInterruption.Store(true)
						}

						select {
						case interruptionChan <- 1:
						default:
						}

						if smm.currentAmbientPlayer.Player.Volume() != 0 {
							remainingTime := smm.currentAmbientPlayer.Duration - smm.currentAmbientPlayer.Player.Position()

							normalized := float64(remainingTime.Milliseconds()) / float64(fadeoutDuration.Milliseconds())

							fadeFactor := math.Max(0, math.Sin((math.Pi/2)*normalized))

							newVolume := float64(smm.currentAmbientPlayer.Player.Volume()) * fadeFactor

							if smm.currentAmbientPlayer.Player.Volume() >= 0 {
								smm.currentAmbientPlayer.Player.SetVolume(newVolume)
							}
						}
					}

					if !smm.currentAmbientPlayer.Player.IsPlaying() {
						if err := smm.currentAmbientPlayer.Player.Close(); err != nil {
							logging.GetInstance().Fatal(errors.Wrap(err, common.ErrSoundPlayerAccess.Error()).Error())
						}

						tempTicker.Stop()

						smm.ambientProcessingInterruption.Store(false)

						smm.currentAmbientPlayer = nil

						return
					}

					tempTicker.Reset(endTempTickerPeriod)
				}
			}
		}()
	}()
}

// StopMusic stops currently playing music stream.
func (smm *SoundMusicManager) StopMusic() {
	if smm.currentMusicPlayer == nil {
		logging.GetInstance().Fatal(ErrMusicIsNotPlaying.Error())
	}

	go func() {
		tempTicker := time.NewTicker(endTempTickerPeriod)

		for smm.currentMusicPlayer.Player.Volume() != 0 {
			select {
			case <-tempTicker.C:
				remainingTime := smm.currentMusicPlayer.Duration - smm.currentMusicPlayer.Player.Position()

				normalized := float64(remainingTime.Milliseconds()) / float64(fadeoutDuration.Milliseconds())

				fadeFactor := math.Max(0, math.Sin((math.Pi/2)*normalized))

				newVolume := float64(smm.currentMusicPlayer.Player.Volume()) * fadeFactor

				if smm.currentMusicPlayer.Player.Volume() >= 0 {
					smm.currentMusicPlayer.Player.SetVolume(newVolume)
				}
			}
		}

		smm.currentMusicPlayer.Player.Pause()

		volume := smm.currentAmbientPlayer.Player.Volume() * float64(config.GetSettingsSoundMusic())

		smm.ambientProcessingInterruption.Store(true)

		for smm.currentAmbientPlayer.Player.Volume()*float64(config.GetSettingsSoundMusic()) != float64(config.GetSettingsSoundMusic()) {
			select {
			case <-tempTicker.C:
				volume += volumeShift

				smm.currentAmbientPlayer.Player.SetVolume(volume / float64(config.GetSettingsSoundMusic()))
			}
		}

		smm.ambientProcessingInterruption.Store(false)

		if err := smm.currentMusicPlayer.Player.Close(); err != nil {
			logging.GetInstance().Fatal(errors.Wrap(err, common.ErrSoundPlayerAccess.Error()).Error())
		}

		tempTicker.Stop()

		smm.currentMusicPlayer = nil
	}()
}

// NewSoundMusicManager initializes SoundMusicManager.
func NewSoundMusicManager(audioContext *audio.Context) *SoundMusicManager {
	return &SoundMusicManager{
		audioContext: audioContext,
		ticker:       time.NewTicker(processingTickerPeriod),
	}
}
