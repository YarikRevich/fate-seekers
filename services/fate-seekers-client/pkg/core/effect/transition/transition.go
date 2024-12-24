package transition

import "github.com/hajimehoshi/ebiten/v2"

// TransitionEffect represents transition effects interface.
type TransitionEffect interface {
	// Done checks if transition has been finished.
	Done() bool

	// OnEnd checks if transition is on end state.
	OnEnd() bool

	// Update handles transition state update.
	Update()

	// Clean performes forced memory cleanup for the transition only.
	Clean()

	// GetOptions retrieves updated image draw options.
	GetOptions() *ebiten.DrawImageOptions
}
