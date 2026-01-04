package collections

import (
	"image/color"
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/sound"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/common"
	componentscommon "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/translation"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/entity"
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
	// GetInstance retrieves instance of the collections component, performing initial creation if needed.
	GetInstance = sync.OnceValue[*CollectionsComponent](newCollectionsComponent)
)

// CollectionsComponent represents component which displays gathered collections.
type CollectionsComponent struct {
	// Represents collections list widget.
	list *widget.List

	// Represents the callback fired when a list item is clicked.
	entrySelectedCallback func(path string)

	// Represents back callback.
	backCallback func()

	// Represents container widget.
	container *widget.Container
}

// SetListsEntries sets lists entries to the list widget.
func (cc *CollectionsComponent) SetListsEntries(value []interface{}) {
	cc.list.SetEntries(value)
}

// SetEntrySelectedCallback modifies the callback executed when a list item is clicked.
func (cc *CollectionsComponent) SetEntrySelectedCallback(callback func(path string)) {
	cc.entrySelectedCallback = callback
}

// SetBackCallback modifies back callback in the container.
func (cc *CollectionsComponent) SetBackCallback(callback func()) {
	cc.backCallback = callback
}

// GetContainer retrieves container widget.
func (cc *CollectionsComponent) GetContainer() *widget.Container {
	return cc.container
}

// newCollectionsComponent creates new collections component.
func newCollectionsComponent() *CollectionsComponent {
	var result *CollectionsComponent

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
			translation.GetInstance().GetTranslation("client.collections.title"),
			generalFont,
			color.White)))

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
				Top: scaler.GetPercentageOf(config.GetWorldHeight(), 6),
			}),
		)))

	listsContainer.AddChild(widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextOpts.Insets(widget.Insets{
			Bottom: 10,
		}),
		widget.TextOpts.Text(
			translation.GetInstance().GetTranslation("client.collections.available-collections"),
			generalFont,
			color.White)))

	list := widget.NewList(
		widget.ListOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.MinSize(
					scaler.GetPercentageOf(config.GetWorldWidth(), 40),
					scaler.GetPercentageOf(config.GetWorldHeight(), 40),
				),
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					MaxWidth:  scaler.GetPercentageOf(config.GetWorldWidth(), 40),
					MaxHeight: scaler.GetPercentageOf(config.GetWorldHeight(), 40),
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
			return (e.(entity.CollectionEntity)).Name
		}),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			if result.entrySelectedCallback != nil {
				result.entrySelectedCallback((args.Entry.(entity.CollectionEntity)).Path)
			}
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

	closeButtonsContainer := widget.NewContainer(
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
				Top: scaler.GetPercentageOf(config.GetWorldHeight(), 5),
			}),
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
			translation.GetInstance().GetTranslation("client.collections.close"),
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
			sound.GetInstance().GetSoundUIFxManager().PushWithHandbrake(loader.ButtonFXSound)

			if result.backCallback != nil {
				result.backCallback()
			}
		}),
	))

	container.AddChild(closeButtonsContainer)

	result = &CollectionsComponent{
		list:      list,
		container: container,
	}

	return result
}
