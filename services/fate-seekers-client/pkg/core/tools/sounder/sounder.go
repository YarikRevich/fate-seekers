package sounder

import (
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
)

var (
	// GetInstance retrieves instance of the sounder, performing initial creation if needed.
	GetInstance = sync.OnceValue[*Sounder](newSounder)
)

// Sounder represents active map sound manager.
type Sounder struct {
	// Represents tertiary static objects mutex.
	soundableTileObjectMutex sync.Mutex

	// Represents soundable tile objects.
	soundableTileObjects []*dto.SoundableTile

	// Represents external trackable objects mutex.
	externalTrackableObjectMutex sync.RWMutex

	// Represents external trackable soundable objects.
	externalTrackableObjects map[string]dto.Position

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

// newRenderer initializes Sounder.
func newSounder() *Sounder {
	return &Sounder{
		externalTrackableObjects: make(map[string]dto.Position),
	}
}
