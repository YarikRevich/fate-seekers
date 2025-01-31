package answerinput

import (
	"image/color"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	// Describes max amount of symbols, which can be entered to input component.
	maxInputSymbols = 20
)

// NewAnswerInputComponent creates new answer input component.
func NewAnswerInputComponent(submitCallback, closeCallback func()) *widget.Container {
	result := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(
				config.GetWorldWidth(),
				config.GetWorldHeight()),
			widget.WidgetOpts.TrackHover(false),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				Padding: widget.Insets{
					Left: scaler.GetPercentageOf(config.GetWorldWidth(), 15),
				},
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
			})),
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

	var answerInput *widget.TextInput

	answerInput = widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch:  true,
			Position: widget.RowLayoutPositionStart,
			MaxWidth: 200,
		})),
		widget.TextInputOpts.Image(&widget.TextInputImage{
			Idle:     image.NewNineSlice(loader.GetInstance().GetStatic(loader.TextInputIdle), [3]int{9, 14, 6}, [3]int{9, 14, 6}),
			Disabled: image.NewNineSlice(loader.GetInstance().GetStatic(loader.TextInputIdle), [3]int{9, 14, 6}, [3]int{9, 14, 6}),
		}),
		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:          color.Black,
			Disabled:      color.Black,
			Caret:         color.Black,
			DisabledCaret: color.Black,
		}),
		widget.TextInputOpts.Padding(widget.Insets{
			Left:   13,
			Right:  13,
			Top:    13,
			Bottom: 13,
		}),
		widget.TextInputOpts.Face(&text.GoTextFace{
			Source: loader.GetInstance().GetFont(loader.KyivRegularFont),
			Size:   28,
		}),
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(generalFont, 4),
		),
		widget.TextInputOpts.AllowDuplicateSubmit(false),
		widget.TextInputOpts.Validation(func(newInputTextRaw string) (bool, *string) {
			newInputText := newInputTextRaw

			parsedNewInputText := newInputText[len(answerInput.GetText()):]

			if len(parsedNewInputText) > 1 {
				newInputText = answerInput.GetText() + parsedNewInputText[:1]
			}

			if len(newInputText) >= maxInputSymbols {
				replacement := answerInput.GetText()

				return false, &replacement
			}

			return false, &newInputText
		}),
		widget.TextInputOpts.Placeholder("Enter text here"))

	result.AddChild(answerInput)

	closeContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(result.GetWidget().MinWidth, 40)),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	buttonIdleIcon := common.GetImageAsNineSlice(loader.ButtonIdleButton, 16, 15)
	buttonHoverIcon := common.GetImageAsNineSlice(loader.ButtonHoverButton, 16, 15)

	buttonsContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				StretchHorizontal:  true,
				StretchVertical:    false,
			}),
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				// MaxWidth: 100,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(30),
		)))

	buttonsContainer.AddChild(widget.NewButton(
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
				StretchHorizontal:  true,
				StretchVertical:    false,
			}),
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			})),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    20,
			Bottom: 20,
		}),
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			closeCallback()
		}),
	))

	closeContainer.AddChild(buttonsContainer)

	result.AddChild(closeContainer)

	return result
}
