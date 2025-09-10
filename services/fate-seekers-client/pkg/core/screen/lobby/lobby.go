package lobby

import (
	"fmt"
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
}

func (ls *LobbyScreen) HandleInput() error {
	if store.GetLobbySetRetrievalStartedNetworking() == value.LOBBY_SET_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE {
		dispatcher.GetInstance().Dispatch(
			action.NewSetLobbySetRetrievalStartedNetworkingAction(value.LOBBY_SET_RETRIEVAL_STARTED_NETWORKING_TRUE_VALUE))

		var sessionID int64

		for _, session := range store.GetRetrievedSessionsMetadata() {
			if session.Name == store.GetSelectedSessionMetadata() {
				sessionID = session.SessionID

				break
			}
		}

		handler.PerformGetLobbySet(sessionID, func(response *metadatav1.GetLobbySetResponse, err error) {
			if err != nil {
				notification.GetInstance().Push(
					common.ComposeMessage(
						translation.GetInstance().GetTranslation("client.networking.get-lobby-set-failure"),
						err.Error()),
					time.Second*3,
					common.NotificationErrorTextColor)

				return
			}

			for _, value := range response.GetLobbySet() {
				fmt.Println(value.GetIssuer(), store.GetRepositoryUUID(), value.GetHost())

				if value.GetIssuer() == store.GetRepositoryUUID() && value.GetHost() {
					lobby.GetInstance().ShowStartButton()
				}
			}

			dispatcher.
				GetInstance().
				Dispatch(
					action.NewSetRetrievedLobbySetMetadata(
						converter.ConvertGetLobbySetResponseToRetrievedLobbySetMetadata(
							response)))

			lobby.GetInstance().SetListsEntries(
				converter.ConvertGetLobbySetResponseToListEntries(response))
		})
	}

	if store.GetSessionMetadataRetrievalStartedNetworking() == value.SESSION_METADATA_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE {
		dispatcher.GetInstance().Dispatch(
			action.NewSetSessionMetadataRetrievalStartedNetworkingAction(
				value.SESSION_METADATA_RETRIEVAL_STARTED_NETWORKING_TRUE_VALUE))

		var sessionID int64

		for _, session := range store.GetRetrievedSessionsMetadata() {
			if session.Name == store.GetSelectedSessionMetadata() {
				sessionID = session.SessionID

				break
			}
		}

		stream.GetGetSessionMetadataSubmitter().Clean(func() {
			stream.GetGetSessionMetadataSubmitter().Submit(
				sessionID, func(response *metadatav1.GetSessionMetadataResponse, err error) bool {
					if store.GetActiveScreen() != value.ACTIVE_SCREEN_LOBBY_VALUE {
						return true
					}

					if response.GetStarted() &&
						store.GetLobbySetRetrievalStartedNetworking() ==
							value.LOBBY_SET_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE {
						notification.GetInstance().Push(
							translation.GetInstance().GetTranslation("client.lobby.transfering-to-session"),
							time.Second*4,
							common.NotificationInfoTextColor)

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

	})

	lobby.GetInstance().SetBackCallback(func() {
		transparentTransitionEffect.Reset()

		lobby.GetInstance().HideStartButton()

		dispatcher.GetInstance().Dispatch(
			action.NewSetLobbySetRetrievalStartedNetworkingAction(value.LOBBY_SET_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE))

		dispatcher.GetInstance().Dispatch(
			action.NewSetSessionMetadataRetrievalStartedNetworkingAction(
				value.SESSION_METADATA_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE))

		dispatcher.GetInstance().Dispatch(
			action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_SELECTOR_VALUE))
	})

	lobby.GetInstance().HideStartButton()

	return &LobbyScreen{
		ui: builder.Build(
			lobby.GetInstance().GetContainer()),
		transparentTransitionEffect: transparentTransitionEffect,
		world:                       ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
	}
}
