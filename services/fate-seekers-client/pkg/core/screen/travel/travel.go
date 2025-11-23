package travel

import (
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/particle"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/particle/loadingstars"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/converter"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/handler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/builder"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/lobby"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/notification"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/translation"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader/utils"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// GetInstance retrieves instance of the travel screen, performing initilization if needed.
	GetInstance = sync.OnceValue[screen.Screen](newTravelScreen)
)

// TravelScreen represents travel screen implementation.
type TravelScreen struct {
	// Represents attached user interface.
	ui *ebitenui.UI

	// Represents global world view.
	world *ebiten.Image

	// Represents session loading stars particle effect.
	loadingStarsParticleEffect particle.ParticleEffect
}

func (ts *TravelScreen) HandleInput() error {
	ts.ui.Update()

	if store.GetResetTravel() == value.RESET_TRAVEL_FALSE_VALUE {
		ts.loadingStarsParticleEffect.Reset()

		ts.loadingStarsParticleEffect.HoldProgression()

		if store.GetStartSessionTravel() == value.START_SESSION_TRAVEL_TRUE_VALUE {
			utils.PerformLoadMap(loader.GetInstance().GetMap(loader.FirstMap), func(spawnables []dto.Position) {
				handler.PerformStartSession(
					store.GetSelectedSessionMetadata().ID,
					store.GetSelectedLobbySetUnitMetadata().ID,
					converter.ConvertPositionsToStartSessionSpawnables(spawnables),
					converter.ConvertPositionsToStartSessionChestLocations(spawnables),
					converter.ConvertPositionsToStartSessionHealthPackLocations(spawnables),
					func(err error) {
						if err != nil {
							notification.GetInstance().Push(
								common.ComposeMessage(
									translation.GetInstance().GetTranslation("client.networking.start-session-failure"),
									err.Error()),
								time.Second*3,
								common.NotificationErrorTextColor)

							lobby.GetInstance().CleanSelection()

							lobby.GetInstance().CleanListsEntries()

							lobby.GetInstance().HideStartButton()

							dispatcher.GetInstance().Dispatch(
								action.NewSetLobbySetRetrievalStartedNetworkingAction(value.LOBBY_SET_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE))

							dispatcher.GetInstance().Dispatch(
								action.NewSetSessionMetadataRetrievalStartedNetworkingAction(
									value.SESSION_METADATA_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE))

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

						lobby.GetInstance().CleanSelection()

						lobby.GetInstance().CleanListsEntries()

						lobby.GetInstance().HideStartButton()

						dispatcher.GetInstance().Dispatch(
							action.NewSetLobbySetRetrievalStartedNetworkingAction(value.LOBBY_SET_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE))

						dispatcher.GetInstance().Dispatch(
							action.NewSetSessionMetadataRetrievalStartedNetworkingAction(
								value.SESSION_METADATA_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE))

						ts.loadingStarsParticleEffect.ResumeProgression()
					})
			})

			dispatcher.GetInstance().Dispatch(
				action.NewSetStartSessionTravel(value.START_SESSION_TRAVEL_FALSE_VALUE))
		} else {
			utils.PerformLoadMap(loader.GetInstance().GetMap(loader.FirstMap), func(_ []dto.Position) {
				ts.loadingStarsParticleEffect.ResumeProgression()
			})
		}

		dispatcher.GetInstance().Dispatch(
			action.NewSetResetTravel(value.RESET_TRAVEL_TRUE_VALUE))
	}

	if !ts.loadingStarsParticleEffect.Done() {
		if !ts.loadingStarsParticleEffect.OnEnd() {
			ts.loadingStarsParticleEffect.Update()
		} else {
			ts.loadingStarsParticleEffect.Clean()

			dispatcher.GetInstance().Dispatch(
				action.NewSetResetTravel(value.RESET_TRAVEL_FALSE_VALUE))

			dispatcher.GetInstance().Dispatch(
				action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_SESSION_VALUE))
		}
	}

	return nil
}

func (ts *TravelScreen) HandleRender(screen *ebiten.Image) {
	ts.world.Clear()

	if !ts.loadingStarsParticleEffect.Done() {
		ts.loadingStarsParticleEffect.Draw(screen)
	}

	ts.ui.Draw(ts.world)

	screen.DrawImage(ts.world, &ebiten.DrawImageOptions{})
}

// newTravelScreen initializes TravelScreen.
func newTravelScreen() screen.Screen {
	return &TravelScreen{
		ui:                         builder.Build(),
		world:                      ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
		loadingStarsParticleEffect: loadingstars.NewStarsParticleEffect(),
	}
}
