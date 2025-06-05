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
	var backgroundAnimationGeometry ebiten.GeoM

	backgroundAnimationGeometry.Scale(
		scaler.GetScaleFactor(config.GetMinStaticWidth(), config.GetWorldWidth()),
		scaler.GetScaleFactor(config.GetMinStaticHeight(), config.GetWorldHeight()))

	shared.GetInstance().GetBackgroundAnimation().DrawTo(ss.world, &ebiten.DrawImageOptions{
		GeoM: backgroundAnimationGeometry,
	})

	ss.ui.Draw(ss.world)

	screen.DrawImage(ss.world, &ebiten.DrawImageOptions{
		ColorM: options.GetTransparentDrawOptions(
			ss.transparentTransitionEffect.GetValue()).ColorM})
}

// newSelectorScreen initializes SelectorScreen.
func newSelectorScreen() screen.Screen {
	return &SelectorScreen{
		ui: builder.Build(
			selector.NewSelectorComponent(
				func(sessionID string) {

				},
				func() {

				},
				func() {

				})),
		transparentTransitionEffect: transparent.NewTransparentTransitionEffect(true, 255, 0, 5, time.Microsecond*10),
		world:                       ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
	}
}
