package menu

import (
	"image/color"
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/sound"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/loader"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/ui/common"
	componentscommon "github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/ui/component/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/ui/manager/translation"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	// GetInstance retrieves instance of the menu component, performing initial creation if needed.
	GetInstance = sync.OnceValue[*MenuComponent](newMenuComponent)
)

// MenuComponent represents component, which contains menu statement.
type MenuComponent struct {
	// Represents start button widget.
	startButtonWidget *widget.Button

	// Represents stop button widget.
	stopButtonWidget *widget.Button

	// Represents start callback.
	startCallback func()

	// Represents stop callback.
	stopCallback func()

	// Represents monitoring callback.
	monitoringCallback func()

	// Represents settings callback.
	settingsCallback func()

	// Represents exit callback.
	exitCallback func()

	// Represents container widget.
	container *widget.Container
}

// EnableStartButton enables start button component in the container.
func (mc *MenuComponent) EnableStartButton() {
	mc.startButtonWidget.GetWidget().Disabled = false
}

// DisableStartButton disables start button component in the container.
func (mc *MenuComponent) DisableStartButton() {
	mc.startButtonWidget.GetWidget().Disabled = true
}

// EnableStopButton enables stop button component in the container.
func (mc *MenuComponent) EnableStopButton() {
	mc.stopButtonWidget.GetWidget().Disabled = false
}

// DisableStopButton disables stop button component in the container.
func (mc *MenuComponent) DisableStopButton() {
	mc.stopButtonWidget.GetWidget().Disabled = true
}

// SetStartCallback modified start callback in the container.
func (mc *MenuComponent) SetStartCallback(callback func()) {
	mc.startCallback = callback
}

// SetStopCallback modified stop callback in the container.
func (mc *MenuComponent) SetStopCallback(callback func()) {
	mc.stopCallback = callback
}

// SetMonitoringCallback modified monitoring callback in the container.
func (mc *MenuComponent) SetMonitoringCallback(callback func()) {
	mc.monitoringCallback = callback
}

// SetSettingsCallback modified settings callback in the container.
func (mc *MenuComponent) SetSettingsCallback(callback func()) {
	mc.settingsCallback = callback
}

// SetExitCallback modified exit callback in the container.
func (mc *MenuComponent) SetExitCallback(callback func()) {
	mc.exitCallback = callback
}

// GetContainer retrieves container widget.
func (mc *MenuComponent) GetContainer() *widget.Container {
	return mc.container
}

