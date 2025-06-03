package logo

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
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/repository"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/repository/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/storage/shared"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// GetInstance retrieves instance of the logo screen, performing initilization if needed.
	GetInstance = sync.OnceValue[screen.Screen](newLogoScreen)
)

// LogoScreen represents logo screen implementation.
type LogoScreen struct {
	// Represents attached user interface.
	ui *ebitenui.UI

	// Represents transparent transition effect.
	transparentTransitionEffect transition.TransitionEffect

	// Represents global world view.
	world *ebiten.Image
}

func (ls *LogoScreen) HandleInput() error {
	_, ok, err := repository.GetFlagsRepository().GetByName(common.UUID_FLAG_NAME)

	dispatcher.GetInstance().Dispatch(
		action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_ENTRY_VALUE))

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

func (ls *LogoScreen) HandleRender(screen *ebiten.Image) {
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

// newLogoScreen initializes LogoScreen.
func newLogoScreen() screen.Screen {
	return &LogoScreen{
		ui:                          builder.Build(),
		transparentTransitionEffect: transparent.NewTransparentTransitionEffect(true, 255, 0, 5, time.Microsecond*10),
		world:                       ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
	}
}
