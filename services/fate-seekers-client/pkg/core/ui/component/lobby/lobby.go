package lobby

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

// Describes all the colors used for list combo definition.
var (
	selectedListColor = color.NRGBA{183, 228, 202, 255}
	focusedListColor  = color.NRGBA{R: 170, G: 170, B: 180, A: 255}
	disabledListColor = color.NRGBA{100, 100, 100, 255}
)

var (
	// GetInstance retrieves instance of the lobby component, performing initial creation if needed.
	GetInstance = sync.OnceValue[*LobbyComponent](newLobbyComponent)
)

// LobbyComponent represents component, which contains lobby menu.
type LobbyComponent struct {
	// Represents session name input widget.
	sessionNameInput *widget.TextInput

	// Represents sessions list widget.
	list *widget.List

	// Represents currently selected session name entry.
	sessionNameEntry string

	// Represents start action button widget.
	startActionButton *widget.Button

	// Represents start callback.
	startCallback func()

	// Represents back callback.
	backCallback func()

	// Represents container widget.
	container *widget.Container
}

// SetListsEntries sets lists entries to the list widget.
func (lc *LobbyComponent) SetListsEntries(value []interface{}) {
	fmt.Println(value)

	lc.list.SetEntries(value)
}

// SetStartCallback modifies start callback in the container.
func (lc *LobbyComponent) SetStartCallback(callback func()) {
	lc.startCallback = callback
}

// SetBackCallback modifies back callback in the container.
func (lc *LobbyComponent) SetBackCallback(callback func()) {
	lc.backCallback = callback
}

// ShowStartButton shows start button widget.
func (lc *LobbyComponent) ShowStartButton() {
	lc.startActionButton.GetWidget().Visibility = widget.Visibility_Show
}

// HideStartButton hides start button widget.
func (lc *LobbyComponent) HideStartButton() {
	lc.startActionButton.GetWidget().Visibility = widget.Visibility_Hide
}

// GetContainer retrieves container widget.
func (lc *LobbyComponent) GetContainer() *widget.Container {
	return lc.container
}

// newLobbyComponent creates new selector component.
func newLobbyComponent() *LobbyComponent {
	var result *LobbyComponent

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
			translation.GetInstance().GetTranslation("client.lobby.title"),
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
			translation.GetInstance().GetTranslation("client.lobby.players"),
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
			return translation.
				GetInstance().
				GetTranslation(
					fmt.Sprintf("client.skin.%d.name", e.(uint64)))
		}),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			// sessionNameEntry := args.Entry.(string)

			// result.sessionNameEntry = sessionNameEntry
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
			translation.GetInstance().GetTranslation("client.lobby.close"),
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
				Left: scaler.GetPercentageOf(config.GetWorldWidth(), 11),
			}),
		)),
	)

	startActionButton := widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:         buttonIdleIcon,
			Hover:        buttonHoverIcon,
			Pressed:      buttonIdleIcon,
			PressedHover: buttonIdleIcon,
			Disabled:     buttonIdleIcon,
		}),
		widget.ButtonOpts.Text(
			translation.GetInstance().GetTranslation("client.lobby.start"),
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
			result.startCallback()
		}),
	)

	actionButtonContainer.AddChild(startActionButton)

	buttonsContainer.AddChild(actionButtonContainer)

	container.AddChild(buttonsContainer)

	result = &LobbyComponent{
		list:              list,
		startActionButton: startActionButton,
		container:         container,
	}

	return result
}
