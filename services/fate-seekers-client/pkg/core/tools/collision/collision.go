package collision

import (
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
)

var (
	// GetInstance retrieves instance of the sounder, performing initial creation if needed.
	GetInstance = sync.OnceValue[*Collision](newCollision)
)

// Collision represents active map collision manager.
type Collision struct {
	// Represents collision objects mutex.
	collisionTileObjectMutex sync.Mutex

	// Represents collision tile objects.
	collisionTileObjects []dto.Position

	// Represents tertiary static objects mutex.
	mainTrackableObjectMutex sync.Mutex

	// Represents main trackable soundable object.
	mainTrackableObject dto.Position
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
func (s *Sounder) SetMainTrackableObject(value dto.Position) {
	s.mainTrackableObjectMutex.Lock()

	s.mainTrackableObject = value

	s.mainTrackableObjectMutex.Unlock()
}

// AddSoundables adds new soundable tile object with the provided value.
func (s *Sounder) AddSoundableTileObject(value *dto.SoundableTile) {
	s.soundableTileObjectMutex.Lock()

	s.soundableTileObjects = append(s.soundableTileObjects, value)

	s.soundableTileObjectMutex.Unlock()
}

// Update performs update operation for all the soundable objects.
func (s *Sounder) Update() {
	s.soundableTileObjectMutex.Lock()

	s.mainTrackableObjectMutex.Lock()

	s.mainTrackableObjectMutex.Unlock()

	s.externalTrackableObjectMutex.RLock()

	s.externalTrackableObjectMutex.RUnlock()

	s.soundableTileObjectMutex.Unlock()
}

// Clean performs clean operation for the configured collision holders.
func (c *Collision) Clean() {
	s.soundableTileObjectMutex.Lock()

	clear(s.soundableTileObjects)

	s.soundableTileObjectMutex.Unlock()

	s.externalTrackableObjectMutex.Lock()

	clear(s.externalTrackableObjects)

	s.externalTrackableObjectMutex.Unlock()
}

// newCollision initializes Collision.
func newCollision() *Collision {
	return &Collision{
		externalTrackableObjects: make(map[string]dto.Position),
	}
}
