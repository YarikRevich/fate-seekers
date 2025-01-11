package runtime

import (
	"github.com/Frabjous-Studios/asebiten"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen/entry"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen/intro"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen/menu"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen/settings"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/builder"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/subtitles"
	subtitlesmanager "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/subtitles"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/application"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
)

// Runtime represents main runtime flow implementation.
type Runtime struct {
	// Represents attached user interface.
	ui *ebitenui.UI

	// Represents currently active screen.
	activeScreen screen.Screen

	// Represents global loader animation.
	loaderAnimation *asebiten.Animation
}

// Update performs logic update operations.
func (r *Runtime) Update() error {
	if store.GetExitApplication() == value.EXIT_APPLICATION_TRUE_VALUE {
		return ebiten.Termination
	}

	if store.GetLoadingApplication() == value.LOADING_APPLICATION_TRUE_VALUE {
		r.loaderAnimation.Update()
	}

	subtitlesmanager.GetInstance().Update()

	switch store.GetActiveScreen() {
	case value.ACTIVE_SCREEN_INTRO_VALUE:
		r.activeScreen = intro.GetInstance()

	case value.ACTIVE_SCREEN_ENTRY_VALUE:
		r.activeScreen = entry.GetInstance()

	case value.ACTIVE_SCREEN_MENU_VALUE:
		r.activeScreen = menu.GetInstance()

	case value.ACTIVE_SCREEN_SETTINGS_VALUE:
		r.activeScreen = settings.GetInstance()
	}

	err := r.activeScreen.HandleInput()
	if err != nil {
		return err
	}

	r.activeScreen.HandleNetworking()

	return nil
}

// Draw performs render operation.
func (r *Runtime) Draw(screen *ebiten.Image) {
	r.activeScreen.HandleRender(screen)

	if store.GetInstance().GetState(application.LOADING_APPLICATION_STATE) ==
		value.LOADING_APPLICATION_TRUE_VALUE {
		var loadingAnimationGeometry ebiten.GeoM

		loadingAnimationGeometry.Translate(
			float64(scaler.GetPercentageOf(config.GetWorldWidth(), 2)),
			float64(scaler.GetPercentageOf(config.GetWorldHeight(), 91)))

		r.loaderAnimation.DrawTo(screen, &ebiten.DrawImageOptions{GeoM: loadingAnimationGeometry})
	}

	r.ui.Draw(screen)
}

// Layout manages virtual world size.
func (r *Runtime) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.GetWorldWidth(), config.GetWorldHeight()
}

// NewRuntime creates new instance of Runtime.
func NewRuntime() *Runtime {
	return &Runtime{
		ui: builder.Build(subtitles.GetInstance().GetContainer()),
		// Guarantees non blocking rendering, if state management fails.
		activeScreen:    entry.GetInstance(),
		loaderAnimation: loader.GetInstance().GetAnimation(loader.LoaderAnimation, true),
	}
}
