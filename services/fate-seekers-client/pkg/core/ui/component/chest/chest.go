package chest

import (
	"fmt"
	"image/color"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// NewChestComponent creates new session chest component.
func NewChestComponent() *widget.Container {
	//

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
		widget.TextOpts.Text("Chest", generalFont, color.White)))

	bc := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(4),
			widget.GridLayoutOpts.Stretch([]bool{true, true, true, true}, nil),
			widget.GridLayoutOpts.Spacing(10, 10),
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top: 20,
			}))))

	buttonIdleIcon := common.GetImageAsNineSlice(loader.ButtonIdleButton, 16, 0)
	buttonHoverIcon := common.GetImageAsNineSlice(loader.ButtonHoverButton, 16, 0)

	i := 0
	for row := 0; row < 3; row++ {
		for col := 0; col < 4; col++ {
			b := widget.NewButton(
				widget.ButtonOpts.Image(&widget.ButtonImage{
					Idle:         buttonIdleIcon,
					Hover:        buttonHoverIcon,
					Pressed:      buttonIdleIcon,
					PressedHover: buttonIdleIcon,
					Disabled:     buttonIdleIcon,
				}),
				widget.ButtonOpts.Text(fmt.Sprintf("%s %d", string(rune('A'+i)), i+1), generalFont, &widget.ButtonTextColor{Idle: color.White}))
			bc.AddChild(b)

			i++
		}
	}

	result.AddChild(bc)

	closeContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(result.GetWidget().MinWidth, 0)),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(widget.Insets{
				Top: 20,
			}),
		)),
	)

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
