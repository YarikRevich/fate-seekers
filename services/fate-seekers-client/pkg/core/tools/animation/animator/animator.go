package animator

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/animation/animator/movable"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/kamera/v2"
)

// Animator represents movable animation animator.
type Animator struct {
	// Represents a set of movable objects, which are used for in the animator processing.
	movables *movable.Movables
}

// GetMovables retrieves configured animator movables holders.
func (a *Animator) GetMovables() *movable.Movables {
	return a.movables
}

// Clean performs clean operation for the configured animator holders.
func (a *Animator) Clean() {
	a.movables.Clean()
}

// Update performs update operation for the configured animator holders.
func (a *Animator) Update() {
	a.movables.Update()
}

// Draw performs draw operation for the provided screens image for configured animator holders.
func (a *Animator) Draw(screen *ebiten.Image, camera *kamera.Camera) {
	a.movables.Draw(screen, camera)
}

// NewAnimator initializes Animator.
func NewAnimator() *Animator {
	return &Animator{
		movables: movable.NewMovables(),
	}
}
