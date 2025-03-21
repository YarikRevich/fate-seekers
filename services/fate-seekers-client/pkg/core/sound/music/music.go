package music

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

var (
	ErrMusicIsNotPlaying     = errors.New("err happened music is not playing")
	ErrMusicIsAlreadyPlaying = errors.New("err happened music is already playing")
)

const (
	// Represents processing ticker duration.
	processingTickerPeriod = time.Millisecond * 400

	// Represents temp ticker period.
	tempTickerPeriod = time.Millisecond * 200

	// Represents fadeout duration.
	fadeoutDuration = time.Second * 5

	// Represents volume decremention step.
	volumeDecrementor = 20.0

	// Represents ambient suspension volume level, given in percentage.
	ambientSuspensionVolume = 30
)

// SoundMusicManager represents sound music manager, used for both ambient and music streams management.
type SoundMusicManager struct {
	// Represents audio context used for stream creation.
	audioContext *audio.Context

	// Represents time ticker used for player processing batch.
	ticker *time.Ticker

	// Represents processing ambient player infinite batch.
	ambientProcessingBatch []*dto.AmbientSoundUnit

	// Represents check if initial current ambient player volume setup was performed.
	currentAmbientPlayerVolumeConfigured bool

	// Represents currently playing ambient player.
	currentAmbientPlayer *dto.AmbientSoundUnit

	// Represents check if initial current music player volume setup was performed.
	currentMusicPlayerVolumeConfigured bool

	// Represents currently playing music player.
	currentMusicPlayer *dto.MusicSoundUnit
}

// Init starts sound music manager processing worker.
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

					go func() {
						tempTicker := time.NewTicker(tempTickerPeriod)

						var volume float64

						for smm.currentAmbientPlayer.Player.Volume()*100 != 100 {
							select {
							case <-tempTicker.C:
								volume += volumeDecrementor

								smm.currentAmbientPlayer.Player.SetVolume(volume / 100)
							}
						}

						smm.currentAmbientPlayerVolumeConfigured = true

						tempTicker.Stop()
					}()

					smm.currentAmbientPlayer.Player.Play()

					go func() {
						tempTicker := time.NewTicker(tempTickerPeriod)

						for {
							select {
							case <-tempTicker.C:
								if smm.currentAmbientPlayer.Duration-smm.currentAmbientPlayer.Player.Position() <= fadeoutDuration && smm.currentAmbientPlayerVolumeConfigured {
									if smm.currentAmbientPlayer.Player.Volume() != 0 {
										smm.currentAmbientPlayer.Player.SetVolume(
											((smm.currentAmbientPlayer.Player.Volume() * 100) - volumeDecrementor) / 100)
									}

									if !smm.currentAmbientPlayer.Player.IsPlaying() {
										if err := smm.currentAmbientPlayer.Player.Close(); err != nil {
											logging.GetInstance().Fatal(errors.Wrap(err, common.ErrSoundPlayerAccess.Error()).Error())
										}
									}

									tempTicker.Stop()

									smm.currentAmbientPlayer = nil

									smm.currentAmbientPlayerVolumeConfigured = false

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
		tempTicker := time.NewTicker(tempTickerPeriod)

		smm.currentMusicPlayer.Player.SetVolume(0)

		var volume float64

		if smm.currentAmbientPlayer != nil {
			volume := smm.currentAmbientPlayer.Player.Volume() * 100

			for smm.currentAmbientPlayer.Player.Volume()*100 != ambientSuspensionVolume {
				select {
				case <-tempTicker.C:
					if volume-volumeDecrementor >= ambientSuspensionVolume {
						volume -= volumeDecrementor
					} else {
						volume = ambientSuspensionVolume
					}

					smm.currentAmbientPlayer.Player.SetVolume(volume / 100)
				}
			}
		}

		smm.currentMusicPlayer.Player.Play()

		volume = smm.currentMusicPlayer.Player.Volume() * 100

		for smm.currentMusicPlayer.Player.Volume()*100 != 100 {
			select {
			case <-tempTicker.C:
				volume += volumeDecrementor

				fmt.Println(volume)

				smm.currentMusicPlayer.Player.SetVolume(volume / 100)
			}
		}

		smm.currentMusicPlayerVolumeConfigured = true

		tempTicker.Stop()
	}()

	go func() {
		tempTicker := time.NewTicker(tempTickerPeriod)

		for {
			select {
			case <-tempTicker.C:
				if smm.currentMusicPlayer.Duration-smm.currentMusicPlayer.Player.Position() <= fadeoutDuration && smm.currentMusicPlayerVolumeConfigured {
					if smm.currentMusicPlayer.Player.Volume() != 0 {
						smm.currentMusicPlayer.Player.SetVolume(
							((smm.currentMusicPlayer.Player.Volume() * 100) - volumeDecrementor) / 100)
					}

					if !smm.currentMusicPlayer.Player.IsPlaying() {
						if err := smm.currentMusicPlayer.Player.Close(); err != nil {
							logging.GetInstance().Fatal(errors.Wrap(err, common.ErrSoundPlayerAccess.Error()).Error())
						}
					}

					tempTicker.Stop()

					smm.currentMusicPlayer = nil

					smm.currentMusicPlayerVolumeConfigured = false

					break
				}
			}
		}
	}()
}

// StopMusic stops currently playing music stream.
func (smm *SoundMusicManager) StopMusic() {
	if smm.currentMusicPlayer == nil {
		logging.GetInstance().Fatal(ErrMusicIsNotPlaying.Error())
	}

	go func() {
		tempTicker := time.NewTicker(tempTickerPeriod)

		volume := smm.currentMusicPlayer.Player.Volume() * 100

		for smm.currentMusicPlayer.Player.Volume() != 0 {
			select {
			case <-tempTicker.C:
				volume -= volumeDecrementor

				smm.currentAmbientPlayer.Player.SetVolume(volume / 100)
			}
		}

		smm.currentMusicPlayer.Player.Pause()

		volume = smm.currentAmbientPlayer.Player.Volume() * 100

		for smm.currentAmbientPlayer.Player.Volume()*100 != 100 {
			select {
			case <-tempTicker.C:
				volume += volumeDecrementor

				smm.currentAmbientPlayer.Player.SetVolume(volume / 100)
			}
		}

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
