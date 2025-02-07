package letter

import (
	"image/color"
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/common"
	componentscommon "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	// GetInstance retrieves instance of the letter component, performing initial creation if needed.
	GetInstance = sync.OnceValue[*LetterComponent](newLetterComponent)
)

// LetterComponent represents component, which contains actual letter.
type LetterComponent struct {
	// Represents text area widget.
	textArea *widget.TextArea

	// Represents attachment value.
	attachmentValue *string

	// Represents attachment callback.
	attachmentCallback func(value string)

	// Represents close callback.
	closeCallback func()

	// Represents container widget.
	container *widget.Container
}

// SetText modifies text component in the container.
func (lc *LetterComponent) SetText(value string) {
	lc.textArea.SetText(value)
}

// GetText retrieves current text.
func (lc *LetterComponent) GetText() string {
	return lc.textArea.GetText()
}

// SetAttachment modified attachment button redirect in the container.
func (lc *LetterComponent) SetAttachment(value string) {
	*lc.attachmentValue = value
}

// GetAttachment retrieves attachment button redirect.
func (lc *LetterComponent) GetAttachment(value string) {
	*lc.attachmentValue = value
}

// SetAttachmentCallback modified close callback in the container.
func (lc *LetterComponent) SetAttachmentCallback(callback func(value string)) {
	lc.attachmentCallback = callback
}

// SetCloseCallback modified close callback in the container.
func (lc *LetterComponent) SetCloseCallback(callback func()) {
	lc.closeCallback = callback
}

// GetContainer retrieves container widget.
func (lc *LetterComponent) GetContainer() *widget.Container {
	return lc.container
}

// newLetterComponent initializes LetterComponent.
func newLetterComponent() *LetterComponent {
	var result *LetterComponent

	container := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(
			scaler.GetPercentageOf(config.GetWorldWidth(), 45),
			scaler.GetPercentageOf(config.GetWorldHeight(), 60))),
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

	container.AddChild(widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextOpts.Text("Letter", generalFont, color.White)))

	textAreaContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
			widget.WidgetOpts.MinSize(
				container.GetWidget().MinWidth,
				scaler.GetPercentageOf(container.GetWidget().MinHeight, 72)),
		),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true}),
			widget.GridLayoutOpts.Spacing(0, 0)),
		),
	)

	textArea := widget.NewTextArea(
		widget.TextAreaOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.MinSize(
					textAreaContainer.GetWidget().MinWidth,
					textAreaContainer.GetWidget().MinHeight),
				widget.WidgetOpts.LayoutData(widget.GridLayoutData{
					MaxWidth:  textAreaContainer.GetWidget().MinWidth,
					MaxHeight: textAreaContainer.GetWidget().MinHeight,
				}))),
		widget.TextAreaOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
			Idle:     image.NewNineSlice(loader.GetInstance().GetStatic(loader.ListIdle), [3]int{25, 12, 22}, [3]int{25, 12, 25}),
			Disabled: image.NewNineSlice(loader.GetInstance().GetStatic(loader.ListDisabled), [3]int{25, 12, 22}, [3]int{25, 12, 25}),
			Mask:     image.NewNineSlice(loader.GetInstance().GetStatic(loader.ListMask), [3]int{26, 10, 23}, [3]int{26, 10, 26}),
		})),
		widget.TextAreaOpts.SliderOpts(
			widget.SliderOpts.Images(
				&widget.SliderTrackImage{
					Idle:     image.NewNineSlice(loader.GetInstance().GetStatic(loader.ListTrackIdle), [3]int{5, 0, 0}, [3]int{25, 12, 25}),
					Hover:    image.NewNineSlice(loader.GetInstance().GetStatic(loader.ListTrackIdle), [3]int{5, 0, 0}, [3]int{25, 12, 25}),
					Disabled: image.NewNineSlice(loader.GetInstance().GetStatic(loader.ListTrackDisabled), [3]int{0, 5, 0}, [3]int{25, 12, 25}),
				},
				&widget.ButtonImage{
					Idle:     image.NewNineSliceSimple(loader.GetInstance().GetStatic(loader.SliderHandleIdle), 0, 5),
					Hover:    image.NewNineSliceSimple(loader.GetInstance().GetStatic(loader.SliderHandleHover), 0, 5),
					Pressed:  image.NewNineSliceSimple(loader.GetInstance().GetStatic(loader.SliderHandleHover), 0, 5),
					Disabled: image.NewNineSliceSimple(loader.GetInstance().GetStatic(loader.SliderHandleIdle), 0, 5),
				}),
			widget.SliderOpts.MinHandleSize(5),
			widget.SliderOpts.TrackPadding(widget.Insets{
				Top:    5,
				Bottom: 24,
			}),
		),
		widget.TextAreaOpts.ShowVerticalScrollbar(),
		widget.TextAreaOpts.VerticalScrollMode(widget.ScrollBeginning),
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
		widget.TextAreaOpts.Text(""),
	)

	textAreaContainer.AddChild(textArea)

	container.AddChild(textAreaContainer)

	bottomContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(
				container.GetWidget().MinWidth,
				scaler.GetPercentageOf(container.GetWidget().MinHeight, 3)),
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionEnd,
				Stretch:  true,
			})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	buttonIdleIcon := common.GetImageAsNineSlice(loader.ButtonIdleButton, 16, 15)
	buttonHoverIcon := common.GetImageAsNineSlice(loader.ButtonHoverButton, 16, 15)

	buttonsContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(
				bottomContainer.GetWidget().MinWidth,
				bottomContainer.GetWidget().MinHeight),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				StretchHorizontal:  true,
				StretchVertical:    false,
			}),
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionEnd,
				Stretch:  true,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(30),
		)))

	attachmentValue := new(string)

	*attachmentValue = value.LETTER_IMAGE_EMPTY_VALUE

	buttonsContainer.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:         buttonIdleIcon,
			Hover:        buttonHoverIcon,
			Pressed:      buttonIdleIcon,
			PressedHover: buttonIdleIcon,
			Disabled:     buttonIdleIcon,
		}),
		widget.ButtonOpts.Text("Attachment", generalFont, &widget.ButtonTextColor{Idle: componentscommon.ButtonTextColor}),
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			result.attachmentCallback(*result.attachmentValue)
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   25,
			Right:  25,
			Top:    25,
			Bottom: 25,
		}),
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
	))

	buttonsContainer.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:         buttonIdleIcon,
			Hover:        buttonHoverIcon,
			Pressed:      buttonIdleIcon,
			PressedHover: buttonIdleIcon,
			Disabled:     buttonIdleIcon,
		}),
		widget.ButtonOpts.Text("Close", generalFont, &widget.ButtonTextColor{Idle: componentscommon.ButtonTextColor}),
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			result.closeCallback()
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   25,
			Right:  25,
			Top:    25,
			Bottom: 25,
		}),
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
	))

	bottomContainer.AddChild(buttonsContainer)

	container.AddChild(bottomContainer)

	result = &LetterComponent{
		textArea:        textArea,
		attachmentValue: attachmentValue,
		container:       container,
	}

	return result
}
