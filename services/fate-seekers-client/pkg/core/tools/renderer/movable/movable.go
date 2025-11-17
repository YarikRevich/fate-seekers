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
	"github.com/tidwall/btree"
)

const (
	// Repreesnts effect update ticker frequency.
	updateTickerFrequency = time.Millisecond * 150
)

type MovableUnit struct {
	// Represents metadata selected for the exact movable object.
	metadata dto.ProcessedMovableMetadataSet

	// Represents camera lock option.
	cameraLock bool

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

	opts ebiten.DrawImageOptions
}

// TriggerNormalHit triggers normal git transparent transition effect.
func (m *MovableUnit) TriggerNormalHit() {
	if m.normalHitTransparentTransitionEffect.Done() {
		m.normalHitTransparentTransitionEffect.Reset()
	}
}

// SetCameraLock sets camera lock value.
func (m *MovableUnit) SetCameraLock(value bool) {
	m.cameraLock = value
}

// SetDirection sets direction value for the movable unit.
func (m *MovableUnit) SetDirection(value string) {
	if m.direction != value {
		m.frame = 0
		m.direction = value
	}
}

// SetStatic sets static value for the movable unit.
func (m *MovableUnit) SetStatic(value bool) {
	if m.static != value {
		m.frame = 0
		m.static = value
	}
}

