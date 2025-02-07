package notification

import (
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/notification"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
)

var (
	// GetInstance retrieves instance of the notification manager, performing initial creation if needed.
	GetInstance = sync.OnceValue[*NotificationManager](newNotificationManager)
)

// NotificationManager represents notification manager, which acts in a queue manner.
type NotificationManager struct {
	// Represents timer for the currently selected subtitle.
	timer *time.Timer

	// Represents notification visibility.
	visible bool

	// Represents check if text has already been updated.
	textUpdated bool

	// Represents queue of notification units.
	queue []*dto.NotificationUnit
}

func (sm *NotificationManager) GetTextUpdated() bool {
	return sm.textUpdated
}

// GetVisible retrieves if notification manager is visible.
func (sm *NotificationManager) GetVisible() bool {
	return sm.visible
}

// ToggleVisible sets notification visibility to be toggled.
func (sm *NotificationManager) ToggleVisible() {
	sm.visible = !sm.visible
}

// Update updates currently shown subtitles.
func (sm *NotificationManager) Update() {
	if len(sm.queue) > 0 {
		notificationUnit := sm.queue[0]

		if !sm.textUpdated {
			notification.GetInstance().SetText(notificationUnit.Text)

			sm.textUpdated = true
		}

		if sm.timer == nil && sm.visible {
			sm.timer = time.NewTimer(notificationUnit.Duration)
		}

		if sm.timer != nil {
			select {
			case <-sm.timer.C:
				sm.queue = sm.queue[1:]

				notification.GetInstance().CleanText()

				sm.textUpdated = false

				sm.visible = false

				sm.timer = nil
			default:
			}
		}
	}
}

// Push pushes new value to the notification queue.
func (sm *NotificationManager) Push(text string, duration time.Duration) {
	sm.queue = append(sm.queue, &dto.NotificationUnit{
		Text:     text,
		Duration: duration,
	})
}

// Reset removes all the values from the queue.
func (sm *NotificationManager) Reset() {
	sm.queue = sm.queue[:0]
}

func (sm *NotificationManager) IsEmpty() bool {
	return len(sm.queue) == 0
}

// newNotificationManager initializes NotificationManager.
func newNotificationManager() *NotificationManager {
	return new(NotificationManager)
}
