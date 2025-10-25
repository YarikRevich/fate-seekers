package session

import (
	"fmt"
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/shader/event/toxicrain"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition/transparent"
	contentstream "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/content/stream"
	metadatav1 "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/api"
	metadatastream "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/stream"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/options"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/builder"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/notification"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/translation"
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

	// Represents session toxic rain event shader effect.
	toxicRainEventShaderEffect *toxicrain.ToxicRainEventEffect
}

// // DirectionFromPoints возвращает направление движения от (px,py) к (x,y).
// // eps — dead zone: если |dx| и |dy| ≤ eps, вернёт "" (нет движения).
// func DirectionFromPoints(px, py, x, y, eps float64) string {
// 	dx := x - px
// 	dy := y - py

// 	if math.Abs(dx) <= eps && math.Abs(dy) <= eps {
// 		return "" // no move
// 	}

// 	ang := math.Atan2(dy, dx) // (-π, π]

// 	step := math.Pi / 4 // 45°
// 	idx := int(math.Round(ang / step))

// 	idx = ((idx % 8) + 8) % 8

// 	switch idx {
// 	case 0:
// 		return RightMovableRotation
// 	case 1:
// 		return UpRightMovableRotation
// 	case 2:
// 		return UpMovableRotation
// 	case 3:
// 		return UpLeftMovableRotation
// 	case 4:
// 		return LeftMovableRotation
// 	case 5:
// 		return DownLeftMovableRotation
// 	case 6:
// 		return DownMovableRotation
// 	case 7:
// 		return DownRightMovableRotation
// 	default:
// 		return ""
// 	}
// }

