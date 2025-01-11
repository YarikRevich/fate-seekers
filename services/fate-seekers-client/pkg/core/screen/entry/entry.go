package entry

import (
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition/transparent"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/builder"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/notification"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/subtitles"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/storage/shared"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// GetInstance retrieves instance of the entry screen, performing initilization if needed.
	GetInstance = sync.OnceValue[screen.Screen](newEntryScreen)
)

// EntryScreen represents entry screen implementation.
type EntryScreen struct {
	// Represents attached user interface.
	ui *ebitenui.UI

	// Represents transparent transition effect.
	transparentTransitionEffect transition.TransitionEffect

	// Represents global world view.
	world *ebiten.Image

	// Represents stub timer, which is used to emulate some delay before screne switch.
	stubTimer *time.Timer
}

func (es *EntryScreen) HandleInput() error {
	if !es.transparentTransitionEffect.Done() {
		if !es.transparentTransitionEffect.OnEnd() {
			es.transparentTransitionEffect.Update()
		} else {
			es.transparentTransitionEffect.Clean()
		}
	}

	shared.GetInstance().GetBackgroundAnimation().Update()

	es.ui.Update()

	return nil
}

func (es *EntryScreen) HandleNetworking() {
	if store.GetEntryHandshakeStartedNetworking() == value.ENTRY_HANDSHAKE_STARTED_NETWORKING_FALSE_VALUE {
		es.stubTimer.Reset(time.Second * 3)

		dispatcher.GetInstance().Dispatch(
			action.NewSetLoadingApplicationAction(value.LOADING_APPLICATION_TRUE_VALUE))

		dispatcher.GetInstance().Dispatch(
			action.NewSetEntryHandshakeStartedNetworkingAction(value.ENTRY_HANDSHAKE_STARTED_NETWORKING_TRUE_VALUE))
	}

	select {
	case <-es.stubTimer.C:
		dispatcher.GetInstance().Dispatch(
			action.NewSetLoadingApplicationAction(value.LOADING_APPLICATION_FALSE_VALUE))

		dispatcher.GetInstance().Dispatch(
			action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))

		subtitles.GetInstance().Push("О, лягушка! Так дивно...", time.Second*6)
		subtitles.GetInstance().Push("'У багатих свої причуди!'", time.Second*6)

		notification.GetInstance().Push("Тестове повідомлення!", time.Second*6)
		notification.GetInstance().Push("Друге повідомлення!", time.Second*6)
	default:
	}
}

func (es *EntryScreen) HandleRender(screen *ebiten.Image) {
	var backgroundAnimationGeometry ebiten.GeoM

	backgroundAnimationGeometry.Scale(
		scaler.GetScaleFactor(config.GetMinStaticWidth(), config.GetWorldWidth()),
		scaler.GetScaleFactor(config.GetMinStaticHeight(), config.GetWorldHeight()))

	shared.GetInstance().GetBackgroundAnimation().DrawTo(es.world, &ebiten.DrawImageOptions{
		GeoM: backgroundAnimationGeometry,
	})

	es.ui.Draw(es.world)

	screen.DrawImage(es.world, &ebiten.DrawImageOptions{
		ColorM: es.transparentTransitionEffect.GetOptions().ColorM})
}

func (es *EntryScreen) Clean() {
}

// newEntryScreen initializes EntryScreen.
func newEntryScreen() screen.Screen {
	stubTimer := time.NewTimer(time.Minute)

	stubTimer.Stop()

	return &EntryScreen{
		ui:                          builder.Build(),
		transparentTransitionEffect: transparent.NewTransparentTransitionEffect(),
		world:                       ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
		stubTimer:                   stubTimer,
	}
}
