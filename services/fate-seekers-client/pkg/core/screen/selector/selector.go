package selector

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition/transparent"
	metadatav1 "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/converter"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/handler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/options"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/builder"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/selector"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/notification"
	selectormanager "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/selector"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/translation"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/storage/shared"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/exp/slices"
)

var (
	// GetInstance retrieves instance of the selector screen, performing initilization if needed.
	GetInstance = sync.OnceValue[screen.Screen](newSelectorScreen)
)

// SelectorScreen represents selector screen implementation.
type SelectorScreen struct {
	// Represents attached user interface.
	ui *ebitenui.UI

	// Represents transparent transition effect.
	transparentTransitionEffect transition.TransitionEffect

	// Represents global world view.
	world *ebiten.Image

	// Represents interface world view.
	interfaceWorld *ebiten.Image
}

func (ss *SelectorScreen) HandleInput() error {
	if store.GetSessionRetrievalStartedNetworking() == value.SESSION_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE {
		dispatcher.GetInstance().Dispatch(
			action.NewSetSessionRetrievalStartedNetworkingAction(value.SESSION_RETRIEVAL_STARTED_NETWORKING_TRUE_VALUE))

		handler.PerformGetUserSessions(func(response *metadatav1.GetUserSessionsResponse, err error) {
			if err != nil {
				notification.GetInstance().Push(
					common.ComposeMessage(
						translation.GetInstance().GetTranslation("client.networking.get-user-sessions-failure"),
						err.Error()),
					time.Second*3,
					common.NotificationErrorTextColor)

				return
			}

			dispatcher.
				GetInstance().
				Dispatch(
					action.NewSetRetrievedSessionsMetadata(
						converter.ConvertGetUserSessionsResponseToRetrievedSessionsMetadata(
							response)))

			selector.GetInstance().SetListsEntries(
				converter.ConvertGetUserSessionsResponseToListEntries(response))
		})
	}

	if !ss.transparentTransitionEffect.Done() {
		if !ss.transparentTransitionEffect.OnEnd() {
			ss.transparentTransitionEffect.Update()
		} else {
			ss.transparentTransitionEffect.Clean()
		}
	}

	shared.GetInstance().GetBackgroundAnimation().Update()

	ss.ui.Update()

	return nil
}

func (ss *SelectorScreen) HandleRender(screen *ebiten.Image) {
	ss.world.Clear()

	ss.interfaceWorld.Clear()

	var backgroundAnimationGeometry ebiten.GeoM

	backgroundAnimationGeometry.Scale(
		scaler.GetScaleFactor(config.GetMinStaticWidth(), config.GetWorldWidth()),
		scaler.GetScaleFactor(config.GetMinStaticHeight(), config.GetWorldHeight()))

	shared.GetInstance().GetBackgroundAnimation().DrawTo(ss.world, &ebiten.DrawImageOptions{
		GeoM: backgroundAnimationGeometry,
	})

	ss.ui.Draw(ss.interfaceWorld)

	ss.world.DrawImage(ss.interfaceWorld, &ebiten.DrawImageOptions{
		ColorM: options.GetTransparentDrawOptions(
			ss.transparentTransitionEffect.GetValue()).ColorM})

	screen.DrawImage(ss.world, &ebiten.DrawImageOptions{})
}

