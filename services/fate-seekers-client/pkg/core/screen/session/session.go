package session

import (
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/particle"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/particle/loadingstars"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/shader/event/toxicrain"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition/transparent"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/sound"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/options"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/builder"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
)

// const (
// 	tileSize   = 16
// 	tileXCount = 25
// )
// worldSizeX  = worldWidth / tileSize

// var (
// 	tilesImage *ebiten.Image
// )

// func init() {
// 	// Decode an image from the image file's byte slice.
// 	img, _, err := image.Decode(bytes.NewReader(images.Tiles_png))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	tilesImage = ebiten.NewImageFromImage(img)
// }

// for _, l := range r.layers {
// 	for i, t := range l {
// 		op := &ebiten.DrawImageOptions{}
// 		op.GeoM.Translate(float64((i%worldSizeX)*tileSize), float64((i/worldSizeX)*tileSize))

// 		sx := (t % tileXCount) * tileSize
// 		sy := (t / tileXCount) * tileSize
// 		r.world.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
// 	}
// 	}

// if ebiten.IsKeyPressed(ebiten.KeyU) {
// 	ms.animation.Update()
// }

// if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
// 	ms.camera.TranslatePositionX(-1)
// }
// if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
// 	ms.camera.TranslatePositionX(1)
// }
// if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
// 	ms.camera.TranslatePositionY(-1)
// }
// if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
// 	ms.camera.TranslatePositionY(1)
// }

// if ebiten.IsKeyPressed(ebiten.KeyQ) {
// 	if ms.camera.GetZoom() > -2400 {
// 		ms.camera.ZoomOutBy(10)
// 	}
// }
// if ebiten.IsKeyPressed(ebiten.KeyE) {
// 	if ms.camera.GetZoom() < 2400 {
// 		ms.camera.ZoomInBy(10)
// 	}
// }

// if ebiten.IsKeyPressed(ebiten.KeyR) {
// 	ms.camera.RotateLeft()
// }

// if ebiten.IsKeyPressed(ebiten.KeyT) {
// 	ms.camera.RotateRight()
// }

// if ebiten.IsKeyPressed(ebiten.KeySpace) {
// 	ms.camera.Reset()
// }

var (
	// GetInstance retrieves instance of the session screen, performing initilization if needed.
	GetInstance = sync.OnceValue[screen.Screen](newSessionScreen)
)

// SessionScreen represents session screen implementation.
type SessionScreen struct {
	// Represents attached user interface.
	ui *ebitenui.UI

	// Represents attached pressed user interface.
	pressedInterface *ebitenui.UI

	// Represents transparent transition effect.
	transparentTransitionEffect transition.TransitionEffect

	// Represents transparent transition effect used for toxic rain event component, when event is started.
	toxicRainEventStartTransparentTransitionEffect transition.TransitionEffect

	// Represents transparent transition effect used for toxic rain event component, when event is ended.
	toxicRainEventEndTransparentTransitionEffect transition.TransitionEffect

	// Represents global world view.
	world *ebiten.Image

	// Represents event world view.
	eventWorld *ebiten.Image

	// Represents session loading stars particle effect.
	loadingStarsParticleEffect particle.ParticleEffect

	// Represents session toxic rain event shader effect.
	toxicRainEventShaderEffect *toxicrain.ToxicRainEventEffect
}

