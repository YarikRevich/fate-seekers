package screen

import "github.com/hajimehoshi/ebiten/v2"

// Reducer represents common screen interface.
type Screen interface {
	// HandleInput handles provided user input.
	HandleInput() error

	// HandleRender handles render operation on the main screen.
	HandleRender(screen *ebiten.Image)
}
