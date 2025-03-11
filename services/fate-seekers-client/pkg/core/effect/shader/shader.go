package shader

import "github.com/hajimehoshi/ebiten/v2"

// ShaderEffect represents shader effects interface.
type ShaderEffect interface {
	// Draw performs draw operation for the shader composition.
	Draw(screen *ebiten.Image)
}
