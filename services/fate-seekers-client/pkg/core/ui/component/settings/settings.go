package settings

import (
	"fmt"
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

// NewSettingsComponent creates new main settings component.
func NewSettingsComponent(
	submitCallback func(soundMusic, soundFX int, networkingHost, networkingEncryptionKey, language string),
	closeCallback func(soundMusic, soundFX int, networkingHost, networkingEncryptionKey, language string)) *widget.Container {
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
			translation.GetInstance().GetTranslation("shared.settings.title"),
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
		widget.TextOpts.Text(
			translation.GetInstance().GetTranslation("client.settings.sound.music"),
			generalFont,
			color.White)))

	soundMusicComponent := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(100, 10)),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(30),
		)),
	)

	soundMusicLabel := widget.NewText(
		widget.TextOpts.Text(
			fmt.Sprintf("%d", config.GetSettingsSoundMusic()),
			generalFont,
			color.White))

	soundMusicSlider := widget.NewSlider(
		widget.SliderOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(
				scaler.GetPercentageOf(config.GetWorldWidth(), 15),
				10),
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch:  true,
				Position: widget.RowLayoutPositionStart,
			}),
		),
		widget.SliderOpts.MinMax(0, 100),
		widget.SliderOpts.Images(&widget.SliderTrackImage{
			Idle:     image.NewNineSlice(loader.GetInstance().GetStatic(loader.SliderTrackIdle), [3]int{0, 19, 0}, [3]int{6, 0, 0}),
			Hover:    image.NewNineSlice(loader.GetInstance().GetStatic(loader.SliderTrackIdle), [3]int{0, 19, 0}, [3]int{6, 0, 0}),
			Disabled: image.NewNineSlice(loader.GetInstance().GetStatic(loader.SliderTrackIdle), [3]int{0, 19, 0}, [3]int{6, 0, 0}),
		}, &widget.ButtonImage{
			Idle:     image.NewNineSliceSimple(loader.GetInstance().GetStatic(loader.SliderHandleIdle), 0, 5),
			Hover:    image.NewNineSliceSimple(loader.GetInstance().GetStatic(loader.SliderHandleHover), 0, 5),
			Pressed:  image.NewNineSliceSimple(loader.GetInstance().GetStatic(loader.SliderHandleHover), 0, 5),
			Disabled: image.NewNineSliceSimple(loader.GetInstance().GetStatic(loader.SliderHandleIdle), 0, 5),
		}),
		widget.SliderOpts.FixedHandleSize(4),
		widget.SliderOpts.TrackOffset(5),
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			soundMusicLabel.Label = fmt.Sprintf("%d", args.Current)
		}),
	)

	soundMusicSlider.Current = config.GetSettingsSoundMusic()

	soundMusicComponent.AddChild(soundMusicSlider)

	soundMusicComponent.AddChild(soundMusicLabel)

	components.AddChild(soundMusicComponent)

	components.AddChild(widget.NewText(
		widget.TextOpts.Text(
			translation.GetInstance().GetTranslation("client.settings.sound.fx"),
			generalFont,
			color.White)))

	soundFXComponent := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(100, 10)),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(30),
		)),
	)

	soundFXLabel := widget.NewText(
		widget.TextOpts.Text(
			fmt.Sprintf("%d", config.GetSettingsSoundFX()),
			generalFont,
			color.White))

	soundFXSlider := widget.NewSlider(
		widget.SliderOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(
				scaler.GetPercentageOf(config.GetWorldWidth(), 15),
				10),
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch:  true,
				Position: widget.RowLayoutPositionStart,
			}),
		),
		widget.SliderOpts.MinMax(0, 100),
		widget.SliderOpts.Images(&widget.SliderTrackImage{
			Idle:     image.NewNineSlice(loader.GetInstance().GetStatic(loader.SliderTrackIdle), [3]int{0, 19, 0}, [3]int{6, 0, 0}),
			Hover:    image.NewNineSlice(loader.GetInstance().GetStatic(loader.SliderTrackIdle), [3]int{0, 19, 0}, [3]int{6, 0, 0}),
			Disabled: image.NewNineSlice(loader.GetInstance().GetStatic(loader.SliderTrackIdle), [3]int{0, 19, 0}, [3]int{6, 0, 0}),
		}, &widget.ButtonImage{
			Idle:     image.NewNineSliceSimple(loader.GetInstance().GetStatic(loader.SliderHandleIdle), 0, 5),
			Hover:    image.NewNineSliceSimple(loader.GetInstance().GetStatic(loader.SliderHandleHover), 0, 5),
			Pressed:  image.NewNineSliceSimple(loader.GetInstance().GetStatic(loader.SliderHandleHover), 0, 5),
			Disabled: image.NewNineSliceSimple(loader.GetInstance().GetStatic(loader.SliderHandleIdle), 0, 5),
		}),
		widget.SliderOpts.FixedHandleSize(4),
		widget.SliderOpts.TrackOffset(5),
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			soundFXLabel.Label = fmt.Sprintf("%d", args.Current)
		}),
	)

	soundFXSlider.Current = config.GetSettingsSoundFX()

	soundFXComponent.AddChild(soundFXSlider)

	soundFXComponent.AddChild(soundFXLabel)

	components.AddChild(soundFXComponent)

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
			} else if len(parsedNewInputText) == 0 {
				return false, &newInputText
			}

			parsedNewInputTextSymbol := rune(parsedNewInputText[0])

			if parsedNewInputTextSymbol < 32 && parsedNewInputTextSymbol > 127 {
				replacement := networkingHostInput.GetText()

				return false, &replacement
			}

			if len(newInputText) > maxInputSymbols {
				replacement := networkingHostInput.GetText()

				return false, &replacement
			}

			return false, &newInputText
		}))

	networkingHostInput.SetText(config.GetSettingsNetworkingServerHost())

	components.AddChild(networkingHostInput)

	components.AddChild(widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextOpts.Text(
			translation.GetInstance().GetTranslation("shared.settings.encryption-key"),
			generalFont,
			color.White)))

	var networkingEncryptionKeyInput *widget.TextInput

	networkingEncryptionKeyInput = widget.NewTextInput(
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

			parsedNewInputText := newInputText[len(networkingEncryptionKeyInput.GetText()):]

			if len(parsedNewInputText) > 1 {
				newInputText = networkingEncryptionKeyInput.GetText() + parsedNewInputText[:1]
			}

			if len(newInputText) >= maxInputSymbols {
				replacement := networkingEncryptionKeyInput.GetText()

				return false, &replacement
			}

			return false, &newInputText
		}))

	networkingEncryptionKeyInput.SetText(config.GetSettingsNetworkingEncryptionKey())

	components.AddChild(networkingEncryptionKeyInput)

	components.AddChild(widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextOpts.Text(
			translation.GetInstance().GetTranslation("shared.settings.language"),
			generalFont,
			color.White)))

	languageComboButton := widget.NewListComboButton(
		widget.ListComboButtonOpts.SelectComboButtonOpts(
			widget.SelectComboButtonOpts.ComboButtonOpts(
				widget.ComboButtonOpts.ButtonOpts(
					widget.ButtonOpts.Image(&widget.ButtonImage{
						Idle:         common.GetImageAsNineSlice(loader.ComboIdleButton, 12, -10),
						Hover:        common.GetImageAsNineSlice(loader.ComboIdleButton, 12, -10),
						Pressed:      common.GetImageAsNineSlice(loader.ComboIdleButton, 12, -10),
						PressedHover: common.GetImageAsNineSlice(loader.ComboIdleButton, 12, -10),
					}),
				),
			),
		),
		widget.ListComboButtonOpts.Text(generalFont, &widget.ButtonImageImage{
			Idle:     loader.GetInstance().GetStatic(loader.ComboArrayIdleButton),
			Disabled: loader.GetInstance().GetStatic(loader.ComboArrayIdleButton),
		}, &widget.ButtonTextColor{
			Idle:     componentscommon.ButtonTextColor,
			Disabled: componentscommon.ButtonTextColor,
			Hover:    componentscommon.ButtonTextColor,
			Pressed:  componentscommon.ButtonTextColor,
		}),
		widget.ListComboButtonOpts.ListOpts(
			widget.ListOpts.Entries([]interface{}{
				config.SETTINGS_LANGUAGE_ENGLISH,
				config.SETTINGS_LANGUAGE_UKRAINIAN,
			}),
			widget.ListOpts.ScrollContainerOpts(
				widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
					Idle:     image.NewNineSlice(loader.GetInstance().GetStatic(loader.ListIdle), [3]int{25, 12, 22}, [3]int{25, 12, 25}),
					Disabled: image.NewNineSlice(loader.GetInstance().GetStatic(loader.ListDisabled), [3]int{25, 12, 22}, [3]int{25, 12, 25}),
					Mask:     image.NewNineSlice(loader.GetInstance().GetStatic(loader.ListMask), [3]int{26, 10, 23}, [3]int{26, 10, 26}),
				}),
			),
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
				widget.SliderOpts.MinHandleSize(4),
				widget.SliderOpts.TrackPadding(widget.Insets{Bottom: 20})),
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
		),
		widget.ListComboButtonOpts.EntryLabelFunc(
			func(e any) string {
				switch e.(string) {
				case config.SETTINGS_LANGUAGE_ENGLISH:
					return translation.GetInstance().GetTranslation("shared.settings.language.english")
				case config.SETTINGS_LANGUAGE_UKRAINIAN:
					return translation.GetInstance().GetTranslation("shared.settings.language.ukrainian")
				default:
					return e.(string)
				}
			},
			func(e any) string {
				switch e.(string) {
				case config.SETTINGS_LANGUAGE_ENGLISH:
					return translation.GetInstance().GetTranslation("shared.settings.language.english")
				case config.SETTINGS_LANGUAGE_UKRAINIAN:
					return translation.GetInstance().GetTranslation("shared.settings.language.ukrainian")
				default:
					return e.(string)
				}
			}))

	languageComboButton.SetSelectedEntry(config.GetSettingsLanguage())

	components.AddChild(languageComboButton)

	result.AddChild(components)

	var buttonsLeftPadding int

	switch config.GetSettingsInitialLanguage() {
	case config.SETTINGS_LANGUAGE_ENGLISH:
		buttonsLeftPadding = scaler.GetPercentageOf(config.GetWorldWidth(), 20)

	case config.SETTINGS_LANGUAGE_UKRAINIAN:
		buttonsLeftPadding = scaler.GetPercentageOf(config.GetWorldWidth(), 15)
	}

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
			widget.RowLayoutOpts.Padding(widget.Insets{
				Left: buttonsLeftPadding,
			}),
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
			translation.GetInstance().GetTranslation("shared.settings.close"),
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
			closeCallback(
				soundMusicSlider.Current,
				soundFXSlider.Current,
				networkingHostInput.GetText(),
				networkingEncryptionKeyInput.GetText(),
				languageComboButton.SelectedEntry().(string))
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
			translation.GetInstance().GetTranslation("shared.settings.submit"),
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
			submitCallback(
				soundMusicSlider.Current,
				soundFXSlider.Current,
				networkingHostInput.GetText(),
				networkingEncryptionKeyInput.GetText(),
				languageComboButton.SelectedEntry().(string))
		}),
	))

	result.AddChild(buttonsContainer)

	return result
}
