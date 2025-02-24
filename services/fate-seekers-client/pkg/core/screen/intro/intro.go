package intro

import (
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition/transparent"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/builder"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/storage/shared"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// GetInstance retrieves instance of the intro screen, performing initilization if needed.
	GetInstance = sync.OnceValue[screen.Screen](newIntroScreen)
)

// IntroScreen represents intro screen implementation.
type IntroScreen struct {
	// Represents attached user interface.
	ui *ebitenui.UI

	// Represents transparent transition effect.
	transparentTransitionEffect transition.TransitionEffect

	// Represents global world view.
	world *ebiten.Image
}

func (es *IntroScreen) HandleInput() error {
	// TODO: check if intro scene has already been played
	// repository.GetFlagsRepository().GetByName()

	dispatcher.GetInstance().Dispatch(
		action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_ENTRY_VALUE))

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

func (is *IntroScreen) HandleNetworking() {
	// dispatcher.GetInstance().Dispatch(
	// 	action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))
}

func (is *IntroScreen) HandleRender(screen *ebiten.Image) {
	var backgroundAnimationGeometry ebiten.GeoM

	backgroundAnimationGeometry.Scale(
		scaler.GetScaleFactor(config.GetMinStaticWidth(), config.GetWorldWidth()),
		scaler.GetScaleFactor(config.GetMinStaticHeight(), config.GetWorldHeight()))

	shared.GetInstance().GetBackgroundAnimation().DrawTo(is.world, &ebiten.DrawImageOptions{
		GeoM: backgroundAnimationGeometry,
	})

	is.ui.Draw(is.world)

	screen.DrawImage(is.world, &ebiten.DrawImageOptions{
		ColorM: is.transparentTransitionEffect.GetOptions().ColorM})
}

func (es *IntroScreen) Clean() {
}

// newIntroScreen initializes IntroScreen.
func newIntroScreen() screen.Screen {
	return &IntroScreen{
		ui:                          builder.Build(),
		transparentTransitionEffect: transparent.NewTransparentTransitionEffect(),
		world:                       ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
	}
}
