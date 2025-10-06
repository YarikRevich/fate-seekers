package loadingstars

import (
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/particle"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/particle/loadingstars/star"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	// Repreesnts effect update ticker frequency.
	updateTickerFrequency = time.Millisecond * 300

	// Represents max and default amount of elements being used by the particle effect.
	defaultCount = 256

	// Represents decrement count value used for particle effect updates.
	countDecrement = 7

	// Represents max and default divider value for elements being used by the particle effect.
	defaultDivider = 32
)

// StarsParticleEffect represents stars particle effect.
type StarsParticleEffect struct {
	// Represents particle effect time ticker used for particle effect progression.
	ticker *time.Ticker

	// Represents amount of particles, which are currently used for the effect.
	count uint16

	// Represents divider value for the particles, which are currently used for the effect.
	divider float32

	// Represents a set of particles, which are used for the particle effect.
	particles []star.StarParticleElement

	// Represents if particle effect has been finished.
	finished bool
}

func (spe *StarsParticleEffect) Done() bool {
	return spe.finished
}

func (spe *StarsParticleEffect) OnEnd() bool {
	return spe.count < countDecrement
}

func (spe *StarsParticleEffect) Clean() {
	spe.particles = spe.particles[:0]

	spe.ticker.Stop()

	spe.ticker = nil

	spe.finished = true
}

func (spe *StarsParticleEffect) Reset() {
	spe.particles = make([]star.StarParticleElement, defaultCount)

	spe.ticker = time.NewTicker(updateTickerFrequency)

	spe.count = defaultCount

	spe.divider = defaultDivider

	spe.finished = false
}

func (spe *StarsParticleEffect) Update() {
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

func (spe *StarsParticleEffect) Draw(screen *ebiten.Image) {
	for i := uint16(0); i < spe.count; i++ {
		spe.particles[i].Draw(screen)
	}
}

// NewStarsParticleEffect initializes StarsParticleEffect.
func NewStarsParticleEffect() particle.ParticleEffect {
	return &StarsParticleEffect{
		count:     defaultCount,
		divider:   defaultDivider,
		particles: make([]star.StarParticleElement, defaultCount),
		ticker:    time.NewTicker(updateTickerFrequency),
	}
}
