package session

import (
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
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/sound"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/animation/animator"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/animation/animator/movable"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/animation/direction"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/options"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/builder"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/bar"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/notification"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/translation"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/kamera/v2"
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

var (
	// Represents shared users metadata issuers map.
	sharedUsersMetadataIssuers = make(map[string]bool)

	// Represents shared users metadata health hits issuers map.
	sharedUsersMetadataHealthHitsIssuers = make(map[string]bool)
)

// SessionScreen represents session screen implementation.
type SessionScreen struct {
	// Represents attached passive user interface.
	passiveUI *ebitenui.UI

	// Represents attached active user interface.
	activeUI *ebitenui.UI

	// Represents attached animator instance.
	animator *animator.Animator

	// Represents attached camera instance.
	camera *kamera.Camera

	// Represents attached pressed user interface.
	pressedInterface *ebitenui.UI

	// Represents transparent transition effect.
	transparentTransitionEffect transition.TransitionEffect

	// Represents transparent transition effect used for toxic rain event component, when event is started.
	toxicRainEventStartTransparentTransitionEffect transition.TransitionEffect

	// Represents transparent transition effect used for toxic rain event component, when event is ended.
	toxicRainEventEndTransparentTransitionEffect transition.TransitionEffect

	// Represents passive interface world view.
	passiveInterfaceWorld *ebiten.Image

	// Represents active interface world view.
	activeInterfaceWorld *ebiten.Image

	// Represents internal world view.
	internalWorld *ebiten.Image

	// Represents event world view.
	eventWorld *ebiten.Image

	// Represents session toxic rain event shader effect.
	toxicRainEventShaderEffect *toxicrain.ToxicRainEventEffect
}

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
					} else if store.GetEventName() != value.EVENT_NAME_EMPTY_VALUE &&
						store.GetEventEnding() != value.EVENT_ENDING_TRUE_VALUE {
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

					if err != nil {
						notification.GetInstance().Push(
							common.ComposeMessage(
								translation.GetInstance().GetTranslation("client.networking.users-metadata-retrieval-failure"),
								err.Error()),
							time.Second*3,
							common.NotificationErrorTextColor)

						return true
					}

					clear(sharedUsersMetadataIssuers)

					clear(sharedUsersMetadataHealthHitsIssuers)

					var animationDirection string

					for _, userMetadata := range response.GetUserMetadata() {
						if userMetadata.Issuer == store.GetRepositoryUUID() && userMetadata.Eliminated {
							dispatcher.GetInstance().Dispatch(
								action.NewSetEventRetrievalStartedNetworking(
									value.EVENT_RETRIEVAL_STARTED_NETWORKING_FALSE_STATE))

							dispatcher.GetInstance().Dispatch(
								action.NewSetUsersMetadataRetrievalStartedNetworking(
									value.USERS_METADATA_RETRIEVAL_STARTED_NETWORKING_FALSE_STATE))

							dispatcher.GetInstance().Dispatch(
								action.NewSetResetDeath(value.RESET_DEATH_TRUE_VALUE))

							dispatcher.GetInstance().Dispatch(
								action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_DEATH_VALUE))

							return true
						}

						store.RetrievedUsersMetadataSessionSyncHelper.Lock()

						retrievedUsersMetadataSession := store.GetRetrievedUsersMetadataSession()

						previousUsersMetadata, ok := retrievedUsersMetadataSession[userMetadata.GetIssuer()]
						if !ok {
							animationDirection = dto.RightMovableRotation
						} else {
							if previousUsersMetadata.Position.X != userMetadata.GetPosition().GetX() ||
								previousUsersMetadata.Position.Y != userMetadata.GetPosition().GetY() {
								animationDirection = direction.GetAnimationDirection(
									previousUsersMetadata.Position.X,
									previousUsersMetadata.Position.Y,
									userMetadata.GetPosition().GetX(),
									userMetadata.GetPosition().GetY())
							} else {
								animationDirection = previousUsersMetadata.AnimationDirection
							}

							if userMetadata.Issuer == store.GetRepositoryUUID() {
								if previousUsersMetadata.Health != userMetadata.Health {
									bar.GetInstance().SetHealthText(userMetadata.Health)
								}
							}

							if previousUsersMetadata.Health != userMetadata.Health {
								sharedUsersMetadataHealthHitsIssuers[userMetadata.GetIssuer()] = true
							}
						}

						retrievedUsersMetadataSession[userMetadata.GetIssuer()] =
							dto.RetrievedUsersMetadataSessionUnit{
								Health:             userMetadata.GetHealth(),
								Skin:               userMetadata.GetSkin(),
								Active:             userMetadata.GetActive(),
								Eliminated:         userMetadata.GetEliminated(),
								AnimationDirection: animationDirection,
								AnimationStatic:    userMetadata.GetStatic(),
								Position: dto.Position{
									X: userMetadata.GetPosition().GetX(),
									Y: userMetadata.GetPosition().GetY(),
								},
							}

						store.RetrievedUsersMetadataSessionSyncHelper.Unlock()

						if userMetadata.GetIssuer() != store.GetRepositoryUUID() {
							sharedUsersMetadataIssuers[userMetadata.GetIssuer()] = true
						}
					}

					ss.animator.GetMovables().PruneSecondary(sharedUsersMetadataIssuers)

					for issuer := range sharedUsersMetadataIssuers {
						retrievedUsersMetadata := store.GetRetrievedUsersMetadataSession()[issuer]

						if !ss.animator.GetMovables().SecondaryExists(issuer) {
							movableUnit := movable.NewMovableUnit(
								loader.GetMovableSkinsPath(retrievedUsersMetadata.Skin))

							movableUnit.SetDirection(retrievedUsersMetadata.AnimationDirection)
							movableUnit.SetStatic(retrievedUsersMetadata.AnimationStatic)
							movableUnit.AddPosition(retrievedUsersMetadata.Position)

							ss.animator.GetMovables().AddSecondary(issuer, movableUnit)
						} else {
							movableUnit := ss.animator.GetMovables().GetSecondary(issuer)

							movableUnit.SetDirection(retrievedUsersMetadata.AnimationDirection)
							movableUnit.SetStatic(retrievedUsersMetadata.AnimationStatic)
							movableUnit.AddPosition(retrievedUsersMetadata.Position)
						}
					}

					var movableUnit *movable.MovableUnit

					for issuer := range sharedUsersMetadataHealthHitsIssuers {
						if ss.animator.GetMovables().SecondaryExists(issuer) {
							movableUnit = ss.animator.GetMovables().GetSecondary(issuer)
						} else {
							if ss.animator.GetMovables().MainExists(issuer) {
								sound.GetInstance().GetSoundFxManager().PushWithHandbrake(loader.ToxicRainFXSound)

								movableUnit = ss.animator.GetMovables().GetMain(issuer)
							} else {
								continue
							}
						}

						movableUnit.TriggerNormalHit()
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

	store.RetrievedUsersMetadataSessionSyncHelper.Lock()

	retrievedUsersMetadataSession := store.GetRetrievedUsersMetadataSession()

	if _, ok := retrievedUsersMetadataSession[store.GetRepositoryUUID()]; ok {
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
	}

	if ebiten.IsKeyPressed(ebiten.KeyB) {
		ss.camera.AddTrauma(0.1)
	}

	ss.camera.LookAt(store.GetPositionSession().X, -store.GetPositionSession().Y)

	store.RetrievedUsersMetadataSessionSyncHelper.Unlock()

	selectedLobbySet := store.GetSelectedLobbySetUnitMetadata()

	if !ss.animator.GetMovables().MainExists(selectedLobbySet.Issuer) {
		movableUnit := movable.NewMovableUnit(
			loader.GetMovableSkinsPath(selectedLobbySet.Skin))

		movableUnit.SetDirection(dto.RightMovableRotation)
		movableUnit.SetCameraLock(true)
		movableUnit.SetStatic(true)
		// movableUnit.AddPosition(store.GetPositionSession())

		ss.animator.GetMovables().AddMain(selectedLobbySet.Issuer, movableUnit)
	} else {
		movableUnit := ss.animator.GetMovables().GetMain(selectedLobbySet.Issuer)

		if store.GetPreviousPositionSession().X != store.GetPositionSession().X ||
			store.GetPreviousPositionSession().Y != store.GetPositionSession().Y {
			movableUnit.SetDirection(direction.GetAnimationDirection(
				store.GetPreviousPositionSession().X,
				store.GetPreviousPositionSession().Y,
				store.GetPositionSession().X,
				store.GetPositionSession().Y))
			movableUnit.SetStatic(false)

		} else {
			movableUnit.SetStatic(true)
		}

		// movableUnit.AddPosition(store.GetPositionSession())
	}

	dispatcher.GetInstance().Dispatch(
		action.NewSyncPreviousPositionSession())

	ss.animator.Update()

	if !ss.transparentTransitionEffect.Done() {
		if !ss.transparentTransitionEffect.OnEnd() {
			ss.transparentTransitionEffect.Update()
		} else {
			ss.transparentTransitionEffect.Clean()
		}
	}

	ss.passiveUI.Update()

	ss.activeUI.Update()

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

func (ss *SessionScreen) HandleRender(screen *ebiten.Image) {
	ss.passiveInterfaceWorld.Clear()

	ss.activeInterfaceWorld.Clear()

	ss.internalWorld.Clear()

	if store.GetEventName() != value.EVENT_NAME_EMPTY_VALUE {
		ss.eventWorld.Clear()
	}

	ss.animator.Draw(ss.internalWorld, ss.camera)

	screen.DrawImage(ss.internalWorld, &ebiten.DrawImageOptions{})

	ss.passiveUI.Draw(ss.passiveInterfaceWorld)

	ss.activeUI.Draw(ss.activeInterfaceWorld)

	screen.DrawImage(ss.passiveInterfaceWorld, &ebiten.DrawImageOptions{
		ColorM: options.GetTransparentDrawOptions(ss.transparentTransitionEffect.GetValue()).ColorM})

	screen.DrawImage(ss.activeInterfaceWorld, &ebiten.DrawImageOptions{
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
	camera := kamera.NewCamera(0, 0, float64(config.GetWorldWidth()), float64(config.GetWorldHeight()))

	camera.ShakeEnabled = true
	camera.SmoothType = kamera.Lerp

	return &SessionScreen{
		passiveUI: builder.Build(
			bar.GetInstance().GetContainer()),
		activeUI: builder.Build(),
		animator: animator.NewAnimator(),
		camera:   camera,
		transparentTransitionEffect: transparent.NewTransparentTransitionEffect(
			true, 255, 0, 5, time.Microsecond*10),
		toxicRainEventStartTransparentTransitionEffect: transparent.NewTransparentTransitionEffect(
			true, 10, 0, 0.5, time.Millisecond*200),
		toxicRainEventEndTransparentTransitionEffect: transparent.NewTransparentTransitionEffect(
			false, 0, 10, 0.5, time.Millisecond*200),
		passiveInterfaceWorld:      ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
		activeInterfaceWorld:       ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
		internalWorld:              ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
		eventWorld:                 ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
		toxicRainEventShaderEffect: toxicrain.NewToxicRainEventEffect(),
	}
}
