package transparent

import (
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition"
)

// TransparentTransitionEffect represents transparent transition effect.
type TransparentTransitionEffect struct {
	// Represents transition direction progression.
	forward bool

	// Represents max state of the transition.
	maxCounter float64

	// Represents min state of the transition.
	minCounter float64

	// Represents provided shift for the transition.
	shift float64

	// Represents transition timer period.
	period time.Duration

	// Represents transition time ticker used for transition progression.
	ticker *time.Ticker

	// Represents current state of the transition.
	counter float64

	// Represents if transition effect has been finished.
	finished bool
}

func (tte *TransparentTransitionEffect) Done() bool {
	return tte.finished
}

func (tte *TransparentTransitionEffect) OnEnd() bool {
	return tte.counter == tte.maxCounter
}

func (tte *TransparentTransitionEffect) Update() {
	select {
	case <-tte.ticker.C:
		tte.ticker.Stop()

		if tte.forward {
			tte.counter += tte.shift
		} else {
			tte.counter -= tte.shift
		}

		tte.ticker.Reset(tte.period)
	default:
	}
}

func (tte *TransparentTransitionEffect) Clean() {
	tte.ticker.Stop()

	tte.ticker = nil

	tte.finished = true
}

func (tte *TransparentTransitionEffect) Reset() {
	tte.ticker = time.NewTicker(tte.period)

	tte.counter = tte.minCounter

	tte.finished = false
}

func (tte *TransparentTransitionEffect) GetValue() float64 {
	return float64(tte.counter)
}

// NewTransparentTransitionEffect initializes TransparentTransitionEffect.
func NewTransparentTransitionEffect(
	forward bool, maxCounter, minCounter, shift float64, period time.Duration) transition.TransitionEffect {
	result := &TransparentTransitionEffect{
		forward:    forward,
		maxCounter: maxCounter,
		minCounter: minCounter,
		counter:    minCounter,
		shift:      shift,
		period:     period,
		ticker:     time.NewTicker(period),
	}

	return result
}
