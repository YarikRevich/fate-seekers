package sounder

import (
	"math"
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/sound"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/setanarut/kamera/v2"
	"github.com/solarlune/resolv"
)

var (
	// GetInstance retrieves instance of the sounder, performing initial creation if needed.
	GetInstance = sync.OnceValue[*Sounder](newSounder)
)

// Represents static sound options.
const (
	CAMERA_SOUND_OFFSET = 60
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
	externalTrackableObjects map[string]*dto.ExternalSounderObject

	// Represents tertiary static objects mutex.
	mainTrackableObjectMutex sync.Mutex

	// Represents main trackable soundable object.
	mainTrackableObject *resolv.ConvexPolygon

	// Represents if main trackable object is updated and should be processed again.
	mainTrackableObjectUpdated bool
}

// PruneExternalTrackableObjects performs clean operation for abondoned external trackable objects.
func (s *Sounder) PruneExternalTrackableObjects(names map[string]bool) {
	s.externalTrackableObjectMutex.Lock()

	for name := range s.externalTrackableObjects {
		if _, ok := names[name]; !ok {
			if sound.GetInstance().SoundSounderExternalFxManagerExists(name) {
				sound.GetInstance().RemoveSoundSounderExternalFxManager(name)
			}

			delete(s.externalTrackableObjects, name)
		}
	}

	s.externalTrackableObjectMutex.Unlock()
}

// ExternalTrackableObjectExists checks if trackable object with provided key exists.
func (s *Sounder) ExternalTrackableObjectExists(key string) bool {
	_, ok := s.externalTrackableObjects[key]

	return ok
}

// AddExternalTrackableObject adds external trackable object with the provided key and value.
func (s *Sounder) AddExternalTrackableObject(key string, value dto.Position, shiftWidth, shiftHeight float64) {
	s.externalTrackableObjectMutex.Lock()

	s.externalTrackableObjects[key] = &dto.ExternalSounderObject{
		Updated: true,
		Polygon: resolv.NewRectangle(
			value.X, value.Y, shiftWidth/2, shiftHeight/2),
	}

	s.externalTrackableObjectMutex.Unlock()
}

// GetExternalTrackableObject retrieves external trackable object with the provided key.
func (s *Sounder) GetExternalTrackableObject(key string) *dto.ExternalSounderObject {
	return s.externalTrackableObjects[key]
}

// SetMainTrackableObject sets main trackable object with the provided value.
func (s *Sounder) SetMainTrackableObject(value dto.Position, shiftWidth, shiftHeight float64) {
	s.mainTrackableObjectMutex.Lock()

	s.mainTrackableObjectUpdated = true

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
	if sound.GetInstance().GetSoundSounderMainFxManager().IsFXPlaying() {
		sound.GetInstance().GetSoundSounderMainFxManager().StopFXPlaying()
	}
}

// InterruptExternalTrackableObject performs sound interruption for external trackable object.
func (s *Sounder) InterruptExternalTrackableObject(issuer string) {
	if sound.GetInstance().SoundSounderExternalFxManagerExists(issuer) {
		if sound.GetInstance().GetSoundSounderExternalFxManager(issuer).IsFXPlaying() {
			sound.GetInstance().GetSoundSounderExternalFxManager(issuer).StopFXPlaying()
		}
	}
}

// Update performs update operation for all the soundable objects.
func (s *Sounder) Update(camera *kamera.Camera) {
	s.collisionPolygonsMutex.Lock()

	s.soundableTileObjectMutex.RLock()

	s.mainTrackableObjectMutex.Lock()

	if s.mainTrackableObjectUpdated {
		s.mainTrackableObjectUpdated = false

		for _, polygon := range s.collisionPolygons {
			if polygon.IsIntersecting(s.mainTrackableObject) {
				if !sound.GetInstance().GetSoundSounderMainFxManager().IsFXPlaying() {
					value, ok := s.soundableTileObjects[polygon.ID()]
					if !ok {
						continue
					}

					switch value.Name {
					case loader.TilemapSoundRockValue:
						sound.GetInstance().GetSoundSounderMainFxManager().PushWithHandbrake(loader.RockFXSound)
					}
				}

				break
			}
		}
	}

	s.mainTrackableObjectMutex.Unlock()

	s.externalTrackableObjectMutex.RLock()

	var wg sync.WaitGroup

	wg.Add(len(s.externalTrackableObjects))

	minCameraViewportWidth := store.GetPositionSession().X - (math.Abs(camera.CenterOffsetX) + CAMERA_SOUND_OFFSET)
	maxCameraViewportWidth := store.GetPositionSession().X + (math.Abs(camera.CenterOffsetX) + CAMERA_SOUND_OFFSET)

	minCameraViewportHeight := store.GetPositionSession().Y - (math.Abs(camera.CenterOffsetY) + CAMERA_SOUND_OFFSET)
	maxCameraViewportHeight := store.GetPositionSession().Y + (math.Abs(camera.CenterOffsetY) + CAMERA_SOUND_OFFSET)

	for key, externalTrackableObject := range s.externalTrackableObjects {
		if externalTrackableObject.Updated {
			externalTrackableObject.Updated = false

			if (externalTrackableObject.Polygon.Position().X >= minCameraViewportWidth && externalTrackableObject.Polygon.Position().X <= maxCameraViewportWidth) &&
				(externalTrackableObject.Polygon.Position().Y >= minCameraViewportHeight && externalTrackableObject.Polygon.Position().Y <= maxCameraViewportHeight) {
				go func(wg *sync.WaitGroup) {
					for _, polygon := range s.collisionPolygons {
						if polygon.IsIntersecting(externalTrackableObject.Polygon) {
							if !sound.GetInstance().SoundSounderExternalFxManagerExists(key) {
								sound.GetInstance().AddSoundSounderExternalFxManager(key)
							}

							if !sound.GetInstance().GetSoundSounderExternalFxManager(key).IsFXPlaying() {
								switch s.soundableTileObjects[polygon.ID()].Name {
								case loader.TilemapSoundRockValue:
									sound.GetInstance().GetSoundSounderExternalFxManager(key).PushWithHandbrake(loader.RockFXSound)
								}
							}

							break
						}
					}

					wg.Done()
				}(&wg)
			} else {
				wg.Done()
			}
		} else {
			wg.Done()
		}
	}

	wg.Wait()

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
		externalTrackableObjects: make(map[string]*dto.ExternalSounderObject),
	}
}