// newSelectorScreen initializes SelectorScreen.
func newSelectorScreen() screen.Screen {
	transparentTransitionEffect := transparent.NewTransparentTransitionEffect(true, 255, 0, 5, time.Microsecond*10)

	selector.GetInstance().SetSubmitCallback(func(sessionName string) {
		if selectormanager.ProcessChanges(sessionName) {
			if store.GetLobbyCreationStartedNetworking() == value.LOBBY_CREATION_STARTED_NETWORKING_FALSE_VALUE {
				dispatcher.GetInstance().Dispatch(
					action.NewSetLobbyCreationStartedNetworkingAction(
						value.LOBBY_CREATION_STARTED_NETWORKING_TRUE_VALUE))

				if slices.ContainsFunc(
					store.GetRetrievedSessionsMetadata(),
					func(value dto.RetrievedSessionMetadata) bool {
						return value.Name == sessionName
					}) {
					var sessionID int64

					for _, session := range store.GetRetrievedSessionsMetadata() {
						if session.Name == sessionName {
							sessionID = session.SessionID

							break
						}
					}

					dispatcher.GetInstance().Dispatch(
						action.NewSetSelectedSessionMetadata(&dto.SelectedSessionMetadata{
							ID:   sessionID,
							Name: sessionName,
						}))

					handler.PerformCreateLobby(sessionID, func(err error) {
						fmt.Println(err)

						if errors.Is(err, handler.ErrLobbyAlreadyExists) {
							notification.GetInstance().Push(
								translation.GetInstance().GetTranslation("client.networking.joining-existing-lobby"),
								time.Second*3,
								common.NotificationInfoTextColor)
						} else if err != nil {
							notification.GetInstance().Push(
								common.ComposeMessage(
									translation.GetInstance().GetTranslation("client.networking.create-lobby-failure"),
									err.Error()),
								time.Second*3,
								common.NotificationErrorTextColor)

							dispatcher.GetInstance().Dispatch(
								action.NewSetLobbyCreationStartedNetworkingAction(
									value.LOBBY_CREATION_STARTED_NETWORKING_FALSE_VALUE))

							return
						}

						transparentTransitionEffect.Reset()

						dispatcher.GetInstance().Dispatch(
							action.NewSetSessionRetrievalStartedNetworkingAction(value.SESSION_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE))

						fmt.Println("BEFORE 10")

						dispatcher.GetInstance().Dispatch(
							action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_LOBBY_VALUE))

						dispatcher.GetInstance().Dispatch(
							action.NewSetLobbyCreationStartedNetworkingAction(
								value.LOBBY_CREATION_STARTED_NETWORKING_FALSE_VALUE))

						selector.GetInstance().CleanInputs()

						selector.GetInstance().ResetDeleteButton()
					})
				} else {
					handler.PerformGetUserSessions(func(response *metadatav1.GetUserSessionsResponse, err error) {
						if err != nil {
							notification.GetInstance().Push(
								common.ComposeMessage(
									translation.GetInstance().GetTranslation("client.networking.get-user-sessions-failure"),
									err.Error()),
								time.Second*3,
								common.NotificationErrorTextColor)

							dispatcher.GetInstance().Dispatch(
								action.NewSetLobbyCreationStartedNetworkingAction(
									value.LOBBY_CREATION_STARTED_NETWORKING_FALSE_VALUE))

							return
						}

						convertedGetSessionsResponse :=
							converter.ConvertGetUserSessionsResponseToRetrievedSessionsMetadata(response)

						dispatcher.
							GetInstance().
							Dispatch(
								action.NewSetRetrievedSessionsMetadata(
									convertedGetSessionsResponse))

						selector.GetInstance().SetListsEntries(
							converter.ConvertGetUserSessionsResponseToListEntries(response))

						var (
							found     bool
							sessionID int64
						)

						for _, session := range convertedGetSessionsResponse {
							if session.Name == sessionName {
								found = true
								sessionID = session.SessionID

								break
							}
						}

						if found {
							dispatcher.GetInstance().Dispatch(
								action.NewSetSelectedSessionMetadata(&dto.SelectedSessionMetadata{
									ID:   sessionID,
									Name: sessionName,
								}))

							handler.PerformCreateLobby(sessionID, func(err error) {
								if errors.Is(err, handler.ErrLobbyAlreadyExists) {
									notification.GetInstance().Push(
										translation.GetInstance().GetTranslation("client.networking.joining-existing-lobby"),
										time.Second*3,
										common.NotificationInfoTextColor)
								} else if err != nil {
									notification.GetInstance().Push(
										common.ComposeMessage(
											translation.GetInstance().GetTranslation("client.networking.create-lobby-failure"),
											err.Error()),
										time.Second*3,
										common.NotificationErrorTextColor)

									dispatcher.GetInstance().Dispatch(
										action.NewSetLobbyCreationStartedNetworkingAction(
											value.LOBBY_CREATION_STARTED_NETWORKING_FALSE_VALUE))

									return
								}

								transparentTransitionEffect.Reset()

								dispatcher.GetInstance().Dispatch(
									action.NewSetSessionRetrievalStartedNetworkingAction(value.SESSION_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE))

								fmt.Println("BEFORE 30")

								dispatcher.GetInstance().Dispatch(
									action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_LOBBY_VALUE))

								dispatcher.GetInstance().Dispatch(
									action.NewSetLobbyCreationStartedNetworkingAction(
										value.LOBBY_CREATION_STARTED_NETWORKING_FALSE_VALUE))

								selector.GetInstance().CleanInputs()

								selector.GetInstance().ResetDeleteButton()
							})
						} else {
							handler.PerformGetFilteredSessions(dto.GetFilteredSessionsRequest{
								Name: sessionName,
							}, func(response *metadatav1.GetFilteredSessionResponse, err error) {
								if errors.Is(err, handler.ErrFilteredSessionDoesNotExist) {
									notification.GetInstance().Push(
										translation.GetInstance().GetTranslation("client.networking.filtered-session-not-found"),
										time.Second*3,
										common.NotificationErrorTextColor)

									dispatcher.GetInstance().Dispatch(
										action.NewSetLobbyCreationStartedNetworkingAction(
											value.LOBBY_CREATION_STARTED_NETWORKING_FALSE_VALUE))

									return
								} else if err != nil {
									notification.GetInstance().Push(
										common.ComposeMessage(
											translation.GetInstance().GetTranslation("client.networking.get-filtered-sessions-failure"),
											err.Error()),
										time.Second*3,
										common.NotificationErrorTextColor)

									dispatcher.GetInstance().Dispatch(
										action.NewSetLobbyCreationStartedNetworkingAction(
											value.LOBBY_CREATION_STARTED_NETWORKING_FALSE_VALUE))

									return
								}

								dispatcher.GetInstance().Dispatch(
									action.NewSetSelectedSessionMetadata(&dto.SelectedSessionMetadata{
										ID:   response.Session.GetSessionId(),
										Name: response.Session.GetName(),
									}))

								handler.PerformCreateLobby(response.GetSession().GetSessionId(), func(err error) {
									if errors.Is(err, handler.ErrLobbyAlreadyExists) {
										notification.GetInstance().Push(
											translation.GetInstance().GetTranslation("client.networking.joining-existing-lobby"),
											time.Second*3,
											common.NotificationInfoTextColor)
									} else if err != nil {
										notification.GetInstance().Push(
											common.ComposeMessage(
												translation.GetInstance().GetTranslation("client.networking.create-lobby-failure"),
												err.Error()),
											time.Second*3,
											common.NotificationErrorTextColor)

										dispatcher.GetInstance().Dispatch(
											action.NewSetLobbyCreationStartedNetworkingAction(
												value.LOBBY_CREATION_STARTED_NETWORKING_FALSE_VALUE))

										return
									}

									transparentTransitionEffect.Reset()

									dispatcher.GetInstance().Dispatch(
										action.NewSetSessionRetrievalStartedNetworkingAction(value.SESSION_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE))

									fmt.Println("BEFORE 20")

									dispatcher.GetInstance().Dispatch(
										action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_LOBBY_VALUE))

									dispatcher.GetInstance().Dispatch(
										action.NewSetLobbyCreationStartedNetworkingAction(
											value.LOBBY_CREATION_STARTED_NETWORKING_FALSE_VALUE))

									selector.GetInstance().CleanInputs()

									selector.GetInstance().ResetDeleteButton()
								})
							})
						}
					})
				}
			}
		}
	})

	selector.GetInstance().SetCreateCallback(func() {
		transparentTransitionEffect.Reset()

		selector.GetInstance().CleanInputs()

		selector.GetInstance().ResetDeleteButton()

		dispatcher.GetInstance().Dispatch(
			action.NewSetSessionRetrievalStartedNetworkingAction(value.SESSION_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE))

		dispatcher.GetInstance().Dispatch(
			action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_CREATOR_VALUE))
	})

	selector.GetInstance().SetDeleteCallback(func(sessionName string) {
		if store.GetSessionRemovalStartedNetworking() == value.SESSION_REMOVAL_STARTED_NETWORKING_FALSE_VALUE {
			dispatcher.GetInstance().Dispatch(
				action.NewSetLobbyCreationStartedNetworkingAction(
					value.SESSION_REMOVAL_STARTED_NETWORKING_TRUE_VALUE))

			if slices.ContainsFunc(
				store.GetRetrievedSessionsMetadata(),
				func(value dto.RetrievedSessionMetadata) bool {
					return value.Name == sessionName
				}) {
				var sessionID int64

				for _, session := range store.GetRetrievedSessionsMetadata() {
					if session.Name == sessionName {
						sessionID = session.SessionID

						break
					}
				}

				handler.PerformRemoveSession(sessionID, func(err error) {
					if err != nil {
						notification.GetInstance().Push(
							common.ComposeMessage(
								translation.GetInstance().GetTranslation("client.networking.remove-session-failure"),
								err.Error()),
							time.Second*3,
							common.NotificationErrorTextColor)

						return
					}

					selector.GetInstance().CleanInputs()

					selector.GetInstance().ResetDeleteButton()

					notification.GetInstance().Push(
						translation.GetInstance().GetTranslation("client.selectormanager.session-removal-in-progress"),
						time.Second*4,
						common.NotificationInfoTextColor)

					dispatcher.GetInstance().Dispatch(
						action.NewSetSessionRetrievalStartedNetworkingAction(value.SESSION_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE))
				})
			} else {
				handler.PerformGetUserSessions(func(response *metadatav1.GetUserSessionsResponse, err1 error) {
					if err1 != nil {
						notification.GetInstance().Push(
							common.ComposeMessage(
								translation.GetInstance().GetTranslation("client.networking.get-user-sessions-failure"),
								err1.Error()),
							time.Second*3,
							common.NotificationErrorTextColor)

						return
					}

					convertedGetSessionsResponse :=
						converter.ConvertGetUserSessionsResponseToRetrievedSessionsMetadata(response)

					dispatcher.
						GetInstance().
						Dispatch(
							action.NewSetRetrievedSessionsMetadata(
								convertedGetSessionsResponse))

					selector.GetInstance().SetListsEntries(
						converter.ConvertGetUserSessionsResponseToListEntries(response))

					var sessionID int64

					for _, session := range convertedGetSessionsResponse {
						if session.Name == sessionName {
							sessionID = session.SessionID

							break
						}
					}

					handler.PerformRemoveSession(sessionID, func(err2 error) {
						if err2 != nil {
							notification.GetInstance().Push(
								common.ComposeMessage(
									translation.GetInstance().GetTranslation("client.networking.remove-session-failure"),
									err2.Error()),
								time.Second*3,
								common.NotificationErrorTextColor)

							return
						}

						selector.GetInstance().CleanInputs()

						selector.GetInstance().ResetDeleteButton()

						notification.GetInstance().Push(
							translation.GetInstance().GetTranslation("client.selectormanager.session-removal-in-progress"),
							time.Second*4,
							common.NotificationInfoTextColor)

						dispatcher.GetInstance().Dispatch(
							action.NewSetSessionRetrievalStartedNetworkingAction(value.SESSION_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE))
					})
				})
			}
		}
	})

	selector.GetInstance().SetBackCallback(func() {
		transparentTransitionEffect.Reset()

		selector.GetInstance().CleanInputs()

		selector.GetInstance().ResetDeleteButton()

		dispatcher.GetInstance().Dispatch(
			action.NewSetSessionRetrievalStartedNetworkingAction(value.SESSION_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE))

		dispatcher.GetInstance().Dispatch(
			action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))
	})

	return &SelectorScreen{
		ui:                          builder.Build(selector.GetInstance().GetContainer()),
		transparentTransitionEffect: transparentTransitionEffect,
		world:                       ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
		interfaceWorld: ebiten.NewImage(
			config.GetWorldWidth(), config.GetWorldHeight()),
	}
}
