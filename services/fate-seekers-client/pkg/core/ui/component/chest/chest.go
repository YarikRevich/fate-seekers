package chest

import (
	"image/color"
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/sound"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/common"
	componentscommon "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/translation"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	// GetInstance retrieves instance of the chest component, performing initial creation if needed.
	GetInstance = sync.OnceValue[*ChestComponent](newChestComponent)
)

// InventoryComponent represents inventory component.
type ChestComponent struct {
	// Represents container instance.
	container *widget.Container

	// Represents general font instance.
	generalFont *text.GoTextFace

	// Represent button icons.
	buttonIdleIcon, buttonHoverIcon *image.NineSlice

	// Represents close button callback.
	closeCallback func()

	// Represents a set of elements to be retrieved.
	elements *widget.Container
}

// AddElements adds new element to chest component.
func (cc *ChestComponent) AddElements(elements []dto.ChestElement) {
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
					if args.Button == ebiten.MouseButtonLeft {
						element.Callback(func() {
							cc.elements.RemoveChild(graphic)
						})
					}
				}),
			),
		)

		cc.elements.AddChild(graphic)
	}
}

// CleanElements removes all the elements.
func (cc *ChestComponent) CleanElements() {
	cc.elements.RemoveChildren()
}

// Show shows chest component.
func (cc *ChestComponent) Show() {
	cc.container.GetWidget().Visibility = widget.Visibility_Show
}

// Hide hides chest component.
func (cc *ChestComponent) Hide() {
	cc.container.GetWidget().Visibility = widget.Visibility_Hide
}

// SetCloseCallback modified close callback in the container.
func (cc *ChestComponent) SetCloseCallback(callback func()) {
	cc.closeCallback = callback
}

// GetContainer retrieves container widget.
func (cc *ChestComponent) GetContainer() *widget.Container {
	return cc.container
}

// newChestComponent creates new session chest component.
func newChestComponent() *ChestComponent {
	var result *ChestComponent

	container := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(
				scaler.GetPercentageOf(config.GetWorldWidth(), 30),
				scaler.GetPercentageOf(config.GetWorldHeight(), 20)),
			widget.WidgetOpts.TrackHover(false),
		),
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
		widget.TextOpts.Text(
			translation.GetInstance().GetTranslation("client.chest.title"),
			generalFont,
			color.White)))

	elements := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(
				scaler.GetPercentageOf(config.GetWorldWidth(), 30),
				scaler.GetPercentageOf(config.GetWorldHeight(), 10)),
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
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

	closeContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(container.GetWidget().MinWidth, 0)),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(widget.Insets{
				Top:    50,
				Bottom: 25,
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
		widget.ButtonOpts.Text(
			translation.GetInstance().GetTranslation("client.chest.close"),
			generalFont,
			&widget.ButtonTextColor{Idle: componentscommon.ButtonTextColor}),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
			})),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    20,
			Bottom: 20,
		}),
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			sound.GetInstance().GetSoundUIFxManager().PushWithHandbrake(loader.ButtonFXSound)

			result.closeCallback()
		}),
	))

	container.AddChild(closeContainer)

	container.GetWidget().Visibility = widget.Visibility_Hide

	result = &ChestComponent{
		container:       container,
		generalFont:     generalFont,
		buttonIdleIcon:  buttonIdleIcon,
		buttonHoverIcon: buttonHoverIcon,
		elements:        elements,
	}

	return result
}
