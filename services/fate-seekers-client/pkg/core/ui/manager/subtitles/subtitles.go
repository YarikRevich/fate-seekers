package subtitles

import (
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/subtitles"
)

var (
	// GetInstance retrieves instance of the subtitles manager, performing initial creation if needed.
	GetInstance = sync.OnceValue[*SubtitlesManager](newSubtitlesManager)
)

// subtitlesUnit represents a subtitle unit.
type subtitlesUnit struct {
	// Represents text for of the subtitles unit.
	text string

	// Represents duration of which subtitles unit will be shown.
	duration time.Duration
}

// SubtitlesManager represents subtitles manager, which acts in a queue manner.
type SubtitlesManager struct {
	// Represents timer for the currently selected subtitle.
	timer *time.Timer

	// Represents check if text has already been updated.
	textUpdated bool

	// Represents queue of subtitles units.
	queue []*subtitlesUnit
}

// Update updates currently shown subtitles.
func (sm *SubtitlesManager) Update() {
	if len(sm.queue) > 0 {
		subtitleUnit := sm.queue[0]

		if sm.timer == nil {
			sm.timer = time.NewTimer(subtitleUnit.duration)
		}

		if !sm.textUpdated {
			subtitles.GetInstance().SetText(subtitleUnit.text)

			sm.textUpdated = true
		}

		select {
		case <-sm.timer.C:
			sm.queue = sm.queue[1:]

			subtitles.GetInstance().CleanText()

			sm.textUpdated = false

			sm.timer = nil
		default:
		}
	}
}

// Push pushes new value to the subtitles queue.
func (sm *SubtitlesManager) Push(text string, duration time.Duration) {
	sm.queue = append(sm.queue, &subtitlesUnit{
		text:     text,
		duration: duration,
	})
}

// Reset removes all the values from the queue.
func (sm *SubtitlesManager) Reset() {
	sm.queue = sm.queue[:0]
}

// newSubtitlesManager initializes SubtitlesManager.
func newSubtitlesManager() *SubtitlesManager {
	return new(SubtitlesManager)
}
