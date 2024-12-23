package builder

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
)

// Build composes user interface from the given components.
func Build(components ...widget.PreferredSizeLocateableWidget) *ebitenui.UI {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(config.GetWorldWidth(), config.GetWorldHeight()),
		),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	rootContainer.AddChild(components...)

	return &ebitenui.UI{
		Container: rootContainer,
	}
}
