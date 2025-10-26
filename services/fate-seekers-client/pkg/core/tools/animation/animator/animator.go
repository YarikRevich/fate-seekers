package animator

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/animation/animator/movable"
	"github.com/hajimehoshi/ebiten/v2"
)

// TODO: create animator with direction objects

// Animator represents movable animation animator.
type Animator struct {
	// Represents a set of movable objects, which are used for in the animator processing.
	movables *movable.Movables
}

func (a *Animator) GetMovables() *movable.Movables {
	return a.movables
}

func (a *Animator) Clean() {
	a.movables.Clean()
}

func (a *Animator) Update() {
	for _, movable := range a.movables {
		movable.Update()
	}
}

func (a *Animator) Draw(screen *ebiten.Image) {
	for _, movable := range a.movables {
		movable.Draw(screen)
	}
}

// NewAnimator initializes Animator.
func NewAnimator() *Animator {
	return &Animator{
		movables: movable.NewMovables(),
	}
}
