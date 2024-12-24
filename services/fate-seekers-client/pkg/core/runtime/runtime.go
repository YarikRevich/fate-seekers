package runtime

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen/entry"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen/menu"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/application"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/hajimehoshi/ebiten/v2"
)

// Runtime represents main runtime flow implementation.
type Runtime struct {
	// Represents currently active screen.
	activeScreen screen.Screen
}

// Update performs logic update operations.
func (r *Runtime) Update() error {
	if store.GetInstance().GetState(application.EXIT_APPLICATION_STATE) ==
		value.EXIT_APPLICATION_TRUE_VALUE {
		return ebiten.Termination
	}

	switch store.GetActiveScreen() {
	case value.ACTIVE_SCREEN_ENTRY_VALUE:
		r.activeScreen = entry.GetInstance()

	case value.ACTIVE_SCREEN_MENU_VALUE:
		r.activeScreen = menu.GetInstance()
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
}

// Layout manages virtual world size.
func (r *Runtime) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.GetWorldWidth(), config.GetWorldHeight()
}

// NewRuntime creates new instance of Runtime.
func NewRuntime() *Runtime {
	return &Runtime{
		// Guarantees non blocking rendering, if state management fails.
		activeScreen: entry.GetInstance(),
	}
}