func (ss *SessionScreen) HandleInput() error {
	if store.GetUpdateUserMetadataPositionsStartedNetworking() == value.UPDATE_USER_METADATA_POSITIONS_STARTED_NETWORKING_FALSE_VALUE {
		dispatcher.GetInstance().Dispatch(
			action.NewSetUpdateUserMetadataPositionsStartedNetworking(
				value.UPDATE_USER_METADATA_POSITIONS_STARTED_NETWORKING_TRUE_VALUE))

		contentstream.GetUpdateUserMetadataPositionsSubmitter().Clean(func() {
			contentstream.GetUpdateUserMetadataPositionsSubmitter().Submit(
				store.GetSelectedLobbySetUnitMetadata().ID, func(err error) bool {
					if store.GetActiveScreen() != value.ACTIVE_SCREEN_SESSION_VALUE {
						dispatcher.GetInstance().Dispatch(
							action.NewSetUpdateUserMetadataPositionsStartedNetworking(
								value.UPDATE_USER_METADATA_POSITIONS_STARTED_NETWORKING_FALSE_VALUE))

						return true
					}

					if err != nil {
						notification.GetInstance().Push(
							common.ComposeMessage(
								translation.GetInstance().GetTranslation("client.networking.update-user-metadata-positions-failure"),
								err.Error()),
							time.Second*3,
							common.NotificationErrorTextColor)

						return true
					}

					return false
				})
		})
	}

	if store.GetEventRetrievalStartedNetworking() == value.EVENT_RETRIEVAL_STARTED_NETWORKING_FALSE_STATE {
		dispatcher.GetInstance().Dispatch(
			action.NewSetEventRetrievalStartedNetworking(
				value.EVENT_RETRIEVAL_STARTED_NETWORKING_TRUE_STATE))

		metadatastream.GetGetEventsSubmitter().Clean(func() {
			metadatastream.GetGetEventsSubmitter().Submit(
				store.GetSelectedSessionMetadata().ID, func(response *metadatav1.GetEventsResponse, err error) bool {
					if store.GetActiveScreen() != value.ACTIVE_SCREEN_SESSION_VALUE {
						dispatcher.GetInstance().Dispatch(
							action.NewSetEventRetrievalStartedNetworking(
								value.EVENT_RETRIEVAL_STARTED_NETWORKING_FALSE_STATE))

						return true
					}

					if err != nil {
						notification.GetInstance().Push(
							common.ComposeMessage(
								translation.GetInstance().GetTranslation("client.networking.event-retrieval-failure"),
								err.Error()),
							time.Second*3,
							common.NotificationErrorTextColor)

						return true
					}

					if len(response.GetName()) != 0 {
						switch response.GetName() {
						case value.EVENT_NAME_TOXIC_RAIN_VALUE:
							if store.GetEventName() != value.EVENT_NAME_TOXIC_RAIN_VALUE {
								notification.GetInstance().Push(
									translation.GetInstance().GetTranslation("client.networking.event-toxic-rain-starated"),
									time.Second*3,
									common.NotificationInfoTextColor)

								dispatcher.GetInstance().Dispatch(
									action.NewSetEventName(value.EVENT_NAME_TOXIC_RAIN_VALUE))
							}
						}
					} else if store.GetEventName() != value.EVENT_NAME_EMPTY_VALUE {
						notification.GetInstance().Push(
							translation.GetInstance().GetTranslation("client.networking.event-finished"),
							time.Second*3,
							common.NotificationInfoTextColor)

						dispatcher.GetInstance().Dispatch(
							action.NewSetEventEnding(value.EVENT_ENDING_TRUE_VALUE))
					}

					return false
				})
		})
	}

	if store.GetUsersMetadataRetrievalStartedNetworking() == value.USERS_METADATA_RETRIEVAL_STARTED_NETWORKING_FALSE_STATE {
		dispatcher.GetInstance().Dispatch(
			action.NewSetUsersMetadataRetrievalStartedNetworking(
				value.USERS_METADATA_RETRIEVAL_STARTED_NETWORKING_TRUE_STATE))

		metadatastream.GetGetUsersMetadataSubmitter().Clean(func() {
			metadatastream.GetGetUsersMetadataSubmitter().Submit(
				store.GetSelectedSessionMetadata().ID, func(response *metadatav1.GetUsersMetadataResponse, err error) bool {
					if store.GetActiveScreen() != value.ACTIVE_SCREEN_SESSION_VALUE {
						dispatcher.GetInstance().Dispatch(
							action.NewSetUsersMetadataRetrievalStartedNetworking(
								value.USERS_METADATA_RETRIEVAL_STARTED_NETWORKING_FALSE_STATE))

						return true
					}

					fmt.Println(response, err)

					if err != nil {
						notification.GetInstance().Push(
							common.ComposeMessage(
								translation.GetInstance().GetTranslation("client.networking.users-metadata-retrieval-failure"),
								err.Error()),
							time.Second*3,
							common.NotificationErrorTextColor)

						return true
					}

					return false
				})
		})
	}

	{
		// if !sound.GetInstance().GetSoundMusicManager().IsMusicPlaying() {
		// 	sound.GetInstance().GetSoundMusicManager().StartMusic(loader.EnergetykMusicSound)
		// }

		// if sound.GetInstance().GetSoundMusicManager().IsMusicPlaying() &&
		// 	!sound.GetInstance().GetSoundMusicManager().IsMusicStopping() {
		// 	sound.GetInstance().GetSoundMusicManager().StopMusic()
		// }
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		dispatcher.GetInstance().Dispatch(
			action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_RESUME_VALUE))
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		dispatcher.GetInstance().Dispatch(action.NewDecrementXPositionSession())

	}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		dispatcher.GetInstance().Dispatch(action.NewIncrementYPositionSession())

	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		dispatcher.GetInstance().Dispatch(action.NewDecrementYPositionSession())

	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		dispatcher.GetInstance().Dispatch(action.NewIncrementXPositionSession())
	}

	// TODO: update animator.

	if !ss.transparentTransitionEffect.Done() {
		if !ss.transparentTransitionEffect.OnEnd() {
			ss.transparentTransitionEffect.Update()
		} else {
			ss.transparentTransitionEffect.Clean()
		}
	}

	ss.ui.Update()

	if store.GetEventName() != value.EVENT_NAME_EMPTY_VALUE {
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
		} else if store.GetEventStarted() == value.EVENT_STARTED_FALSE_VALUE {
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
		}
	}

	// TODO: click on the letter.
	// dispatcher.GetInstance().Dispatch(action.NewSetLetterNameAction(""))

	// dispatcher.GetInstance().Dispatch(action.NewSetLetterImageAction(""))

	// TODO: click on the chest.
	// dispatcher.GetInstance().Dispatch(action.New)

	return nil
}

// objects
// map(may include )

func (ss *SessionScreen) HandleRender(screen *ebiten.Image) {
	ss.world.Clear()

	if store.GetEventName() != value.EVENT_NAME_EMPTY_VALUE {
		ss.eventWorld.Clear()
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
		toxicRainEventShaderEffect: toxicrain.NewToxicRainEventEffect(),
	}
}
