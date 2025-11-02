package lobby

import (
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition/transparent"
	metadatav1 "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/converter"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/handler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/stream"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/options"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/builder"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/lobby"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/notification"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/translation"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/storage/shared"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// GetInstance retrieves instance of the lobby screen, performing initilization if needed.
	GetInstance = sync.OnceValue[screen.Screen](newLobbyScreen)
)

// LobbyScreen represents lobby screen implementation.
type LobbyScreen struct {
	// Represents attached user interface.
	ui *ebitenui.UI

	// Represents transparent transition effect.
	transparentTransitionEffect transition.TransitionEffect

	// Represents global world view.
	world *ebiten.Image

	// Represents internal state reset flag.
	stateReset bool
}

func (ls *LobbyScreen) HandleInput() error {
	if store.GetLobbySetRetrievalStartedNetworking() == value.LOBBY_SET_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE {
		dispatcher.GetInstance().Dispatch(
			action.NewSetLobbySetRetrievalStartedNetworkingAction(value.LOBBY_SET_RETRIEVAL_STARTED_NETWORKING_TRUE_VALUE))

		stream.GetGetLobbySetSubmitter().Clean(func() {
			stream.GetGetLobbySetSubmitter().Submit(
				store.GetSelectedSessionMetadata().ID, func(response *metadatav1.GetLobbySetResponse, err error) bool {
					if store.GetActiveScreen() != value.ACTIVE_SCREEN_LOBBY_VALUE {
						dispatcher.GetInstance().Dispatch(
							action.NewSetLobbySetRetrievalStartedNetworkingAction(
								value.LOBBY_SET_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE))

						return true
					}

					if response == nil || err != nil {
						notification.GetInstance().Push(
							common.ComposeMessage(
								translation.GetInstance().GetTranslation("client.networking.get-lobby-set-failure"),
								err.Error()),
							time.Second*3,
							common.NotificationErrorTextColor)

						dispatcher.
							GetInstance().
							Dispatch(
								action.NewSetStateResetApplicationAction(
									value.STATE_RESET_APPLICATION_FALSE_VALUE))

						dispatcher.
							GetInstance().
							Dispatch(
								action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))

						dispatcher.GetInstance().Dispatch(
							action.NewSetLobbySetRetrievalStartedNetworkingAction(
								value.LOBBY_SET_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE))

						return true
					}

					var isLobbySetUpdated bool

					if len(store.GetRetrievedLobbySetMetadata()) != len(response.GetLobbySet()) {
						isLobbySetUpdated = true
					} else {
						for index, value := range store.GetRetrievedLobbySetMetadata() {
							if value.Issuer != response.GetLobbySet()[index].GetIssuer() {
								isLobbySetUpdated = true
							}
						}
					}

					if isLobbySetUpdated || ls.stateReset {
						ls.stateReset = false

						for _, value := range response.GetLobbySet() {
							if value.GetIssuer() == store.GetRepositoryUUID() {
								dispatcher.
									GetInstance().
									Dispatch(
										action.NewSetSelectedLobbySetUnitMetadata(
											&dto.SelectedLobbySetUnitMetadata{
												ID:     value.GetLobbyId(),
												Issuer: value.GetIssuer(),
												Skin:   value.GetSkin(),
												Host:   value.GetHost(),
											}))
							}
						}

						dispatcher.
							GetInstance().
							Dispatch(
								action.NewSetRetrievedLobbySetMetadata(
									converter.ConvertGetLobbySetResponseToRetrievedLobbySetMetadata(
										response)))

						var (
							selectedPlayer *metadatav1.LobbySetUnit
							otherPlayers   []*metadatav1.LobbySetUnit
						)

						for _, lobbySetUnit := range response.GetLobbySet() {
							if lobbySetUnit.Issuer == store.GetRepositoryUUID() {
								selectedPlayer = lobbySetUnit
							} else {
								otherPlayers = append(otherPlayers, lobbySetUnit)
							}
						}

						lobby.GetInstance().SetSelectionBySkin(
							selectedPlayer.GetSkin())

						lobby.GetInstance().SetListsEntries(
							converter.ConvertGetLobbySetResponseToListEntries(otherPlayers))

						if store.GetSessionAlreadyStartedMetadata() == value.SESSION_ALREADY_STARTED_METADATA_STATE_FALSE_VALUE {
							for _, value := range response.GetLobbySet() {
								if value.GetIssuer() == store.GetRepositoryUUID() && value.GetHost() {
									lobby.GetInstance().ShowStartButton()
								}
							}
						}

						if store.GetLobbySetRetrievalCycleFinishedNetworking() == value.LOBBY_SET_RETRIEVAL_CYCLE_FINISHED_NETWORKING_FALSE_VALUE {
							dispatcher.GetInstance().Dispatch(
								action.NewSetLobbySetRetrievalCycleFinishedNetworkingAction(
									value.LOBBY_SET_RETRIEVAL_CYCLE_FINISHED_NETWORKING_TRUE_VALUE))
						}
					}

					return false
				})
		})
	}

	if (store.GetSessionMetadataRetrievalStartedNetworking() == value.SESSION_METADATA_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE) &&
		(store.GetLobbySetRetrievalCycleFinishedNetworking() == value.LOBBY_SET_RETRIEVAL_CYCLE_FINISHED_NETWORKING_TRUE_VALUE) {
		dispatcher.GetInstance().Dispatch(
			action.NewSetSessionMetadataRetrievalStartedNetworkingAction(
				value.SESSION_METADATA_RETRIEVAL_STARTED_NETWORKING_TRUE_VALUE))

		stream.GetGetSessionMetadataSubmitter().Clean(func() {
			stream.GetGetSessionMetadataSubmitter().Submit(
				store.GetSelectedSessionMetadata().ID, func(response *metadatav1.GetSessionMetadataResponse, err error) bool {
					if store.GetActiveScreen() != value.ACTIVE_SCREEN_LOBBY_VALUE {
						return true
					}

					if response == nil || err != nil {
						notification.GetInstance().Push(
							common.ComposeMessage(
								translation.GetInstance().GetTranslation("client.networking.get-session-metadata-failure"),
								err.Error()),
							time.Second*3,
							common.NotificationErrorTextColor)

						dispatcher.
							GetInstance().
							Dispatch(
								action.NewSetStateResetApplicationAction(
									value.STATE_RESET_APPLICATION_FALSE_VALUE))

						dispatcher.
							GetInstance().
							Dispatch(
								action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))

						return true
					}

					if response.GetStarted() {
						notification.GetInstance().Push(
							translation.GetInstance().GetTranslation("client.lobby.transfering-to-session"),
							time.Second*4,
							common.NotificationInfoTextColor)

						if store.GetSessionAlreadyStartedMetadata() == value.SESSION_ALREADY_STARTED_METADATA_STATE_TRUE_VALUE {
							dispatcher.GetInstance().Dispatch(
								action.NewSetSessionAlreadyStartedMetadata(
									value.SESSION_ALREADY_STARTED_METADATA_STATE_FALSE_VALUE))
						}

						if store.GetLobbySetRetrievalCycleFinishedNetworking() == value.LOBBY_SET_RETRIEVAL_CYCLE_FINISHED_NETWORKING_TRUE_VALUE {
							dispatcher.GetInstance().Dispatch(
								action.NewSetLobbySetRetrievalCycleFinishedNetworkingAction(
									value.LOBBY_SET_RETRIEVAL_CYCLE_FINISHED_NETWORKING_FALSE_VALUE))
						}

						dispatcher.
							GetInstance().
							Dispatch(
								action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_TRAVEL_VALUE))

						return true
					}

					return false
				})
		})
	}

	if !ls.transparentTransitionEffect.Done() {
		if !ls.transparentTransitionEffect.OnEnd() {
			ls.transparentTransitionEffect.Update()
		} else {
			ls.transparentTransitionEffect.Clean()
		}
	}

	shared.GetInstance().GetBackgroundAnimation().Update()

	ls.ui.Update()

	return nil
}

