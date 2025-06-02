package answerinput

import (
	"fmt"
	"image/color"
	"sync"

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
	maxInputSymbols = 20
)

var (
	// GetInstance retrieves instance of the letter component, performing initial creation if needed.
	GetInstance = sync.OnceValue[*AnswerInputComponent](newAnswerInputComponent)
)

// AnswerInputComponent represents component, which contains actual answer input.
type AnswerInputComponent struct {
	// Represents text widget.
	text *widget.Text

	// Represents submit callback.
	submitCallback func(valueRaw string)

	// Represents close callback.
	closeCallback func()

	// Represents container widget.
	container *widget.Container
}

// SetText modifies text component in the container.
func (aic *AnswerInputComponent) SetText(value string) {
	aic.text.Label = fmt.Sprintf(
		"%s:   %s",
		translation.GetInstance().GetTranslation("client.answerinput.solvetext"),
		value)
}

// SetSubmitCallback modified submit callback in the container.
func (aic *AnswerInputComponent) SetSubmitCallback(callback func(valueRaw string)) {
	aic.submitCallback = callback
}

// SetCloseCallback modified close callback in the container.
func (aic *AnswerInputComponent) SetCloseCallback(callback func()) {
	aic.closeCallback = callback
}

// GetContainer retrieves container widget.
func (aic *AnswerInputComponent) GetContainer() *widget.Container {
	return aic.container
}

// newAnswerInputComponent creates new answer input component.
func newAnswerInputComponent() *AnswerInputComponent {
	var result *AnswerInputComponent

	container := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(
				config.GetWorldWidth(),
				config.GetWorldHeight()),
			widget.WidgetOpts.TrackHover(false),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
			})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	textWidget := widget.NewText(
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
				Padding: widget.Insets{
					Bottom: scaler.GetPercentageOf(config.GetWorldHeight(), 40),
				},
			})),
		widget.TextOpts.Text("", &text.GoTextFace{
			Source: loader.GetInstance().GetFont(loader.KyivRegularFont),
			Size:   40,
		}, componentscommon.ButtonTextColor))

	container.AddChild(textWidget)

	generalFont := &text.GoTextFace{
		Source: loader.GetInstance().GetFont(loader.KyivRegularFont),
		Size:   20,
	}

	var answerInput *widget.TextInput

	answerInput = widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(
				scaler.GetPercentageOf(config.GetWorldWidth(), 45), 0),
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
		widget.TextInputOpts.Placeholder(translation.GetInstance().GetTranslation("client.answerinput.entertext")))

	container.AddChild(answerInput)

	buttonIdleIcon := common.GetImageAsNineSlice(loader.ButtonIdleButton, 16, 15)
	buttonHoverIcon := common.GetImageAsNineSlice(loader.ButtonHoverButton, 16, 15)

	var buttonsLeftPadding int

	switch config.GetSettingsInitialLanguage() {
	case config.SETTINGS_LANGUAGE_ENGLISH:
		buttonsLeftPadding = scaler.GetPercentageOf(config.GetWorldWidth(), 73)

	case config.SETTINGS_LANGUAGE_UKRAINIAN:
		buttonsLeftPadding = scaler.GetPercentageOf(config.GetWorldWidth(), 68)
	}

	buttonsContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(
				container.GetWidget().MinWidth,
				container.GetWidget().MinHeight),
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
			widget.RowLayoutOpts.Padding(widget.Insets{
				Left:   buttonsLeftPadding,
				Bottom: scaler.GetPercentageOf(config.GetWorldHeight(), 9),
			}),
		)))

	buttonsContainer.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:         buttonIdleIcon,
			Hover:        buttonHoverIcon,
			Pressed:      buttonIdleIcon,
			PressedHover: buttonIdleIcon,
			Disabled:     buttonIdleIcon,
		}),
		widget.ButtonOpts.Text(
			translation.GetInstance().GetTranslation("client.answerinput.submit"),
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
			result.submitCallback(answerInput.GetText())
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
			translation.GetInstance().GetTranslation("client.answerinput.close"),
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
			result.closeCallback()
		}),
	))

	container.AddChild(buttonsContainer)

	result = &AnswerInputComponent{
		text:      textWidget,
		container: container,
	}

	return result
}
