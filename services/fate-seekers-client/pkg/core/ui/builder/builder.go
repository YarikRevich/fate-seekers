package builder

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
)

// Build composes user interface from the given components.
func Build(components ...widget.PreferredSizeLocateableWidget) *ebitenui.UI {
	rootContainer := widget.NewContainer()

	rootContainer.AddChild(components...)

	return &ebitenui.UI{
		Container: rootContainer,
	}
}
