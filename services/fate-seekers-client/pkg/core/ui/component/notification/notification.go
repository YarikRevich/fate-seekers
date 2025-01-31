package notification

import (
	"image/color"
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	// GetInstance retrieves instance of the notification component, performing initial creation if needed.
	GetInstance = sync.OnceValue[*NotificationComponent](newNotificationComponent)
)

// NotificationComponent represents component, which contains actual notification.
type NotificationComponent struct {
	container *widget.Container
}

// SetText modifies text component in the container.
func (sc *NotificationComponent) SetText(value string) {
	sc.container.GetWidget().Visibility = widget.Visibility_Show

	sc.container.AddChild(widget.NewText(
		widget.TextOpts.MaxWidth(float64(scaler.GetPercentageOf(config.GetWorldWidth(), 30))),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionStart,
			StretchHorizontal:  false,
			StretchVertical:    false,
		})),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionStart),
		widget.TextOpts.Insets(widget.Insets{
			Top:    20,
			Bottom: 20,
			Left:   20,
			Right:  20,
		}),
		widget.TextOpts.ProcessBBCode(true),
		widget.TextOpts.Text(
			value,
			&text.GoTextFace{
				Source: loader.GetInstance().GetFont(loader.KyivRegularFont),
				Size:   20,
			},
			color.White)))
}

// CleanText cleans text component in the container.
func (sc *NotificationComponent) CleanText() {
	sc.container.GetWidget().Visibility = widget.Visibility_Hide

	sc.container.RemoveChildren()
}

// GetContainer retrieves container widget.
func (sc *NotificationComponent) GetContainer() *widget.Container {
	return sc.container
}

// newNotificationComponent initializes NotificationComponent.
func newNotificationComponent() *NotificationComponent {
	container := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(common.GetImageAsNineSlice(loader.PanelIdlePanel, 10, 10)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.TrackHover(false),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				Padding:            widget.Insets{Top: 10, Right: 10},
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	container.GetWidget().Visibility = widget.Visibility_Hide

	return &NotificationComponent{container: container}
}
