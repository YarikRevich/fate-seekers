package inventory

import (
	"image/color"
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	// GetInstance retrieves instance of the inventory component, performing initial creation if needed.
	GetInstance = sync.OnceValue[*InventoryComponent](newInventoryComponent)
)

// InventoryComponent represents inventory component.
type InventoryComponent struct {
	container *widget.Container

	generalFont *text.GoTextFace

	buttonIdleIcon, buttonHoverIcon *image.NineSlice

	elements *widget.Container
}

func (ic *InventoryComponent) AddElements(elements []dto.InventoryElement) {
	for _, element := range elements {
		var graphic *widget.Graphic

		graphic = widget.NewGraphic(
			widget.GraphicOpts.Image(element.Image),
			widget.GraphicOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
					VerticalPosition:   widget.AnchorLayoutPositionCenter,
				}),
				widget.WidgetOpts.MouseButtonPressedHandler(func(args *widget.WidgetMouseButtonPressedEventArgs) {
					if args.Button == ebiten.MouseButtonRight {
						element.RemoveCallback(func() {
							ic.elements.RemoveChild(graphic)
						})
					} else if args.Button == ebiten.MouseButtonLeft {
						element.ApplyCallback()
					}
				}),
			),
		)

		ic.elements.AddChild(graphic)
	}
}

func (ic *InventoryComponent) CleanElements() {
	ic.elements.RemoveChildren()
}

func (ic *InventoryComponent) Show() {
	ic.container.GetWidget().Visibility = widget.Visibility_Show
}

func (ic *InventoryComponent) Hide() {
	ic.container.GetWidget().Visibility = widget.Visibility_Hide
}

// GetContainer retrieves container widget.
func (ic *InventoryComponent) GetContainer() *widget.Container {
	return ic.container
}

// newInventoryComponent creates new session inventory component.
func newInventoryComponent() *InventoryComponent {
	var result *InventoryComponent

	container := widget.NewContainer(
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

	container.AddChild(widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextOpts.Text("Inventory", generalFont, color.White)))

	elements := widget.NewContainer(
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

	buttonIdleIcon := common.GetImageAsNineSlice(loader.ButtonIdleButton, 16, 15)
	buttonHoverIcon := common.GetImageAsNineSlice(loader.ButtonHoverButton, 16, 15)

	container.AddChild(elements)

	container.GetWidget().Visibility = widget.Visibility_Hide

	result = &InventoryComponent{
		container:       container,
		generalFont:     generalFont,
		buttonIdleIcon:  buttonIdleIcon,
		buttonHoverIcon: buttonHoverIcon,
		elements:        elements,
	}

	return result
}
