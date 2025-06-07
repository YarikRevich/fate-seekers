package selector

import (
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition/transparent"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/options"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/builder"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/selector"
	selectormanager "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/selector"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/storage/shared"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
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

	selector.GetInstance().SetSubmitCallback(func(sessionID string) {
		if selectormanager.ProcessChanges(sessionID) {
			transparentTransitionEffect.Reset()

			dispatcher.GetInstance().Dispatch(
				action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_LOBBY_VALUE))
		}
	})

	selector.GetInstance().SetCreateCallback(func() {
		transparentTransitionEffect.Reset()

		dispatcher.GetInstance().Dispatch(
			action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_CREATOR_VALUE))
	})

	selector.GetInstance().SetBackCallback(func() {
		transparentTransitionEffect.Reset()

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
