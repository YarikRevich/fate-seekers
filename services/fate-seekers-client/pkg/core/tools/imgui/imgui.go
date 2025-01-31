package imgui

import (
	"fmt"
	"sync"

	imgui "github.com/gabstv/cimgui-go"
	ebimgui "github.com/gabstv/ebiten-imgui/v3"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// GetInstance retrieves instance of the imgui implementation, performing initilization if needed.
	GetInstance = sync.OnceValue[*ImGUI](newImGUI)
)

// ImGUI represents a wrapper for imgui implementation.
type ImGUI struct {
}

// Update updates imgui state.
func (i *ImGUI) Update() {
	ebimgui.Update(1.0 / 60.0)
	ebimgui.SetClipMask(!ebimgui.ClipMask())

	ebimgui.BeginFrame()

	imgui.Text(fmt.Sprintf("FPS: %f", ebiten.ActualFPS()))

	ebimgui.EndFrame()
}

// Draw performs imgui interface rendering.
func (i *ImGUI) Draw(screen *ebiten.Image) {
	ebimgui.Draw(screen)
}

// Layout performs imgui interface resizing.
func (i *ImGUI) Layout(outsideWidth, outsideHeight int) {
	ebimgui.SetDisplaySize(float32(outsideWidth), float32(outsideHeight))
}

// newImGUI initializes ImGUI.
func newImGUI() *ImGUI {
	return new(ImGUI)
}
