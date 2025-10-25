package animator

import (
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/particle/loadingstars/star"
	"github.com/hajimehoshi/ebiten"
)

// TODO: create animator with direction objects

// Animator represents movable animation animator.
type Animator struct {
	// Represents particle effect time ticker used for particle effect progression.
	ticker *time.Ticker

	// Represents amount of particles, which are currently used for the effect.
	count uint16

	// Represents divider value for the particles, which are currently used for the effect.
	divider float32

	// Represents a set of particles, which are used for the particle effect.
	animations []star.StarParticleElement

	// Represents if particle effect has been finished.
	finished bool
}

func (a *Animator) Clean() {
	spe.particles = spe.particles[:0]

	spe.ticker.Stop()

	spe.ticker = nil

	spe.finished = true
}

func (a *Animator) Reset() {
	spe.particles = make([]star.StarParticleElement, defaultCount)

	spe.ticker = time.NewTicker(updateTickerFrequency)

	spe.count = defaultCount

	spe.divider = defaultDivider

	spe.finished = false
}

func (a *Animator) Update() {
	select {
	case <-spe.ticker.C:
		spe.ticker.Stop()

		if spe.count >= countDecrement {
			spe.count -= countDecrement
		}

		if spe.divider > 1 {
			spe.divider--
		}

		spe.ticker.Reset(updateTickerFrequency)
	default:
	}

	for i := uint16(0); i < spe.count; i++ {
		if spe.particles[i].GetDivider() != spe.divider {
			spe.particles[i].SetDivider(spe.divider)
		}

		spe.particles[i].Update()
	}
}

func (a *Animator) Draw(screen *ebiten.Image) {
	for i := uint16(0); i < spe.count; i++ {
		spe.particles[i].Draw(screen)
	}
}

// NewStarsParticleEffect initializes StarsParticleEffect.
func NewAnimator() Animator {
	return &StarsParticleEffect{
		count:     defaultCount,
		divider:   defaultDivider,
		particles: make([]star.StarParticleElement, defaultCount),
		ticker:    time.NewTicker(updateTickerFrequency),
	}
}
