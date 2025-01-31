package travel

import (
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/builder"
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
}

func (ts *TravelScreen) HandleInput() error {
	ts.ui.Update()

	return nil
}

func (ts *TravelScreen) HandleNetworking() {

}

func (ts *TravelScreen) HandleRender(screen *ebiten.Image) {
	ts.world.Clear()

	ts.ui.Draw(ts.world)

	screen.DrawImage(ts.world, &ebiten.DrawImageOptions{})
}

func (ts *TravelScreen) Clean() {

}

// newTravelScreen initializes TravelScreen.
func newTravelScreen() screen.Screen {
	return &TravelScreen{
		ui:    builder.Build(),
		world: ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
	}
}
