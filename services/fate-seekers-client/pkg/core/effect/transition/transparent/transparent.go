package transparent

import (
	"image/color"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition"
	"github.com/hajimehoshi/ebiten/v2"
)

// TransparentTransitionEffect represents transparent transition effect.
type TransparentTransitionEffect struct {
	// Represents transition time ticker used for transition progression.
	ticker *time.Ticker

	// Represents current state of the transition.
	counter uint8

	// Represents if transition effect has been finished.
	finished bool
}

func (tte *TransparentTransitionEffect) Done() bool {
	return tte.finished
}

func (tte *TransparentTransitionEffect) OnEnd() bool {
	return tte.counter == 255
}

func (tte *TransparentTransitionEffect) Update() {
	select {
	case <-tte.ticker.C:
		tte.ticker.Stop()

		tte.counter += 5

		tte.ticker.Reset(time.Microsecond * 10)
	}
}

func (tte *TransparentTransitionEffect) Clean() {
	tte.ticker.Stop()

	tte.ticker = nil

	tte.finished = true
}

func (tte *TransparentTransitionEffect) Reset() {
	tte.ticker = time.NewTicker(time.Microsecond * 10)

	tte.counter = 0

	tte.finished = false
}

func (tte *TransparentTransitionEffect) GetOptions() *ebiten.DrawImageOptions {
	var c ebiten.ColorM

	c.ScaleWithColor(color.RGBA{R: 255, G: 255, B: 255, A: tte.counter})

	return &ebiten.DrawImageOptions{ColorM: c}
}

// NewTransparentTransitionEffect initializes TransparentTransitionEffect.
func NewTransparentTransitionEffect() transition.TransitionEffect {
	return &TransparentTransitionEffect{
		ticker: time.NewTicker(time.Microsecond * 10),
	}
}
