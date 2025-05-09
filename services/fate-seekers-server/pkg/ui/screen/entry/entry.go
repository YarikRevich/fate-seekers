package entry

import (
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/storage/shared"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/effect/transition"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/effect/transition/transparent"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/value"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/tools/options"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/ui/builder"
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

	// select {
	// case <-es.stubTimer.C:
	// 	dispatcher.GetInstance().Dispatch(
	// 		action.NewSetLoadingApplicationAction(value.LOADING_APPLICATION_FALSE_VALUE))

	dispatcher.GetInstance().Dispatch(
		action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))

	// 	subtitles.GetInstance().Push("О, лягушка! Так дивно...", time.Second*6)
	// 	subtitles.GetInstance().Push("'У багатих свої причуди!'", time.Second*6)

	// 	// notification.GetInstance().Push("Тестове повідомлення!", time.Second*6)
	// 	// notification.GetInstance().Push("Друге повідомлення!", time.Second*6)
	// default:
	// }

	es.ui.Update()

	return nil
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
		ColorM: options.GetTransparentDrawOptions(
			es.transparentTransitionEffect.GetValue()).ColorM})
}

// newEntryScreen initializes EntryScreen.
func newEntryScreen() screen.Screen {
	stubTimer := time.NewTimer(time.Minute)

	stubTimer.Stop()

	return &EntryScreen{
		ui:                          builder.Build(),
		transparentTransitionEffect: transparent.NewTransparentTransitionEffect(true, 255, 0, 5, time.Microsecond*10),
		world:                       ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
		stubTimer:                   stubTimer,
	}
}
