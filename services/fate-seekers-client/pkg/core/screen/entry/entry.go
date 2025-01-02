package entry

import (
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition/transparent"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/builder"
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
	// dispatcher.GetInstance().Dispatch(
	// 	action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))
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
	return &EntryScreen{
		ui:                          builder.Build(),
		transparentTransitionEffect: transparent.NewTransparentTransitionEffect(),
		world:                       ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
	}
}
