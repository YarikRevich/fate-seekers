package selector

import (
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
	maxInputSymbols = 8
)

// Describes all the colors used for list combo definition.
var (
	selectedListColor = color.NRGBA{183, 228, 202, 255}
	focusedListColor  = color.NRGBA{R: 170, G: 170, B: 180, A: 255}
	disabledListColor = color.NRGBA{100, 100, 100, 255}
)

var (
	// GetInstance retrieves instance of the selector component, performing initial creation if needed.
	GetInstance = sync.OnceValue[*SelectorComponent](newSelectorComponent)
)

// SelectorComponent represents component, which contains selector menu.
type SelectorComponent struct {
	// Represents name input widget.
	sessionIDInput *widget.TextInput

	// Represents list widget.
	list *widget.List

	// Represents submit callback.
	submitCallback func(sessionID string)

	// Represents create callback.
	createCallback func()

	// Represents back callback.
	backCallback func()

	// Represents container widget.
	container *widget.Container
}

// CleanInputs cleans all the inputs in the container.
func (sc *SelectorComponent) CleanInputs() {
	sc.sessionIDInput.SetText("")
}

// SetListsEntries sets lists entries to the list widget.
func (sc *SelectorComponent) SetListsEntries(value []interface{}) {
	sc.list.SetEntries(value)
}

// SetSubmitCallback modified submit callback in the container.
func (sc *SelectorComponent) SetSubmitCallback(callback func(sessionID string)) {
	sc.submitCallback = callback
}

// SetCreateCallback modified create callback in the container.
func (sc *SelectorComponent) SetCreateCallback(callback func()) {
	sc.createCallback = callback
}

// SetBackCallback modified back callback in the container.
func (sc *SelectorComponent) SetBackCallback(callback func()) {
	sc.backCallback = callback
}

// GetContainer retrieves container widget.
func (sc *SelectorComponent) GetContainer() *widget.Container {
	return sc.container
}

// newSelectorComponent creates new selector component.
func newSelectorComponent() *SelectorComponent {
	var result *SelectorComponent

	container := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(
				scaler.GetPercentageOf(config.GetWorldWidth(), 20),
				scaler.GetPercentageOf(config.GetWorldHeight(), 30)),
			widget.WidgetOpts.TrackHover(false),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				Padding: widget.Insets{
					Left: scaler.GetPercentageOf(config.GetWorldWidth(), 6),
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

	container.AddChild(widget.NewText(
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
			translation.GetInstance().GetTranslation("client.selector.session_id"),
			generalFont,
			color.White)))

	var sessionIDInput *widget.TextInput

	sessionIDInput = widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
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

			parsedNewInputText := newInputText[len(sessionIDInput.GetText()):]

			if len(parsedNewInputText) > 1 {
				newInputText = sessionIDInput.GetText() + parsedNewInputText[:1]
			} else if len(parsedNewInputText) == 0 {
				return false, &newInputText
			}

			parsedNewInputTextSymbol := rune(parsedNewInputText[0])

			if parsedNewInputTextSymbol < 32 && parsedNewInputTextSymbol > 127 {
				replacement := sessionIDInput.GetText()

				return false, &replacement
			}

			if len(newInputText) > maxInputSymbols {
				replacement := sessionIDInput.GetText()

				return false, &replacement
			}

			return false, &newInputText
		}))

	components.AddChild(sessionIDInput)

	container.AddChild(components)

	listsContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch:  true,
				Position: widget.RowLayoutPositionCenter,
			}),
			widget.WidgetOpts.MinSize(
				container.GetWidget().MinWidth,
				container.GetWidget().MinHeight,
			),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top: 40,
			}),
		)))

	listsContainer.AddChild(widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextOpts.Insets(widget.Insets{
			Bottom: 20,
		}),
		widget.TextOpts.Text(
			translation.GetInstance().GetTranslation("client.selector.sessions"),
			generalFont,
			color.White)))

	list := widget.NewList(
		widget.ListOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.MinSize(
					scaler.GetPercentageOf(config.GetWorldWidth(), 40),
					scaler.GetPercentageOf(config.GetWorldHeight(), 30),
				),
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					MaxWidth:  scaler.GetPercentageOf(config.GetWorldWidth(), 40),
					MaxHeight: scaler.GetPercentageOf(config.GetWorldHeight(), 30),
					Position:  widget.RowLayoutPositionCenter,
				}))),
		widget.ListOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
			Idle:     image.NewNineSlice(loader.GetInstance().GetStatic(loader.ListIdle), [3]int{25, 12, 22}, [3]int{25, 12, 25}),
			Disabled: image.NewNineSlice(loader.GetInstance().GetStatic(loader.ListDisabled), [3]int{25, 12, 22}, [3]int{25, 12, 25}),
			Mask:     image.NewNineSlice(loader.GetInstance().GetStatic(loader.ListMask), [3]int{26, 10, 23}, [3]int{26, 10, 26}),
		})),
		widget.ListOpts.SliderOpts(
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
			widget.SliderOpts.MinHandleSize(8),
			widget.SliderOpts.TrackPadding(widget.Insets{Bottom: 20}),
		),
		widget.ListOpts.AllowReselect(),
		widget.ListOpts.HideHorizontalSlider(),
		widget.ListOpts.Entries([]interface{}{}),
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			return e.(string)
		}),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			sessionIDInput.SetText(args.Entry.(string))
		}),
		widget.ListOpts.EntryFontFace(generalFont),
		widget.ListOpts.EntryColor(&widget.ListEntryColor{
			Selected:                   componentscommon.ButtonTextColor,
			Unselected:                 selectedListColor,
			SelectedBackground:         selectedListColor,
			SelectedFocusedBackground:  selectedListColor,
			FocusedBackground:          focusedListColor,
			DisabledUnselected:         disabledListColor,
			DisabledSelected:           disabledListColor,
			DisabledSelectedBackground: disabledListColor,
		}),
		widget.ListOpts.EntryTextPadding(widget.Insets{
			Top:    15,
			Left:   40,
			Right:  40,
			Bottom: 15,
		}),
	)

	listsContainer.AddChild(list)

	container.AddChild(listsContainer)

	buttonsContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(
				container.GetWidget().MinWidth,
				0),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top: -scaler.GetPercentageOf(config.GetWorldHeight(), 10),
			}),
		)),
	)

	closeButtonsContainer := widget.NewContainer(
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
		)),
	)

	buttonIdleIcon := common.GetImageAsNineSlice(loader.ButtonIdleButton, 16, 15)
	buttonHoverIcon := common.GetImageAsNineSlice(loader.ButtonHoverButton, 16, 15)

	closeButtonsContainer.AddChild(widget.NewButton(
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
			result.backCallback()
		}),
	))

	buttonsContainer.AddChild(closeButtonsContainer)

	actionButtonContainer := widget.NewContainer(
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
				Left: scaler.GetPercentageOf(config.GetWorldWidth(), 4),
			}),
		)),
	)

	actionButtonContainer.AddChild(widget.NewButton(
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
			result.createCallback()
		}),
	))

	actionButtonContainer.AddChild(widget.NewButton(
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
			result.submitCallback(sessionIDInput.GetText())
		}),
	))

	buttonsContainer.AddChild(actionButtonContainer)

	container.AddChild(buttonsContainer)

	result = &SelectorComponent{
		sessionIDInput: sessionIDInput,
		list:           list,
		container:      container,
	}

	return result
}
