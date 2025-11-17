package renderer

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/renderer/movable"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/kamera/v2"
)

// Renderer represents object renderer.
type Renderer struct {
	// Represents a set of movable objects, which are used for in the animator processing.
	movables *movable.Movables
}

// GetMovables retrieves configured animator movables holders.
func (r *Renderer) GetMovables() *movable.Movables {
	return r.movables
}

// Clean performs clean operation for the configured animator holders.
func (r *Renderer) Clean() {
	r.movables.Clean()
}

// Update performs update operation for the configured animator holders.
func (r *Renderer) Update() {
	r.movables.Update()
}

// Draw performs draw operation for the provided screens image for configured animator holders.
func (r *Renderer) Draw(screen *ebiten.Image, camera *kamera.Camera) {
	r.movables.Draw(screen, camera)
}

// NewRenderer initializes Renderer.
func NewRenderer() *Renderer {
	return &Renderer{
		movables: movable.NewMovables(),
	}
}