func (ls *LobbyScreen) HandleRender(screen *ebiten.Image) {
	var backgroundAnimationGeometry ebiten.GeoM

	backgroundAnimationGeometry.Scale(
		scaler.GetScaleFactor(config.GetMinStaticWidth(), config.GetWorldWidth()),
		scaler.GetScaleFactor(config.GetMinStaticHeight(), config.GetWorldHeight()))

	shared.GetInstance().GetBackgroundAnimation().DrawTo(ls.world, &ebiten.DrawImageOptions{
		GeoM: backgroundAnimationGeometry,
	})

	ls.ui.Draw(ls.world)

	screen.DrawImage(ls.world, &ebiten.DrawImageOptions{
		ColorM: options.GetTransparentDrawOptions(
			ls.transparentTransitionEffect.GetValue()).ColorM})
}

// newLobbyScreen initializes LobbyScreen.
func newLobbyScreen() screen.Screen {
	transparentTransitionEffect := transparent.NewTransparentTransitionEffect(true, 255, 0, 5, time.Microsecond*10)

	lobby.GetInstance().SetStartCallback(func() {
		handler.PerformStartSession(
			store.GetSelectedSessionMetadata().ID,
			store.GetSelectedLobbySetUnitMetadata().ID,
			func(err error) {
				if err != nil {
					notification.GetInstance().Push(
						common.ComposeMessage(
							translation.GetInstance().GetTranslation("client.networking.start-session-failure"),
							err.Error()),
						time.Second*3,
						common.NotificationErrorTextColor)

					dispatcher.
						GetInstance().
						Dispatch(
							action.NewSetStateResetApplicationAction(
								value.STATE_RESET_APPLICATION_FALSE_VALUE))

					dispatcher.
						GetInstance().
						Dispatch(
							action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))

					return
				}

				notification.GetInstance().Push(
					translation.GetInstance().GetTranslation("client.networking.start-session-processing"),
					time.Second*2,
					common.NotificationInfoTextColor)

				notification.GetInstance().Push(
					translation.GetInstance().GetTranslation("client.lobby.transfering-to-session"),
					time.Second*4,
					common.NotificationInfoTextColor)

				dispatcher.
					GetInstance().
					Dispatch(
						action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_TRAVEL_VALUE))
			})
	})

	instance := &LobbyScreen{
		ui: builder.Build(
			lobby.GetInstance().GetContainer()),
		transparentTransitionEffect: transparentTransitionEffect,
		world:                       ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
	}

	lobby.GetInstance().SetBackCallback(func() {
		handler.PerformRemoveLobby(store.GetSelectedSessionMetadata().ID, func(err error) {
			transparentTransitionEffect.Reset()

			instance.stateReset = true

			lobby.GetInstance().HideStartButton()

			if err != nil {
				notification.GetInstance().Push(
					common.ComposeMessage(
						translation.GetInstance().GetTranslation("client.networking.remove-lobby-failure"),
						err.Error()),
					time.Second*3,
					common.NotificationErrorTextColor)
			}

			dispatcher.GetInstance().Dispatch(
				action.NewSetLobbySetRetrievalStartedNetworkingAction(value.LOBBY_SET_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE))

			dispatcher.GetInstance().Dispatch(
				action.NewSetSessionMetadataRetrievalStartedNetworkingAction(
					value.SESSION_METADATA_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE))

			dispatcher.GetInstance().Dispatch(
				action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_SELECTOR_VALUE))
		})
	})

	lobby.GetInstance().HideStartButton()

	return instance
}
