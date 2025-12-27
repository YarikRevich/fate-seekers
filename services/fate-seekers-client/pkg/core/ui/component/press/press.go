package press

import (
	"image/color"
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/translation"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	// GetInstance retrieves instance of the press component, performing initial creation if needed.
	GetInstance = sync.OnceValue[*PressComponent](newPressComponent)
)

// Describes all the available press component text types.
const (
	GAMEPAD = iota
	KEYBOARD
)

// PressComponent represents component, which contains user press.
type PressComponent struct {
	// Represents text.
	text *widget.Text

	// Represents container widget.
	container *widget.Container
}

// SetPressType sets press type by the provided value for press widget.
func (pc *PressComponent) SetPressType(value int) {
	switch value {
	case GAMEPAD:
		pc.text.Label = translation.GetInstance().GetTranslation("client.press.gamepad")

	case KEYBOARD:
		pc.text.Label = translation.GetInstance().GetTranslation("client.press.keyboard")
	}
}

// GetContainer retrieves container widget.
func (pc *PressComponent) GetContainer() *widget.Container {
	return pc.container
}

// newPressComponent creates new press component.
func newPressComponent() *PressComponent {
	var result *PressComponent

	container := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(
				config.GetWorldWidth(),
				scaler.GetPercentageOf(config.GetWorldHeight(), 8),
			),
			widget.WidgetOpts.TrackHover(false),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:  widget.AnchorLayoutPositionStart,
				StretchHorizontal: false,
				StretchVertical:   false,
			})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(
				scaler.GetPercentageOf(config.GetWorldWidth(), 58)),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top:    scaler.GetPercentageOf(config.GetWorldWidth(), 1),
				Left:   scaler.GetPercentageOf(config.GetWorldWidth(), 3),
				Right:  scaler.GetPercentageOf(config.GetWorldWidth(), 3),
				Bottom: scaler.GetPercentageOf(config.GetWorldWidth(), 1),
			}),
		)))

	generalFont := &text.GoTextFace{
		Source: loader.GetInstance().GetFont(loader.KyivRegularFont),
		Size:   20,
	}

	pressTextContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.TrackHover(false),
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionEnd,
				Stretch:  true,
			})),
		widget.ContainerOpts.BackgroundImage(common.GetImageAsNineSlice(loader.PanelIdlePanel, 10, 10)),
		widget.ContainerOpts.Layout(
			widget.NewAnchorLayout(
				widget.AnchorLayoutOpts.Padding(widget.Insets{
					Left:  scaler.GetPercentageOf(config.GetWorldWidth(), 2),
					Right: scaler.GetPercentageOf(config.GetWorldWidth(), 2),
				}))))

	text := widget.NewText(
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.TextOpts.Text(
			"",
			generalFont,
			color.White))

	pressTextContainer.AddChild(text)

	container.AddChild(pressTextContainer)

	result = &PressComponent{
		text:      text,
		container: container,
	}

	return result
}