// newMenuComponent creates menu component.
func newMenuComponent() *MenuComponent {
	var result *MenuComponent

	container := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(
				scaler.GetPercentageOf(config.GetWorldWidth(), 20),
				scaler.GetPercentageOf(config.GetWorldHeight(), 40),
			),
			widget.WidgetOpts.TrackHover(false),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				Padding: widget.Insets{
					Left: scaler.GetPercentageOf(config.GetWorldWidth(), 17),
				},
				VerticalPosition:  widget.AnchorLayoutPositionCenter,
				StretchHorizontal: false,
				StretchVertical:   false,
			})),
		widget.ContainerOpts.BackgroundImage(common.GetImageAsNineSlice(loader.PanelIdlePanel, 10, 10)),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Left:   100,
				Right:  100,
				Top:    40,
				Bottom: 40,
			}),
		)))

	buttonsContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: widget.RowLayoutPositionCenter,
		})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(30),
		)))

	buttonIdleIcon := common.GetImageAsNineSlice(loader.ButtonIdleButton, 16, 15)
	buttonHoverIcon := common.GetImageAsNineSlice(loader.ButtonHoverButton, 16, 15)

	buttonFont := &text.GoTextFace{
		Source: loader.GetInstance().GetFont(loader.KyivRegularFont),
		Size:   20,
	}

	startButtonTooltipContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(
			image.NewNineSlice(loader.GetInstance().GetStatic(loader.ToolTip), [3]int{19, 6, 13}, [3]int{19, 5, 13})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Left:   15,
				Right:  15,
				Top:    10,
				Bottom: 10,
			}),
			widget.RowLayoutOpts.Spacing(2),
		)))

	generalFont := &text.GoTextFace{
		Source: loader.GetInstance().GetFont(loader.KyivRegularFont),
		Size:   20,
	}

	startButtonTooltipContainer.AddChild(widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextOpts.Text(
			translation.GetInstance().GetTranslation("server.menu.start.disabled"),
			generalFont,
			color.White)))

	var startButtonWidget *widget.Button

	startButtonWidget = widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
			widget.WidgetOpts.ToolTip(widget.NewToolTip(
				widget.ToolTipOpts.Content(startButtonTooltipContainer),
				widget.ToolTipOpts.ToolTipUpdater(func(c *widget.Container) {
					if !startButtonWidget.GetWidget().Disabled &&
						c.GetWidget().Visibility != widget.Visibility_Hide {
						c.GetWidget().Visibility = widget.Visibility_Hide
					} else if startButtonWidget.GetWidget().Disabled && c.GetWidget().Visibility != widget.Visibility_Show {
						c.GetWidget().Visibility = widget.Visibility_Show
					}
				}),
			)),
		),
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:         buttonIdleIcon,
			Hover:        buttonHoverIcon,
			Pressed:      buttonIdleIcon,
			PressedHover: buttonIdleIcon,
			Disabled:     buttonIdleIcon,
		}),
		widget.ButtonOpts.Text(
			translation.GetInstance().GetTranslation("shared.menu.start"),
			buttonFont,
			&widget.ButtonTextColor{Idle: componentscommon.ButtonTextColor}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    20,
			Bottom: 20,
		}),
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			sound.GetInstance().GetSoundFxManager().PushWithHandbrake(loader.ButtonFXSound)

			result.startCallback()
		}),
	)

	buttonsContainer.AddChild(startButtonWidget)

	stopButtonTooltipContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(
			image.NewNineSlice(loader.GetInstance().GetStatic(loader.ToolTip), [3]int{19, 6, 13}, [3]int{19, 5, 13})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Left:   15,
				Right:  15,
				Top:    10,
				Bottom: 10,
			}),
			widget.RowLayoutOpts.Spacing(2),
		)))

	stopButtonTooltipContainer.AddChild(widget.NewText(
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextOpts.Text(
			translation.GetInstance().GetTranslation("server.menu.stop.disabled"),
			generalFont,
			color.White)))

	var stopButtonWidget *widget.Button

	stopButtonWidget = widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
			widget.WidgetOpts.ToolTip(widget.NewToolTip(
				widget.ToolTipOpts.Content(stopButtonTooltipContainer),
				widget.ToolTipOpts.ToolTipUpdater(func(c *widget.Container) {
					if !stopButtonWidget.GetWidget().Disabled &&
						c.GetWidget().Visibility != widget.Visibility_Hide {
						c.GetWidget().Visibility = widget.Visibility_Hide
					} else if stopButtonWidget.GetWidget().Disabled && c.GetWidget().Visibility != widget.Visibility_Show {
						c.GetWidget().Visibility = widget.Visibility_Show
					}
				}),
			)),
		),
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:         buttonIdleIcon,
			Hover:        buttonHoverIcon,
			Pressed:      buttonIdleIcon,
			PressedHover: buttonIdleIcon,
			Disabled:     buttonIdleIcon,
		}),
		widget.ButtonOpts.Text(
			translation.GetInstance().GetTranslation("server.menu.stop"),
			buttonFont,
			&widget.ButtonTextColor{Idle: componentscommon.ButtonTextColor}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    20,
			Bottom: 20,
		}),
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			sound.GetInstance().GetSoundFxManager().PushWithHandbrake(loader.ButtonFXSound)

			result.stopCallback()
		}),
	)

	buttonsContainer.AddChild(stopButtonWidget)

	buttonsContainer.AddChild(widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:         buttonIdleIcon,
			Hover:        buttonHoverIcon,
			Pressed:      buttonIdleIcon,
			PressedHover: buttonIdleIcon,
			Disabled:     buttonIdleIcon,
		}),
		widget.ButtonOpts.Text(
			translation.GetInstance().GetTranslation("server.menu.monitoring"),
			buttonFont,
			&widget.ButtonTextColor{Idle: componentscommon.ButtonTextColor}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    20,
			Bottom: 20,
		}),
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			sound.GetInstance().GetSoundFxManager().PushWithHandbrake(loader.ButtonFXSound)

			result.monitoringCallback()
		}),
	))

	buttonsContainer.AddChild(widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:         buttonIdleIcon,
			Hover:        buttonHoverIcon,
			Pressed:      buttonIdleIcon,
			PressedHover: buttonIdleIcon,
			Disabled:     buttonIdleIcon,
		}),
		widget.ButtonOpts.Text(
			translation.GetInstance().GetTranslation("shared.menu.settings"),
			buttonFont,
			&widget.ButtonTextColor{Idle: componentscommon.ButtonTextColor}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    20,
			Bottom: 20,
		}),
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			sound.GetInstance().GetSoundFxManager().PushWithHandbrake(loader.ButtonFXSound)

			result.settingsCallback()
		}),
	))

	buttonsContainer.AddChild(widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:         buttonIdleIcon,
			Hover:        buttonHoverIcon,
			Pressed:      buttonIdleIcon,
			PressedHover: buttonIdleIcon,
			Disabled:     buttonIdleIcon,
		}),
		widget.ButtonOpts.Text(
			translation.GetInstance().GetTranslation("shared.menu.exit"),
			buttonFont,
			&widget.ButtonTextColor{Idle: componentscommon.ButtonTextColor}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    20,
			Bottom: 20,
		}),
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			sound.GetInstance().GetSoundFxManager().PushWithHandbrake(loader.ButtonFXSound)

			result.exitCallback()
		}),
	))

	container.AddChild(buttonsContainer)

	result = &MenuComponent{
		startButtonWidget: startButtonWidget,
		stopButtonWidget:  stopButtonWidget,
		container:         container,
	}

	result.DisableStopButton()

	return result
}
