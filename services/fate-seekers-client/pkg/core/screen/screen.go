package screen

import "github.com/hajimehoshi/ebiten/v2"

// Reducer represents common screen interface.
type Screen interface {
	// HandleInput handles provided user input.
	HandleInput() error

	// HandleNetworking handles network communication.
	HandleNetworking()

	// HandleRender handles render operation on the main screen.
	HandleRender(screen *ebiten.Image)

	// Clean performes forced memory cleanup for the screen only.
	Clean() // TODO: remove
}