func (ss *SessionScreen) HandleInput() error {
	// Debug statements
	{
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			if !sound.GetInstance().GetSoundMusicManager().IsMusicPlaying() {
				sound.GetInstance().GetSoundMusicManager().StartMusic(loader.EnergetykMusicSound)
			}
		}

		if ebiten.IsKeyPressed(ebiten.KeyY) {
			if sound.GetInstance().GetSoundMusicManager().IsMusicPlaying() &&
				!sound.GetInstance().GetSoundMusicManager().IsMusicStopping() {
				sound.GetInstance().GetSoundMusicManager().StopMusic()
			}
		}

		if ebiten.IsKeyPressed(ebiten.KeyD) {
			dispatcher.GetInstance().Dispatch(action.NewSetEventName(value.EVENT_NAME_TOXIC_RAIN_VALUE))
		} else {
			dispatcher.GetInstance().Dispatch(action.NewSetEventEnding(value.EVENT_ENDING_TRUE_VALUE))
		}
	}

	if !ss.transparentTransitionEffect.Done() {
		if !ss.transparentTransitionEffect.OnEnd() {
			ss.transparentTransitionEffect.Update()
		} else {
			ss.transparentTransitionEffect.Clean()
		}
	}

	ss.ui.Update()

	if !ss.loadingStarsParticleEffect.Done() {
		if !ss.loadingStarsParticleEffect.OnEnd() {
			ss.loadingStarsParticleEffect.Update()
		} else {
			ss.loadingStarsParticleEffect.Clean()
		}
	}

	if store.GetEventName() != value.EVENT_NAME_EMPTY_VALUE {
		if store.GetEventStarted() == value.EVENT_STARTED_FALSE_VALUE {
			switch store.GetEventName() {
			case value.EVENT_NAME_TOXIC_RAIN_VALUE:
				if !ss.toxicRainEventStartTransparentTransitionEffect.Done() {
					if !ss.toxicRainEventStartTransparentTransitionEffect.OnEnd() {
						ss.toxicRainEventStartTransparentTransitionEffect.Update()
					} else {
						ss.toxicRainEventStartTransparentTransitionEffect.Clean()

						dispatcher.GetInstance().Dispatch(
							action.NewSetEventStarted(value.EVENT_STARTED_TRUE_VALUE))
					}
				}
			}
		} else {
			if store.GetEventEnding() == value.EVENT_ENDING_TRUE_VALUE {
				switch store.GetEventName() {
				case value.EVENT_NAME_TOXIC_RAIN_VALUE:
					if !ss.toxicRainEventEndTransparentTransitionEffect.Done() {
						if !ss.toxicRainEventEndTransparentTransitionEffect.OnEnd() {
							ss.toxicRainEventEndTransparentTransitionEffect.Update()
						} else {
							dispatcher.GetInstance().Dispatch(
								action.NewSetEventName(value.EVENT_NAME_EMPTY_VALUE))

							dispatcher.GetInstance().Dispatch(
								action.NewSetEventStarted(value.EVENT_STARTED_FALSE_VALUE))

							dispatcher.GetInstance().Dispatch(
								action.NewSetEventEnding(value.EVENT_ENDING_FALSE_VALUE))

							ss.toxicRainEventStartTransparentTransitionEffect.Reset()

							ss.toxicRainEventEndTransparentTransitionEffect.Reset()
						}
					}
				}
			}

		}
	}

	// TODO: click on the letter.
	// dispatcher.GetInstance().Dispatch(action.NewSetLetterNameAction(""))

	// dispatcher.GetInstance().Dispatch(action.NewSetLetterImageAction(""))

	// TODO: click on the chest.
	// dispatcher.GetInstance().Dispatch(action.New)

	return nil
}

func (ss *SessionScreen) HandleNetworking() {

}

func (ss *SessionScreen) HandleRender(screen *ebiten.Image) {
	ss.world.Clear()

	if store.GetEventName() != value.EVENT_NAME_EMPTY_VALUE {
		ss.eventWorld.Clear()
	}

	if !ss.loadingStarsParticleEffect.Done() {
		ss.loadingStarsParticleEffect.Draw(screen)
	}

	ss.ui.Draw(ss.world)

	screen.DrawImage(ss.world, &ebiten.DrawImageOptions{
		ColorM: options.GetTransparentDrawOptions(
			ss.transparentTransitionEffect.GetValue()).ColorM})

	if store.GetEventName() != value.EVENT_NAME_EMPTY_VALUE {
		switch store.GetEventName() {
		case value.EVENT_NAME_TOXIC_RAIN_VALUE:
			if store.GetEventStarted() == value.EVENT_STARTED_TRUE_VALUE && store.GetEventEnding() == value.EVENT_ENDING_TRUE_VALUE {
				ss.toxicRainEventShaderEffect.Draw(
					ss.eventWorld, ss.toxicRainEventEndTransparentTransitionEffect.GetValue())
			} else {
				ss.toxicRainEventShaderEffect.Draw(
					ss.eventWorld, ss.toxicRainEventStartTransparentTransitionEffect.GetValue())
			}
		}

		screen.DrawImage(ss.eventWorld, &ebiten.DrawImageOptions{})
	}
}

func (ss *SessionScreen) Clean() {

}

// newSessionScreen initializes SessionScreen.
func newSessionScreen() screen.Screen {
	return &SessionScreen{
		ui: builder.Build(),
		transparentTransitionEffect: transparent.NewTransparentTransitionEffect(
			true, 255, 0, 5, time.Microsecond*10),
		toxicRainEventStartTransparentTransitionEffect: transparent.NewTransparentTransitionEffect(
			true, 10, 0, 0.5, time.Millisecond*200),
		toxicRainEventEndTransparentTransitionEffect: transparent.NewTransparentTransitionEffect(
			false, 0, 10, 0.5, time.Millisecond*200),
		world:                      ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
		eventWorld:                 ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
		loadingStarsParticleEffect: loadingstars.NewStarsParticleEffect(),
		toxicRainEventShaderEffect: toxicrain.NewToxicRainEventEffect(),
	}
}
