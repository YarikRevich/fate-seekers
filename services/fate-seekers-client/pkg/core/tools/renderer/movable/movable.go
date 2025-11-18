package movable

import (
	"image/color"
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition/transparent"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/interpolation"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/kamera/v2"
)

const (
	// Repreesnts effect update ticker frequency.
	updateTickerFrequency = time.Millisecond * 150
)

// Represents movable object to be rendered.
type Movable struct {
	// Represents metadata selected for the exact movable object.
	metadata dto.ProcessedMovableMetadataSet

	// Represents currently selected direction.
	direction string

	// Represents if selected animation is static.
	static bool

	// Represents current frame to be updated.
	frame int

	// Represents mutex used for delayed positions.
	delayedMutex sync.Mutex

	// Represents delayed movable positions.
	delayedPositions []dto.Position

	// Represents current movable position.
	position dto.Position

	// Represents ticker used for movable unit frame updates.
	ticker *time.Ticker

	// Represents hit transparent transition effect.
	normalHitTransparentTransitionEffect transition.TransitionEffect

	// Represents accumulated draw image options used for camera processing.
	opts ebiten.DrawImageOptions
}

// TriggerNormalHit triggers normal git transparent transition effect.
func (m *Movable) TriggerNormalHit() {
	if m.normalHitTransparentTransitionEffect.Done() {
		m.normalHitTransparentTransitionEffect.Reset()
	}
}

// SetDirection sets direction value for the movable unit.
func (m *Movable) SetDirection(value string) {
	if m.direction != value {
		m.frame = 0
		m.direction = value
	}
}

// SetStatic sets static value for the movable unit.
func (m *Movable) SetStatic(value bool) {
	if m.static != value {
		m.frame = 0
		m.static = value
	}
}

// AddPosition adds position value for the movable unit.
func (m *Movable) AddPosition(value dto.Position) {
	m.delayedMutex.Lock()

	var delayedPositions []dto.Position

	if len(m.delayedPositions) != 0 {
		delayedPositions = interpolation.GetDelayedPositions(m.delayedPositions[len(m.delayedPositions)-1], value)
	} else {
		delayedPositions = interpolation.GetDelayedPositions(m.position, value)
	}

	if len(delayedPositions) != 0 {
		m.delayedPositions = append(m.delayedPositions, delayedPositions...)
	}

	m.delayedMutex.Unlock()
}

// GetPosition retrieves current position.
func (m *Movable) GetPosition() dto.Position {
	return m.position
}

// GetShiftBounds retrieves animation shift bounds.
func (m *Movable) GetShiftBounds() (float64, float64) {
	var shiftWidth, shiftHeight int

	if m.static {
		shiftWidth = m.metadata[m.direction].Rotation.Bounds().Dx()
		shiftHeight = m.metadata[m.direction].Rotation.Bounds().Dy()
	} else {
		shiftWidth = m.metadata[m.direction].Frames[m.frame].Bounds().Dx()
		shiftHeight = m.metadata[m.direction].Frames[m.frame].Bounds().Dy()
	}

	return float64(shiftWidth), float64(shiftHeight)
}

// Update performs update operation for the movable unit.
func (m *Movable) Update() {
	if !m.static {
		select {
		case <-m.ticker.C:
			m.ticker.Stop()

			m.frame = (m.frame + 1) % (len(m.metadata[m.direction].Frames) - 1)

			m.ticker.Reset(updateTickerFrequency)
		default:
		}
	}

	if !m.normalHitTransparentTransitionEffect.Done() {
		if !m.normalHitTransparentTransitionEffect.OnEnd() {
			m.normalHitTransparentTransitionEffect.Update()
		} else {
			m.normalHitTransparentTransitionEffect.Clean()
		}
	}
}

// Draw performs draw operation for the movable unit.
func (m *Movable) Draw(screen *ebiten.Image, centered bool, camera *kamera.Camera) {
	m.delayedMutex.Lock()

	if len(m.delayedPositions) != 0 {
		m.position = m.delayedPositions[0]

		m.delayedPositions = m.delayedPositions[1:]
	}

	m.delayedMutex.Unlock()

	m.opts.GeoM.Reset()

	m.opts.ColorM.Reset()

	if !centered {
		m.opts.GeoM.Translate(m.position.X, -m.position.Y)
	} else {
		shiftWidth, shiftHeight := m.GetShiftBounds()

		m.opts.GeoM.Translate(
			((float64(config.GetWorldWidth()))/2)-(shiftWidth/2),
			(float64(config.GetWorldHeight())/2)-(shiftHeight/2))
	}

	if !m.normalHitTransparentTransitionEffect.Done() {
		m.opts.ColorM.ScaleWithColor(
			color.RGBA{
				R: uint8(m.normalHitTransparentTransitionEffect.GetValue()),
				G: 0,
				B: 0,
				A: 255})
	}

	if m.static {
		if !centered {
			camera.Draw(m.metadata[m.direction].Rotation, &m.opts, screen)
		} else {
			screen.DrawImage(m.metadata[m.direction].Rotation, &m.opts)
		}
	} else {
		if !centered {
			camera.Draw(m.metadata[m.direction].Frames[m.frame], &m.opts, screen)
		} else {
			screen.DrawImage(m.metadata[m.direction].Frames[m.frame], &m.opts)
		}
	}
}

// NewMovable creates new Movable instance.
func NewMovable(path string) *Movable {
	normalHitTransparentTransitionEffect := transparent.NewTransparentTransitionEffect(false, 150, 255, 4, time.Millisecond*10)

	normalHitTransparentTransitionEffect.Clean()

	return &Movable{
		metadata:                             loader.GetInstance().GetMovable(path),
		ticker:                               time.NewTicker(updateTickerFrequency), // TODO: add some random coefficient
		normalHitTransparentTransitionEffect: normalHitTransparentTransitionEffect,
	}
}
