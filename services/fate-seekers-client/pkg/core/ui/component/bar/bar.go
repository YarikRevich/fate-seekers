package bar

import (
	"fmt"
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
	// GetInstance retrieves instance of the bar component, performing initial creation if needed.
	GetInstance = sync.OnceValue[*BarComponent](newBarComponent)
)

// BarComponent represents component, which contains user bar.
type BarComponent struct {
	// Represents health text.
	healthText *widget.Text

	// Represents container widget.
	container *widget.Container
}

// SetHealthText sets label by the provided value for health text widget.
func (bc *BarComponent) SetHealthText(value uint64) {
	bc.healthText.Label = fmt.Sprintf("%d%%", value)
}

// GetContainer retrieves container widget.
func (bc *BarComponent) GetContainer() *widget.Container {
	return bc.container
}

// newBarComponent creates new bar component.
func newBarComponent() *BarComponent {
	var result *BarComponent

	container := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(
				config.GetWorldWidth(),
				scaler.GetPercentageOf(config.GetWorldHeight(), 8),
			),
			widget.WidgetOpts.TrackHover(false),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:  widget.AnchorLayoutPositionEnd,
				StretchHorizontal: false,
				StretchVertical:   false,
			})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(
				scaler.GetPercentageOf(config.GetWorldWidth(), 69)),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Left:   scaler.GetPercentageOf(config.GetWorldWidth(), 3),
				Right:  scaler.GetPercentageOf(config.GetWorldWidth(), 3),
				Bottom: scaler.GetPercentageOf(config.GetWorldWidth(), 1),
			}),
		)))

	health := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.TrackHover(false),
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionStart,
				MaxWidth: scaler.GetPercentageOf(config.GetWorldWidth(), 10),
				Stretch:  true,
			})),
		widget.ContainerOpts.BackgroundImage(common.GetImageAsNineSlice(loader.PanelIdlePanel, 10, 10)),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(
				scaler.GetPercentageOf(config.GetWorldWidth(), 1)),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Left:  scaler.GetPercentageOf(config.GetWorldWidth(), 2),
				Right: scaler.GetPercentageOf(config.GetWorldWidth(), 2),
			}),
		)))

	health.AddChild(widget.NewGraphic(
		widget.GraphicOpts.Image(loader.GetInstance().GetStatic(loader.Heart)),
		widget.GraphicOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionStart,
				Stretch:  true,
			}),
		),
	))

	generalFont := &text.GoTextFace{
		Source: loader.GetInstance().GetFont(loader.KyivRegularFont),
		Size:   20,
	}

	healthText := widget.NewText(
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: widget.RowLayoutPositionEnd,
			Stretch:  true,
		})),
		widget.TextOpts.Text(
			"100%",
			generalFont,
			color.White))

	health.AddChild(healthText)

	container.AddChild(health)

	weapon := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(
				scaler.GetPercentageOf(config.GetWorldWidth(), 15),
				scaler.GetPercentageOf(config.GetWorldHeight(), 7),
			),
			widget.WidgetOpts.TrackHover(false),
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionEnd,
				MaxWidth: scaler.GetPercentageOf(config.GetWorldWidth(), 15),
				Stretch:  true,
			})),
		widget.ContainerOpts.BackgroundImage(common.GetImageAsNineSlice(loader.PanelIdlePanel, 10, 10)))

	container.AddChild(weapon)

	result = &BarComponent{
		healthText: healthText,
		container:  container,
	}

	return result
}
