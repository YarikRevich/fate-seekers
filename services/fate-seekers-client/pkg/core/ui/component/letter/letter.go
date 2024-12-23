package letter

import (
	"image/color"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// NewLetterComponent creates new session letter component.
func NewLetterComponent() *widget.Container {
	result := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(400, 300)),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.TrackHover(false)),
		widget.ContainerOpts.BackgroundImage(common.GetImageAsNineSlice(loader.PanelIdlePanel, 10, 10)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Left:  30,
				Right: 30,
				Top:   30,
			}),
		)))

	generalFont := &text.GoTextFace{
		Source: loader.GetInstance().GetFont(loader.KyivRegularFont),
		Size:   20,
	}

	result.AddChild(widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextOpts.Text("Letter", generalFont, color.White)))

	result.AddChild(widget.NewTextArea(
		widget.TextAreaOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.GridLayoutData{
					MaxHeight: 220,
				}))),
		widget.TextAreaOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
			Idle:     common.GetImageAsNineSlice(loader.ButtonIdleButton, 16, 0),
			Mask:     common.GetImageAsNineSlice(loader.ButtonIdleButton, 16, 0),
			Disabled: common.GetImageAsNineSlice(loader.ButtonIdleButton, 16, 0),
		})),
		widget.TextAreaOpts.SliderOpts(
			widget.SliderOpts.Images(nil, &widget.ButtonImage{
				Idle:         common.GetImageAsNineSlice(loader.ButtonIdleButton, 16, 0),
				Hover:        common.GetImageAsNineSlice(loader.ButtonIdleButton, 16, 0),
				Pressed:      common.GetImageAsNineSlice(loader.ButtonIdleButton, 16, 0),
				PressedHover: common.GetImageAsNineSlice(loader.ButtonIdleButton, 16, 0),
				Disabled:     common.GetImageAsNineSlice(loader.ButtonIdleButton, 16, 0),
			}),
			widget.SliderOpts.MinHandleSize(20),
			// widget.SliderOpts.TrackPadding(widget.Insets{}),
		),
		widget.TextAreaOpts.ShowVerticalScrollbar(),
		widget.TextAreaOpts.VerticalScrollMode(widget.PositionAtEnd),
		widget.TextAreaOpts.ProcessBBCode(true),
		widget.TextAreaOpts.FontFace(&text.GoTextFace{
			Source: loader.GetInstance().GetFont(loader.KyivRegularFont),
			Size:   25,
		}),
		widget.TextAreaOpts.FontColor(color.White),
		widget.TextAreaOpts.TextPadding(widget.Insets{
			Top:    20,
			Bottom: 20,
			Left:   20,
			Right:  20,
		}),
		widget.TextAreaOpts.Text("it wofjdkfjflfjlsfjlfjlfjldfjdlkfjdfjdkjfdkjfkdjfkdjfkjfkrks"),
	))

	closeContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(result.GetWidget().MinWidth, 0)),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(widget.Insets{
				Top: 20,
			}),
		)),
	)

	buttonIdleIcon := common.GetImageAsNineSlice(loader.ButtonIdleButton, 16, 0)
	buttonHoverIcon := common.GetImageAsNineSlice(loader.ButtonHoverButton, 16, 0)

	closeContainer.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:         buttonIdleIcon,
			Hover:        buttonHoverIcon,
			Pressed:      buttonIdleIcon,
			PressedHover: buttonIdleIcon,
			Disabled:     buttonIdleIcon,
		}),
		widget.ButtonOpts.Text("Close", generalFont, &widget.ButtonTextColor{Idle: color.White}),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
	))

	result.AddChild(closeContainer)

	return result
}
