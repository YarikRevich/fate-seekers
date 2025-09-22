package prompt

import (
	"image/color"
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/sound"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/common"
	componentscommon "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/translation"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	// GetInstance retrieves instance of the prompt component, performing initial creation if needed.
	GetInstance = sync.OnceValue[*PromptComponent](newPromptComponent)
)

// PromptComponent represents component, which contains prompt statement.
type PromptComponent struct {
	// Represents text widget.
	text *widget.Text

	// Represents submit callback.
	submitCallback func()

	// Represents close callback.
	closeCallback func()

	// Represents container widget.
	container *widget.Container
}

// SetText modifies text component in the container.
func (pc *PromptComponent) SetText(value string) {
	pc.text.Label = value
}

// SetSubmitCallback modified submit callback in the container.
func (pc *PromptComponent) SetSubmitCallback(callback func()) {
	pc.submitCallback = callback
}

// SetCloseCallback modified close callback in the container.
func (pc *PromptComponent) SetCloseCallback(callback func()) {
	pc.closeCallback = callback
}

// GetContainer retrieves container widget.
func (pc *PromptComponent) GetContainer() *widget.Container {
	return pc.container
}

// newPromptComponent creates prompt component.
func newPromptComponent() *PromptComponent {
	var result *PromptComponent

	container := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(
				scaler.GetPercentageOf(config.GetWorldWidth(), 40),
				scaler.GetPercentageOf(config.GetWorldHeight(), 20)),
			widget.WidgetOpts.TrackHover(false),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
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
			Size:   25,
		}, color.White))

	container.AddChild(textWidget)

	generalFont := &text.GoTextFace{
		Source: loader.GetInstance().GetFont(loader.KyivRegularFont),
		Size:   20,
	}

	buttonIdleIcon := common.GetImageAsNineSlice(loader.ButtonIdleButton, 16, 15)
	buttonHoverIcon := common.GetImageAsNineSlice(loader.ButtonHoverButton, 16, 15)

	var buttonsLeftPadding int

	switch config.GetSettingsInitialLanguage() {
	case config.SETTINGS_LANGUAGE_ENGLISH:
		buttonsLeftPadding = scaler.GetPercentageOf(config.GetWorldWidth(), 34)

	case config.SETTINGS_LANGUAGE_UKRAINIAN:
		buttonsLeftPadding = scaler.GetPercentageOf(config.GetWorldWidth(), 20)
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
				Left: buttonsLeftPadding,
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
			translation.GetInstance().GetTranslation("shared.prompt.submit"),
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
			sound.GetInstance().GetSoundFxManager().PushWithHandbrake(loader.ButtonFXSound)

			result.submitCallback()
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
			translation.GetInstance().GetTranslation("shared.prompt.close"),
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
			sound.GetInstance().GetSoundFxManager().PushWithHandbrake(loader.ButtonFXSound)

			result.closeCallback()
		}),
	))

	container.AddChild(buttonsContainer)

	result = &PromptComponent{
		text:      textWidget,
		container: container,
	}

	return result
}
