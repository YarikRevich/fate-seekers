package selector

import (
	"image/color"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/common"
	componentscommon "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/translation"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	// Describes max amount of symbols, which can be entered to input component.
	maxInputSymbols = 30
)

// Describes all the colors used for list combo definition.
var (
	selectedListColor = color.NRGBA{183, 228, 202, 255}
	focusedListColor  = color.NRGBA{R: 170, G: 170, B: 180, A: 255}
	disabledListColor = color.NRGBA{100, 100, 100, 255}
)

// NewSelectorComponent creates new selector component.
func NewSelectorComponent(
	submitCallback func(sessionID string),
	createCallback func(),
	backCallback func()) *widget.Container {
	result := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(
				scaler.GetPercentageOf(config.GetWorldWidth(), 40),
				scaler.GetPercentageOf(config.GetWorldHeight(), 30)),
			widget.WidgetOpts.TrackHover(false),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				Padding: widget.Insets{
					Left: scaler.GetPercentageOf(config.GetWorldWidth(), 9),
				},
				VerticalPosition:  widget.AnchorLayoutPositionCenter,
				StretchHorizontal: false,
				StretchVertical:   false,
			})),
		widget.ContainerOpts.BackgroundImage(common.GetImageAsNineSlice(loader.PanelIdlePanel, 10, 10)),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Left:   30,
				Right:  30,
				Top:    30,
				Bottom: 30,
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
		widget.TextOpts.Text(
			translation.GetInstance().GetTranslation("client.selector.title"),
			generalFont,
			color.White)))

	components := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{true, true}, nil),
			widget.GridLayoutOpts.Spacing(
				10, scaler.GetPercentageOf(config.GetWorldHeight(), 5)),
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top: scaler.GetPercentageOf(config.GetWorldHeight(), 6),
			}))))

	components.AddChild(widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextOpts.Text(
			translation.GetInstance().GetTranslation("client.settings.host"),
			generalFont,
			color.White)))

	var networkingHostInput *widget.TextInput

	networkingHostInput = widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(
				scaler.GetPercentageOf(config.GetWorldWidth(), 20), 0),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
				Padding: widget.Insets{
					Bottom: scaler.GetPercentageOf(config.GetWorldHeight(), 20),
				},
			})),
		widget.TextInputOpts.Image(&widget.TextInputImage{
			Idle:     image.NewNineSlice(loader.GetInstance().GetStatic(loader.TextInputIdle), [3]int{9, 14, 6}, [3]int{9, 14, 6}),
			Disabled: image.NewNineSlice(loader.GetInstance().GetStatic(loader.TextInputIdle), [3]int{9, 14, 6}, [3]int{9, 14, 6}),
		}),
		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:          color.White,
			Disabled:      color.White,
			Caret:         color.White,
			DisabledCaret: color.White,
		}),
		widget.TextInputOpts.Padding(widget.Insets{
			Left:   13,
			Right:  13,
			Top:    13,
			Bottom: 13,
		}),
		widget.TextInputOpts.Face(&text.GoTextFace{
			Source: loader.GetInstance().GetFont(loader.KyivRegularFont),
			Size:   20,
		}),
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(generalFont, 4),
		),
		widget.TextInputOpts.AllowDuplicateSubmit(false),
		widget.TextInputOpts.Validation(func(newInputTextRaw string) (bool, *string) {
			newInputText := newInputTextRaw

			parsedNewInputText := newInputText[len(networkingHostInput.GetText()):]

			if len(parsedNewInputText) > 1 {
				newInputText = networkingHostInput.GetText() + parsedNewInputText[:1]
			}

			if len(newInputText) >= maxInputSymbols {
				replacement := networkingHostInput.GetText()

				return false, &replacement
			}

			return false, &newInputText
		}))

	networkingHostInput.SetText(config.GetSettingsNetworkingServerHost())

	components.AddChild(networkingHostInput)

	buttonsContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(
				result.GetWidget().MinWidth,
				result.GetWidget().MinHeight),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(13),
		)),
	)

	buttonIdleIcon := common.GetImageAsNineSlice(loader.ButtonIdleButton, 16, 15)
	buttonHoverIcon := common.GetImageAsNineSlice(loader.ButtonHoverButton, 16, 15)

	buttonsContainer.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:         buttonIdleIcon,
			Hover:        buttonHoverIcon,
			Pressed:      buttonIdleIcon,
			PressedHover: buttonIdleIcon,
			Disabled:     buttonIdleIcon,
		}),
		widget.ButtonOpts.Text(
			translation.GetInstance().GetTranslation("client.selector.submit"),
			generalFont,
			&widget.ButtonTextColor{Idle: componentscommon.ButtonTextColor}),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionEnd,
			}),
		),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    20,
			Bottom: 20,
		}),
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			submitCallback("it works")
		}),
	))

	buttonsContainer.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:         buttonIdleIcon,
			Hover:        buttonHoverIcon,
			Pressed:      buttonIdleIcon,
			PressedHover: buttonIdleIcon,
			Disabled:     buttonIdleIcon,
		}),
		widget.ButtonOpts.Text(
			translation.GetInstance().GetTranslation("client.selector.create"),
			generalFont,
			&widget.ButtonTextColor{Idle: componentscommon.ButtonTextColor}),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionEnd,
			}),
		),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    20,
			Bottom: 20,
		}),
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			createCallback()
		}),
	))

	buttonsContainer.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:         buttonIdleIcon,
			Hover:        buttonHoverIcon,
			Pressed:      buttonIdleIcon,
			PressedHover: buttonIdleIcon,
			Disabled:     buttonIdleIcon,
		}),
		widget.ButtonOpts.Text(
			translation.GetInstance().GetTranslation("client.selector.close"),
			generalFont,
			&widget.ButtonTextColor{Idle: componentscommon.ButtonTextColor}),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionEnd,
			})),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    20,
			Bottom: 20,
		}),
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			backCallback()
		}),
	))

	result.AddChild(buttonsContainer)

	return result
}
