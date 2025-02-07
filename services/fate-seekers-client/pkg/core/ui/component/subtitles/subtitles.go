package subtitles

import (
	"image/color"
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	// GetInstance retrieves instance of the subtitles component, performing initial creation if needed.
	GetInstance = sync.OnceValue[*SubtitlesComponent](newSubtitlesComponent)
)

// SubtitlesComponent represents component, which contains actual subtitles.
type SubtitlesComponent struct {
	container *widget.Container
}

// SetText modifies text component in the container.
func (sc *SubtitlesComponent) SetText(value string) {
	sc.container.AddChild(widget.NewText(
		widget.TextOpts.MaxWidth(float64(scaler.GetPercentageOf(config.GetWorldWidth(), 30))),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionStart,
			StretchHorizontal:  false,
			StretchVertical:    false,
		})),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionStart),
		widget.TextOpts.Insets(widget.Insets{
			Top:    20,
			Bottom: 20,
			Left:   20,
			Right:  20,
		}),
		widget.TextOpts.ProcessBBCode(true),
		widget.TextOpts.Text(
			value,
			&text.GoTextFace{
				Source: loader.GetInstance().GetFont(loader.KyivRegularFont),
				Size:   20,
			},
			color.White)))
}

// CleanText cleans text component in the container.
func (sc *SubtitlesComponent) CleanText() {
	sc.container.RemoveChildren()
}

// GetContainer retrieves container widget.
func (sc *SubtitlesComponent) GetContainer() *widget.Container {
	return sc.container
}

// newSubtitlesComponent initializes SubtitlesComponent.
func newSubtitlesComponent() *SubtitlesComponent {
	container := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(common.GetImageAsNineSlice(loader.PanelIdlePanel, 10, 10)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.TrackHover(false),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				Padding:            widget.Insets{Bottom: 10},
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	return &SubtitlesComponent{container: container}
}
