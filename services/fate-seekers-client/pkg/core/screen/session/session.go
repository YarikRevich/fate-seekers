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
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/sound"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/animation/direction"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/collision"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/gamepad"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/options"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/renderer"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/renderer/movable"
	selected "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/selecter"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/sounder"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/builder"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/bar"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/press"
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
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/setanarut/kamera/v2"
)

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
	// Represents attached passive transition user interface.
	passiveTransitionUI *ebitenui.UI

	// Represents attached passive user interface.
	passiveUI *ebitenui.UI

	// Represents attached active user interface.
	activeUI *ebitenui.UI

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

	// Represents transparent transition effect used for press component, when there is no selectable anymore.
	pressTransparentTransitionEffect transition.TransitionEffect

	// Represents passive transition interface world view.
	passiveTransitionInterfaceWorld *ebiten.Image

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
	if store.GetResetSession() == value.RESET_SESSION_TRUE_VALUE {
		selected.GetInstance().Clean()

		collision.GetInstance().Clean()

		sounder.GetInstance().Clean()

		renderer.GetInstance().Clean()

		dispatcher.GetInstance().Dispatch(
			action.NewSetResetSession(value.RESET_SESSION_FALSE_VALUE))
	}

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

								ss.camera.AddTrauma(0.8)

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
								action.NewSetChestsRetrievalStartedNetworking(
									value.CHESTS_RETRIEVAL_STARTED_NETWORKING_FALSE_STATE))

							dispatcher.GetInstance().Dispatch(
								action.NewSetHealthPacksRetrievalStartedNetworking(
									value.HEALTH_PACKS_RETRIEVAL_STARTED_NETWORKING_FALSE_STATE))

							// dispatcher.GetInstance().Dispatch(
							// 	action.NewSetUserInventoryRetrievalStartedNetworking(
							// 		value.USER_INVENTORY_RETRIEVAL_STARTED_NETWORKING_FALSE_STATE))

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

							if userMetadata.GetIssuer() == store.GetRepositoryUUID() {
								if previousUsersMetadata.Health != userMetadata.GetHealth() {
									bar.GetInstance().SetHealthText(userMetadata.GetHealth())
								}
							}

							if previousUsersMetadata.Health != userMetadata.Health {
								sharedUsersMetadataHealthHitsIssuers[userMetadata.GetIssuer()] = true
							}
						}

						if userMetadata.GetIssuer() == store.GetRepositoryUUID() {
							if _, ok := retrievedUsersMetadataSession[userMetadata.GetIssuer()]; !ok {
								dispatcher.GetInstance().Dispatch(
									action.NewSetPositionSession(dto.Position{
										X: userMetadata.GetPosition().GetX(),
										Y: userMetadata.GetPosition().GetY(),
									}),
								)

								bar.GetInstance().SetHealthText(userMetadata.GetHealth())

								ss.camera.SetCenter(
									userMetadata.GetPosition().GetX(),
									-userMetadata.GetPosition().GetY())
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

					selected.GetInstance().PruneExternalMovableObjects(sharedUsersMetadataIssuers)

					sounder.GetInstance().PruneExternalTrackableObjects(sharedUsersMetadataIssuers)

					renderer.GetInstance().PruneSecondaryExternalMovableObjects(sharedUsersMetadataIssuers)

					for issuer := range sharedUsersMetadataIssuers {
						retrievedUsersMetadata := store.GetRetrievedUsersMetadataSession()[issuer]

						sounder.GetInstance().SetExternalTrackableObject(issuer, retrievedUsersMetadata.Position)

						var movableUnit *movable.Movable

						if !renderer.GetInstance().SecondaryExternalMovableObjectExists(issuer) {
							movableUnit = movable.NewMovable(
								loader.GetMovableSkinsPath(retrievedUsersMetadata.Skin))

							movableUnit.SetDirection(retrievedUsersMetadata.AnimationDirection)
							movableUnit.SetStatic(retrievedUsersMetadata.AnimationStatic)
							movableUnit.SetPosition(retrievedUsersMetadata.Position)

							renderer.GetInstance().AddSecondaryExternalMovableObject(issuer, movableUnit)
						} else {
							movableUnit = renderer.GetInstance().GetSecondaryExternalMovableObject(issuer)

							movableUnit.SetDirection(retrievedUsersMetadata.AnimationDirection)
							movableUnit.SetStatic(retrievedUsersMetadata.AnimationStatic)
							movableUnit.AddPosition(retrievedUsersMetadata.Position)
						}

						if !selected.GetInstance().ExternalMovableObjectExists(issuer) {
							shiftWidth, shiftHeight := movableUnit.GetShiftBounds()

							selected.GetInstance().AddExternalMovableObject(
								issuer,
								retrievedUsersMetadata.Position,
								shiftWidth,
								shiftHeight)
						} else {
							selectedUnit := selected.GetInstance().GetExternalMovableObject(issuer)

							selectedUnit.SetPosition(
								retrievedUsersMetadata.Position.X, retrievedUsersMetadata.Position.Y)
						}
					}

					var movableUnit *movable.Movable

					for issuer := range sharedUsersMetadataHealthHitsIssuers {
						if renderer.GetInstance().SecondaryExternalMovableObjectExists(issuer) {
							movableUnit = renderer.GetInstance().GetSecondaryExternalMovableObject(issuer)
						} else {
							if renderer.GetInstance().MainCenteredMovableObjectExists(issuer) {
								sound.GetInstance().GetSoundEventsFxManager().PushWithHandbrake(loader.ToxicRainFXSound)

								movableUnit = renderer.GetInstance().GetMainCenteredMovableObject(issuer)
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

	if store.GetUserInventoryRetrievalStartedNetworking() == value.USER_INVENTORY_RETRIEVAL_STARTED_NETWORKING_FALSE_STATE {
		dispatcher.GetInstance().Dispatch(
			action.NewSetUserInventoryRetrievalStartedNetworking(
				value.USER_INVENTORY_RETRIEVAL_STARTED_NETWORKING_TRUE_STATE))

		metadatastream.GetGetUserInventorySubmitter().Clean(func() {
			metadatastream.GetGetUserInventorySubmitter().Submit(
				store.GetSelectedSessionMetadata().ID, func(response *metadatav1.GetUserInventoryResponse, err error) bool {
					fmt.Println(response.GetUserInventoryItems(), err, "INVENTORY ITEMS")

					return false
				})
		})
	}

	if store.GetChestsRetrievalStartedNetworking() == value.CHESTS_RETRIEVAL_STARTED_NETWORKING_FALSE_STATE {
		dispatcher.GetInstance().Dispatch(
			action.NewSetChestsRetrievalStartedNetworking(
				value.CHESTS_RETRIEVAL_STARTED_NETWORKING_TRUE_STATE))

		fmt.Println("STARTING CHESTS RETRIEVAL")

		metadatastream.GetGetChestsSubmitter().Clean(func() {
			metadatastream.GetGetChestsSubmitter().Submit(
				store.GetSelectedSessionMetadata().ID, func(response *metadatav1.GetChestsResponse, err error) bool {
					fmt.Println(response.GetChests(), err, "CHESTS")

					for _, chest := range response.GetChests() {
						for _, chestItem := range chest.GetChestItems() {
							fmt.Println(chestItem.GetName())
						}
					}

					return false
				})
		})
	}

	if store.GetHealthPacksRetrievalStartedNetworking() == value.HEALTH_PACKS_RETRIEVAL_STARTED_NETWORKING_FALSE_STATE {
		dispatcher.GetInstance().Dispatch(
			action.NewSetHealthPacksRetrievalStartedNetworking(
				value.HEALTH_PACKS_RETRIEVAL_STARTED_NETWORKING_TRUE_STATE))

		metadatastream.GetGetHealthPacksSubmitter().Clean(func() {
			metadatastream.GetGetHealthPacksSubmitter().Submit(
				store.GetSelectedSessionMetadata().ID, func(response *metadatav1.GetHealthPacksResponse, err error) bool {
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

	selectedLobbySet := store.GetSelectedLobbySetUnitMetadata()

	retrievedUsersMetadataSession := store.GetRetrievedUsersMetadataSession()

	if _, ok := retrievedUsersMetadataSession[store.GetRepositoryUUID()]; ok {
		if store.GetApplicationStateGamepadEnabled() == value.GAMEPAD_ENABLED_APPLICATION_TRUE_VALUE && ebiten.IsFocused() {
			gamepadID := ebiten.GamepadIDs()[0]

			direction := gamepad.GetGamepadLeftStickDirection(gamepadID)

			switch direction {
			case dto.DirUp:
				dispatcher.GetInstance().Dispatch(action.NewIncrementYPositionSession())

			case dto.DirDown:
				dispatcher.GetInstance().Dispatch(action.NewDecrementYPositionSession())

			case dto.DirLeft:
				dispatcher.GetInstance().Dispatch(action.NewDecrementXPositionSession())

			case dto.DirRight:
				dispatcher.GetInstance().Dispatch(action.NewIncrementXPositionSession())

			case dto.DirUpLeft:
				dispatcher.GetInstance().Dispatch(action.NewDiagonalUpLeftPositionSession())

			case dto.DirUpRight:
				dispatcher.GetInstance().Dispatch(action.NewDiagonalUpRightPositionSession())

			case dto.DirDownLeft:
				dispatcher.GetInstance().Dispatch(action.NewDiagonalDownLeftPositionSession())

			case dto.DirDownRight:
				dispatcher.GetInstance().Dispatch(action.NewDiagonalDownRightPositionSession())

			}
		} else {
			if ebiten.IsKeyPressed(ebiten.KeyEscape) {
				dispatcher.GetInstance().Dispatch(
					action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_RESUME_VALUE))
			}

			if ebiten.IsKeyPressed(ebiten.KeyW) && ebiten.IsKeyPressed(ebiten.KeyA) {
				dispatcher.GetInstance().Dispatch(action.NewDiagonalUpLeftPositionSession())

			} else if ebiten.IsKeyPressed(ebiten.KeyW) && ebiten.IsKeyPressed(ebiten.KeyD) {
				dispatcher.GetInstance().Dispatch(action.NewDiagonalUpRightPositionSession())

			} else if ebiten.IsKeyPressed(ebiten.KeyS) && ebiten.IsKeyPressed(ebiten.KeyA) {
				dispatcher.GetInstance().Dispatch(action.NewDiagonalDownLeftPositionSession())

			} else if ebiten.IsKeyPressed(ebiten.KeyS) && ebiten.IsKeyPressed(ebiten.KeyD) {
				dispatcher.GetInstance().Dispatch(action.NewDiagonalDownRightPositionSession())

			} else if ebiten.IsKeyPressed(ebiten.KeyA) {
				dispatcher.GetInstance().Dispatch(action.NewDecrementXPositionSession())

			} else if ebiten.IsKeyPressed(ebiten.KeyW) {
				dispatcher.GetInstance().Dispatch(action.NewIncrementYPositionSession())

			} else if ebiten.IsKeyPressed(ebiten.KeyS) {
				dispatcher.GetInstance().Dispatch(action.NewDecrementYPositionSession())

			} else if ebiten.IsKeyPressed(ebiten.KeyD) {
				dispatcher.GetInstance().Dispatch(action.NewIncrementXPositionSession())

			}
		}

		if store.GetStagePositionSession().X != store.GetPositionSession().X ||
			store.GetStagePositionSession().Y != store.GetPositionSession().Y {
			if renderer.GetInstance().MainCenteredMovableObjectExists(selectedLobbySet.Issuer) {
				movableUnit := renderer.GetInstance().GetMainCenteredMovableObject(selectedLobbySet.Issuer)

				shiftWidth, shiftHeight := movableUnit.GetShiftBounds()

				if store.GetStagePositionSession().X != store.GetPositionSession().X {
					collision.GetInstance().SetMainTrackableObject(
						dto.Position{
							X: store.GetStagePositionSession().X,
							Y: store.GetPositionSession().Y},
						shiftWidth, shiftHeight)

					if collision.GetInstance().IsColliding() {
						dispatcher.GetInstance().Dispatch(action.NewRevertStagePositionXSession())
					} else {
						dispatcher.GetInstance().Dispatch(action.NewSyncStagePositionXSession())
					}
				}

				if store.GetStagePositionSession().Y != store.GetPositionSession().Y {
					collision.GetInstance().SetMainTrackableObject(
						dto.Position{
							X: store.GetPositionSession().X,
							Y: store.GetStagePositionSession().Y},
						shiftWidth, shiftHeight)

					if collision.GetInstance().IsColliding() {
						dispatcher.GetInstance().Dispatch(action.NewRevertStagePositionYSession())
					} else {
						dispatcher.GetInstance().Dispatch(action.NewSyncStagePositionYSession())
					}
				}

				sounder.GetInstance().SetMainTrackableObject(
					store.GetPositionSession(), shiftWidth, shiftHeight)
			} else {
				dispatcher.GetInstance().Dispatch(action.NewSyncStagePositionXSession())
				dispatcher.GetInstance().Dispatch(action.NewSyncStagePositionYSession())
			}
		} else {
			sounder.GetInstance().InterruptMainTrackableObject()
		}
	}

	store.RetrievedUsersMetadataSessionSyncHelper.Unlock()

	var movableUnit *movable.Movable

	if !renderer.GetInstance().MainCenteredMovableObjectExists(selectedLobbySet.Issuer) {
		movableUnit = movable.NewMovable(
			loader.GetMovableSkinsPath(selectedLobbySet.Skin))

		movableUnit.SetDirection(dto.RightMovableRotation)
		movableUnit.SetStatic(true)

		renderer.GetInstance().AddMainCenteredMovableObject(selectedLobbySet.Issuer, movableUnit)
	} else {
		movableUnit = renderer.GetInstance().GetMainCenteredMovableObject(selectedLobbySet.Issuer)

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
	}

	movableUnit.SetPosition(store.GetPositionSession())

	shiftWidth, shiftHeight := movableUnit.GetShiftBounds()

	ss.camera.LookAt(
		store.GetPositionSession().X+(shiftWidth/2),
		-store.GetPositionSession().Y+(shiftHeight/2))

	dispatcher.GetInstance().Dispatch(
		action.NewSyncPreviousPositionSession())

	selectedPosition, ok := selected.GetInstance().Scan(ss.camera)
	if ok {
		ss.pressTransparentTransitionEffect.Reset()

		renderer.GetInstance().SetSelectedObject(selectedPosition)

		var pressed bool

		if store.GetApplicationStateGamepadEnabled() == value.GAMEPAD_ENABLED_APPLICATION_TRUE_VALUE && ebiten.IsFocused() {
			press.GetInstance().SetPressType(press.GAMEPAD)

			if inpututil.IsStandardGamepadButtonJustPressed(ebiten.GamepadIDs()[0], ebiten.StandardGamepadButtonRightRight) {
				pressed = true
			}
		} else {
			press.GetInstance().SetPressType(press.KEYBOARD)

			if ebiten.IsKeyPressed(ebiten.KeyE) {
				pressed = true
			}
		}

		if pressed {
			// TODO: check for saved location associations in memory.
		}
	} else {
		if !ss.pressTransparentTransitionEffect.Done() {
			if !ss.pressTransparentTransitionEffect.OnEnd() {
				ss.pressTransparentTransitionEffect.Update()
			} else {
				ss.pressTransparentTransitionEffect.Clean()
			}
		}

		renderer.GetInstance().DisableSelectedObject()
	}

	renderer.GetInstance().Update(ss.camera)

	sounder.GetInstance().Update()

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

	// fmt.Println(
	// 	math.Round(ss.camera.X*math.Pow(10, 2))/math.Pow(10, 2),
	// 	math.Round(ss.camera.Y*math.Pow(10, 2))/math.Pow(10, 2))

	// TODO: click on the letter.

	// dispatcher.GetInstance().Dispatch(action.NewSetLetterImageAction(""))

	// TODO: click on the chest.
	// dispatcher.GetInstance().Dispatch(action.New)

	return nil
}

func (ss *SessionScreen) HandleRender(screen *ebiten.Image) {
	// TODO: refactor to remove session sync helper lock usage.
	store.RetrievedUsersMetadataSessionSyncHelper.Lock()

	retrievedUsersMetadataSession := store.GetRetrievedUsersMetadataSession()

	if _, ok := retrievedUsersMetadataSession[store.GetRepositoryUUID()]; ok {
		ss.passiveTransitionInterfaceWorld.Clear()

		ss.passiveInterfaceWorld.Clear()

		ss.activeInterfaceWorld.Clear()

		ss.internalWorld.Clear()

		if store.GetEventName() != value.EVENT_NAME_EMPTY_VALUE {
			ss.eventWorld.Clear()
		}

		renderer.GetInstance().Draw(ss.internalWorld, ss.camera)

		screen.DrawImage(ss.internalWorld, &ebiten.DrawImageOptions{})

		ss.passiveTransitionUI.Draw(ss.passiveTransitionInterfaceWorld)

		ss.passiveUI.Draw(ss.passiveInterfaceWorld)

		ss.activeUI.Draw(ss.activeInterfaceWorld)

		if !ss.pressTransparentTransitionEffect.Done() {
			screen.DrawImage(ss.passiveTransitionInterfaceWorld, &ebiten.DrawImageOptions{
				ColorM: options.GetTransparentDrawOptions(ss.pressTransparentTransitionEffect.GetValue()).ColorM})
		}

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

	store.RetrievedUsersMetadataSessionSyncHelper.Unlock()
}

// newSessionScreen initializes SessionScreen.
func newSessionScreen() screen.Screen {
	camera := kamera.NewCamera(0, 0, float64(config.GetWorldWidth()), float64(config.GetWorldHeight()))

	camera.ShakeEnabled = true
	camera.SmoothType = kamera.Lerp

	pressTransparentTransitionEffect := transparent.NewTransparentTransitionEffect(false, 0, 255, 5, time.Microsecond*10)
	pressTransparentTransitionEffect.Clean()

	return &SessionScreen{
		passiveTransitionUI: builder.Build(
			press.GetInstance().GetContainer()),
		passiveUI: builder.Build(
			bar.GetInstance().GetContainer()),
		activeUI: builder.Build(),
		camera:   camera,
		transparentTransitionEffect: transparent.NewTransparentTransitionEffect(
			true, 255, 0, 5, time.Microsecond*10),
		toxicRainEventStartTransparentTransitionEffect: transparent.NewTransparentTransitionEffect(
			true, 10, 0, 0.5, time.Millisecond*200),
		toxicRainEventEndTransparentTransitionEffect: transparent.NewTransparentTransitionEffect(
			false, 0, 10, 0.5, time.Millisecond*200),
		pressTransparentTransitionEffect: pressTransparentTransitionEffect,
		passiveTransitionInterfaceWorld:  ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
		passiveInterfaceWorld:            ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
		activeInterfaceWorld:             ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
		internalWorld:                    ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
		eventWorld:                       ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
		toxicRainEventShaderEffect:       toxicrain.NewToxicRainEventEffect(),
	}
}