// AddPosition adds position value for the movable unit.
func (m *MovableUnit) AddPosition(value dto.Position) {
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
func (m *MovableUnit) GetPosition() dto.Position {
	return m.position
}

// GetShiftBounds retrieves animation shift bounds.
func (m *MovableUnit) GetShiftBounds() (float64, float64) {
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
func (m *MovableUnit) Update() {
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
func (m *MovableUnit) Draw(screen *ebiten.Image, camera *kamera.Camera) {
	m.delayedMutex.Lock()

	if len(m.delayedPositions) != 0 {
		m.position = m.delayedPositions[0]

		m.delayedPositions = m.delayedPositions[1:]
	}

	m.delayedMutex.Unlock()

	m.opts.GeoM.Reset()

	m.opts.ColorM.Reset()

	if !m.cameraLock {
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
		if !m.cameraLock {
			camera.Draw(m.metadata[m.direction].Rotation, &m.opts, screen)
		} else {
			screen.DrawImage(m.metadata[m.direction].Rotation, &m.opts)
		}
	} else {
		if !m.cameraLock {
			camera.Draw(m.metadata[m.direction].Frames[m.frame], &m.opts, screen)
		} else {
			screen.DrawImage(m.metadata[m.direction].Frames[m.frame], &m.opts)
		}
	}
}

// NewMovableUnit creates new MovableUnit instance.
func NewMovableUnit(path string) *MovableUnit {
	normalHitTransparentTransitionEffect := transparent.NewTransparentTransitionEffect(false, 150, 255, 4, time.Millisecond*10)

	normalHitTransparentTransitionEffect.Clean()

	return &MovableUnit{
		metadata:                             loader.GetInstance().GetMovable(path),
		ticker:                               time.NewTicker(updateTickerFrequency), // TODO: add some random coefficient
		normalHitTransparentTransitionEffect: normalHitTransparentTransitionEffect,
	}
}

// Movables represents movable objects holder.
type Movables struct {
	// Represents secondary objects mutex.
	secondaryObjectsMutex sync.RWMutex

	// Represents secondary objects to be rendered in the background.
	secondaryObjects map[string]*MovableUnit

	// Represents main objects mutex.
	mainObjectsMutex sync.RWMutex

	// Represents main movable objects to be rendered in front
	// (in this case main means character and visible inventory)
	mainObjects map[string]*MovableUnit

	// Represents objects position mutex.
	objectPositionMutex sync.RWMutex

	// Represents objects positions, which define rendering order.
	objectPosition *btree.Map[float64, []dto.AnimatorMovablePositionItem]
}

// Clean performs removal for all the configured objects.
func (m *Movables) Clean() {
	m.secondaryObjectsMutex.Lock()

	clear(m.secondaryObjects)

	m.secondaryObjectsMutex.Unlock()

	m.mainObjectsMutex.Lock()

	clear(m.mainObjects)

	m.mainObjectsMutex.Unlock()

	m.objectPositionMutex.Lock()

	m.objectPosition.Clear()

	m.objectPositionMutex.Unlock()
}

// PruneSecondary performs clean operation for abondoned secondary movables.
func (m *Movables) PruneSecondary(names map[string]bool) {
	m.secondaryObjectsMutex.Lock()

	for name := range m.secondaryObjects {
		if _, ok := names[name]; !ok {
			delete(m.secondaryObjects, name)
		}
	}

	m.secondaryObjectsMutex.Unlock()
}

// SecondaryExists checks if secondary movable object with the provided name exists.
func (m *Movables) SecondaryExists(name string) bool {
	m.secondaryObjectsMutex.RLock()

	_, ok := m.secondaryObjects[name]

	m.secondaryObjectsMutex.RUnlock()

	return ok
}

// AddSecondary adds new secondary movable object with the provided name and value.
func (m *Movables) AddSecondary(name string, value *MovableUnit) {
	m.secondaryObjectsMutex.Lock()

	m.secondaryObjects[name] = value

	m.secondaryObjectsMutex.Unlock()
}

// GetSecondary retrieves secondary movable object with the provided name.
func (m *Movables) GetSecondary(name string) *MovableUnit {
	m.secondaryObjectsMutex.RLock()

	result, _ := m.secondaryObjects[name]

	m.secondaryObjectsMutex.RUnlock()

	return result
}

// MainExists checks if main movable object with the provided name exists.
func (m *Movables) MainExists(name string) bool {
	m.mainObjectsMutex.RLock()

	_, ok := m.mainObjects[name]

	m.mainObjectsMutex.RUnlock()

	return ok
}

// AddMain adds new main movable object with the provided name and value.
func (m *Movables) AddMain(name string, value *MovableUnit) {
	m.mainObjectsMutex.Lock()

	m.mainObjects[name] = value

	m.mainObjectsMutex.Unlock()
}

// GetMain retrieves main movable object with the provided name.
func (m *Movables) GetMain(name string) *MovableUnit {
	m.mainObjectsMutex.RLock()

	result, _ := m.mainObjects[name]

	m.mainObjectsMutex.RUnlock()

	return result
}

// Update performs update operation and position rearangemenet for all the configured objects.
func (m *Movables) Update() {
	m.secondaryObjectsMutex.RLock()

	m.objectPosition.Clear()

	var (
		presentObjectPositions []dto.AnimatorMovablePositionItem
		ok                     bool
	)

	for issuer, movable := range m.secondaryObjects {
		movable.Update()

		m.objectPositionMutex.RLock()

		presentObjectPositions, ok = m.objectPosition.Get(movable.GetPosition().Y)
		if ok {
			presentObjectPositions = append(
				presentObjectPositions,
				dto.AnimatorMovablePositionItem{
					Issuer: issuer,
					Type:   dto.AnimatorMovablePositionItemSecondary})
		} else {
			presentObjectPositions = []dto.AnimatorMovablePositionItem{
				dto.AnimatorMovablePositionItem{
					Issuer: issuer,
					Type:   dto.AnimatorMovablePositionItemSecondary}}
		}

		m.objectPositionMutex.RUnlock()

		m.objectPositionMutex.Lock()

		m.objectPosition.Set(movable.GetPosition().Y, presentObjectPositions)

		m.objectPositionMutex.Unlock()
	}

	m.secondaryObjectsMutex.RUnlock()

	m.mainObjectsMutex.RLock()

	for issuer, movable := range m.mainObjects {
		movable.Update()

		m.objectPositionMutex.RLock()

		// var shiftHeight int

		// if movable.static {
		// 	shiftHeight = movable.metadata[movable.direction].Rotation.Bounds().Dy()
		// } else {
		// 	shiftHeight = movable.metadata[movable.direction].Frames[movable.frame].Bounds().Dy()
		// }

		presentObjectPositions, ok = m.objectPosition.Get(movable.position.Y) //float64(config.GetWorldHeight()-shiftHeight) / 2)
		if ok {
			presentObjectPositions = append(
				presentObjectPositions,
				dto.AnimatorMovablePositionItem{
					Issuer: issuer,
					Type:   dto.AnimatorMovablePositionItemMain})
		} else {
			presentObjectPositions = []dto.AnimatorMovablePositionItem{
				dto.AnimatorMovablePositionItem{
					Issuer: issuer,
					Type:   dto.AnimatorMovablePositionItemMain}}
		}

		m.objectPositionMutex.RUnlock()

		m.objectPositionMutex.Lock()

		m.objectPosition.Set(movable.position.Y, presentObjectPositions)

		m.objectPositionMutex.Unlock()
	}

	m.mainObjectsMutex.RUnlock()
}

// Draw performs draw operation for all the configured objects.
func (m *Movables) Draw(screen *ebiten.Image, camera *kamera.Camera) {
	m.objectPositionMutex.RLock()

	m.objectPosition.Reverse(func(key float64, value []dto.AnimatorMovablePositionItem) bool {
		for _, movable := range value {
			switch movable.Type {
			case dto.AnimatorMovablePositionItemSecondary:
				m.secondaryObjects[movable.Issuer].Draw(screen, camera)

			case dto.AnimatorMovablePositionItemMain:
				m.mainObjects[movable.Issuer].Draw(screen, camera)

			}
		}

		return true
	})

	m.objectPositionMutex.RUnlock()
}

// NewMovables creates new Movables holder.
func NewMovables() *Movables {
	return &Movables{
		secondaryObjects: make(map[string]*MovableUnit),
		mainObjects:      make(map[string]*MovableUnit),
		objectPosition:   btree.NewMap[float64, []dto.AnimatorMovablePositionItem](32),
	}
}
