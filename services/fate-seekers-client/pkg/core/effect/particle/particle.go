package particle

import "github.com/hajimehoshi/ebiten/v2"

// ParticleEffect represents particles effects interface.
type ParticleEffect interface {
	// Done checks if particle effect has been finished.
	Done() bool

	// OnEnd checks if particle effect is on end state.
	OnEnd() bool

	// Clean performes forced memory cleanup for the particle effect only.
	Clean()

	// Reset performes particle effect state reset.
	Reset()

	// Update performs update operation for all particles.
	Update()

	// Draw performs draw operation for all particles.
	Draw(screen *ebiten.Image)
}
