package movable

import (
	"fmt"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	// Repreesnts effect update ticker frequency.
	updateTickerFrequency = time.Millisecond * 100
)

type MovableUnit struct {
	// Represents metadata selected for the exact movable object.
	metadata dto.ProcessedMovableMetadataSet

	// Represents currently selected direction.
	direction string

	// Represents if selected animation is static.
	static bool

	// Represents current frame to be updated.
	frame int

	// Represents current movable position.
	position dto.Position

	// Represents ticker used for movable unit frame updates.
	ticker *time.Ticker
}

// SetDirection sets direction value for the movable unit.
func (m *MovableUnit) SetDirection(value string) {
	if m.direction != value {
		fmt.Println(value)

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

// SetPosition sets position value for the movable unit.
func (m *MovableUnit) SetPosition(value dto.Position) {
	m.position = value
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
}

// Draw performs draw operation for the movable unit.
func (m *MovableUnit) Draw(screen *ebiten.Image) {
	var opts ebiten.DrawImageOptions

	opts.GeoM.Translate(m.position.X, -m.position.Y)

	if m.static {
		screen.DrawImage(m.metadata[m.direction].Rotation, &opts)
	} else {
		screen.DrawImage(m.metadata[m.direction].Frames[m.frame], &opts)
	}
}

// NewMovableUnit creates new MovableUnit instance.
func NewMovableUnit(path string) *MovableUnit {
	return &MovableUnit{
		metadata: loader.GetInstance().GetMovable(path),
		ticker:   time.NewTicker(updateTickerFrequency), // TODO: add some random coefficient
	}
}

// Movables represents movable objects holder.
type Movables struct {
	objects map[string]*MovableUnit
}

// Clean performs removal for all the configured objects.
func (m *Movables) Clean() {
	clear(m.objects)
}

// Prune performs clean operation for abondoned animations.
func (m *Movables) Prune(names map[string]bool) {
	for name := range m.objects {
		if _, ok := names[name]; !ok {
			delete(m.objects, name)
		}
	}
}

// Exists checks if a movable object with the provided name exists.
func (m *Movables) Exists(name string) bool {
	_, ok := m.objects[name]

	return ok
}

// Add adds new movable object with the provided name and value.
func (m *Movables) Add(name string, value *MovableUnit) {
	m.objects[name] = value
}

// Get retrieves movable object with the provided name.
func (m *Movables) Get(name string) *MovableUnit {
	return m.objects[name]
}

// Update performs update operation for all the configured objects.
func (m *Movables) Update() {
	for _, movable := range m.objects {
		movable.Update()
	}
}

// Draw performs draw operation for all the configured objects.
func (m *Movables) Draw(screen *ebiten.Image) {
	for _, movable := range m.objects {
		movable.Draw(screen)
	}
}

// NewMovables creates new Movables holder.
func NewMovables() *Movables {
	return &Movables{
		objects: make(map[string]*MovableUnit),
	}
}
