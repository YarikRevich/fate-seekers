package sounder

import (
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/sound"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/solarlune/resolv"
)

var (
	// GetInstance retrieves instance of the sounder, performing initial creation if needed.
	GetInstance = sync.OnceValue[*Sounder](newSounder)
)

// Sounder represents active map sound manager.
type Sounder struct {
	// Represents collision polygons mutex.
	collisionPolygonsMutex sync.Mutex

	// Represents collision polygons.
	collisionPolygons []*resolv.ConvexPolygon

	// Represents soundable objects mutex.
	soundableTileObjectMutex sync.RWMutex

	// Represents soundable tile objects.
	soundableTileObjects map[uint32]*dto.SoundableTile

	// Represents external trackable objects mutex.
	externalTrackableObjectMutex sync.RWMutex

	// Represents external trackable soundable objects.
	externalTrackableObjects map[string]dto.Position

	// Represents tertiary static objects mutex.
	mainTrackableObjectMutex sync.Mutex

	// Represents main trackable soundable object.
	mainTrackableObject *resolv.ConvexPolygon

	//
	mainTrackableObjectUpdated bool
}

// PruneExternalTrackableObjects performs clean operation for abondoned external trackable objects.
func (s *Sounder) PruneExternalTrackableObjects(names map[string]bool) {
	s.externalTrackableObjectMutex.Lock()

	for name := range s.externalTrackableObjects {
		if _, ok := names[name]; !ok {
			delete(s.externalTrackableObjects, name)
		}
	}

	s.externalTrackableObjectMutex.Unlock()
}

// SetExternalTrackableObject sets external trackable object with the provided key and value.
func (s *Sounder) SetExternalTrackableObject(key string, value dto.Position) {
	s.externalTrackableObjectMutex.Lock()

	s.externalTrackableObjects[key] = value

	s.externalTrackableObjectMutex.Unlock()
}

// SetMainTrackableObject sets main trackable object with the provided value.
func (s *Sounder) SetMainTrackableObject(value dto.Position, shiftWidth, shiftHeight float64) {
	s.mainTrackableObjectMutex.Lock()

	if s.mainTrackableObject == nil ||
		(s.mainTrackableObject.Position().X != value.X ||
			s.mainTrackableObject.Position().Y != value.Y) {
		s.mainTrackableObjectUpdated = true
	}

	s.mainTrackableObject = resolv.NewRectangle(
		value.X, value.Y, shiftWidth/2, shiftHeight/2)

	s.mainTrackableObjectMutex.Unlock()
}

// AddSoundables adds new soundable tile object with the provided value.
func (s *Sounder) AddSoundableTileObject(value *dto.SoundableTile) {
	s.collisionPolygonsMutex.Lock()

	collider := resolv.NewConvexPolygon(
		value.Position.X-(float64(value.TileWidth)/2), value.Position.Y-(float64(value.TileHeight)/2),
		[]float64{
			float64(value.TileWidth) / 2.0, 0,
			float64(value.TileWidth), float64(value.TileHeight) / 2.0,
			float64(value.TileWidth) / 2.0, float64(value.TileHeight),
			0, float64(value.TileHeight) / 2.0,
		},
	)

	s.soundableTileObjectMutex.Lock()

	s.soundableTileObjects[collider.ID()] = value

	s.soundableTileObjectMutex.Unlock()

	s.collisionPolygons = append(s.collisionPolygons, collider)

	s.collisionPolygonsMutex.Unlock()
}

// InterruptMainTrackableObject performs sound interruption for main trackable object.
func (s *Sounder) InterruptMainTrackableObject() {
	sound.GetInstance().GetSoundSounderMainFxManager().Stop()
}

// Update performs update operation for all the soundable objects.
func (s *Sounder) Update() {
	s.collisionPolygonsMutex.Lock()

	s.soundableTileObjectMutex.RLock()

	s.mainTrackableObjectMutex.Lock()

	if s.mainTrackableObjectUpdated {
		s.mainTrackableObjectUpdated = false

		for _, polygon := range s.collisionPolygons {
			if polygon.IsIntersecting(s.mainTrackableObject) {
				if !sound.GetInstance().GetSoundSounderMainFxManager().IsFXPlaying() {
					switch s.soundableTileObjects[polygon.ID()].Name {
					case loader.TilemapSoundRockValue:
						sound.GetInstance().GetSoundSounderMainFxManager().PushWithHandbrake(loader.RockFXSound)
					}
				}

				break
			}
		}
	}

	// TODO: do main objects tracking.

	s.mainTrackableObjectMutex.Unlock()

	s.externalTrackableObjectMutex.RLock()

	// TODO: do external objects tracking.

	s.externalTrackableObjectMutex.RUnlock()

	s.soundableTileObjectMutex.RUnlock()

	s.collisionPolygonsMutex.Unlock()
}

// Clean performs clean operation for the configured soundable holders.
func (s *Sounder) Clean() {
	s.soundableTileObjectMutex.Lock()

	clear(s.soundableTileObjects)

	s.soundableTileObjectMutex.Unlock()

	s.mainTrackableObjectMutex.Lock()

	s.mainTrackableObject = nil

	s.mainTrackableObjectUpdated = false

	s.mainTrackableObjectMutex.Unlock()

	s.externalTrackableObjectMutex.Lock()

	clear(s.externalTrackableObjects)

	s.externalTrackableObjectMutex.Unlock()
}

// newSounder initializes Sounder.
func newSounder() *Sounder {
	return &Sounder{
		soundableTileObjects:     make(map[uint32]*dto.SoundableTile),
		externalTrackableObjects: make(map[string]dto.Position),
	}
}
