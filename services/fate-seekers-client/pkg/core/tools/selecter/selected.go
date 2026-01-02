package selected

import (
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/kamera/v2"
	"github.com/solarlune/resolv"
)

var (
	// GetInstance retrieves instance of the selected, performing initial creation if needed.
	GetInstance = sync.OnceValue[*Selected](newSelected)
)

// Describes all the available configurations related to selected processing.
const (
	MIN_DISTANCE_TO_LOCAL_STATIC = 65
)

// Selected represents active map selected manager.
type Selected struct {
	// Represents tertiary static objects mutex.
	mainTrackableObjectMutex sync.Mutex

	// Represents main trackable soundable object.
	mainTrackableObject *resolv.ConvexPolygon

	// Represents selectable tile objects mutex.
	selectableTileObjectsMutex sync.Mutex

	// Represents selectable tile objects.
	selectableTileObjects []*resolv.ConvexPolygon

	// Represents local static movable objects mutex.
	localStaticObjectsMutex sync.RWMutex

	// Represents local static selectable objects.
	localStaticObjects map[string]*resolv.ConvexPolygon

	// Represents external movable objects mutex.
	externalMovableObjectsMutex sync.RWMutex

	// Represents external movable selectable objects.
	externalMovableObjects map[string]*resolv.ConvexPolygon

	// Represents cursor trackable object.
	cursorTrackableObject *resolv.ConvexPolygon
}

// SetMainTrackableObject sets main trackable object with the provided value.
func (s *Selected) SetMainTrackableObject(value dto.Position, shiftWidth, shiftHeight float64) {
	s.mainTrackableObjectMutex.Lock()

	s.mainTrackableObject = resolv.NewRectangle(
		value.X, value.Y, shiftWidth/2, shiftHeight/2)

	s.mainTrackableObjectMutex.Unlock()
}

// MainTrackableObjectExists checks if main trackable object exists with the provided value.
func (s *Selected) MainTrackableObjectExists() bool {
	s.mainTrackableObjectMutex.Lock()

	value := s.mainTrackableObject

	s.mainTrackableObjectMutex.Unlock()

	return value != nil
}

// PruneExternalMovableObjects performs clean operation for abondoned external movables.
func (s *Selected) PruneExternalMovableObjects(names map[string]bool) {
	s.externalMovableObjectsMutex.Lock()

	for name := range s.externalMovableObjects {
		if _, ok := names[name]; !ok {
			delete(s.externalMovableObjects, name)
		}
	}

	s.externalMovableObjectsMutex.Unlock()
}

// ExternalMovableObjectExists checks if external movable object with the provided name exists.
func (s *Selected) ExternalMovableObjectExists(name string) bool {
	s.externalMovableObjectsMutex.RLock()

	_, ok := s.externalMovableObjects[name]

	s.externalMovableObjectsMutex.RUnlock()

	return ok
}

// AddExternalMovableObject adds external movable object as a selectable one with the provided value.
func (s *Selected) AddExternalMovableObject(name string, value dto.Position, shiftWidth, shiftHeight float64) {
	s.externalMovableObjectsMutex.Lock()

	s.externalMovableObjects[name] = resolv.NewRectangle(
		value.X, value.Y, shiftWidth/2, shiftHeight/2)

	s.externalMovableObjectsMutex.Unlock()
}

// GetExternalMovableObject retrieves external movable object with the provided name.
func (s *Selected) GetExternalMovableObject(name string) *resolv.ConvexPolygon {
	s.externalMovableObjectsMutex.RLock()

	result, _ := s.externalMovableObjects[name]

	s.externalMovableObjectsMutex.RUnlock()

	return result
}

// AddSelectableTileObject adds new selectable tile object with the provided value.
func (s *Selected) AddSelectableTileObject(value *dto.SelectableTile) {
	s.selectableTileObjectsMutex.Lock()

	selected := resolv.NewConvexPolygon(
		value.Position.X, value.Position.Y,
		[]float64{
			float64(value.TileWidth/2) / 2.0, 0,
			float64(value.TileWidth / 2), float64(value.TileHeight/2) / 2.0,
			float64(value.TileWidth/2) / 2.0, float64(value.TileHeight / 2),
			0, float64(value.TileHeight/2) / 2.0,
		},
	)

	s.selectableTileObjects = append(s.selectableTileObjects, selected)

	s.selectableTileObjectsMutex.Unlock()
}

// AddSelectableStaticObject adds new selectable static object with the provided value.
func (s *Selected) AddSelectableStaticObject(key string, value *dto.SelectableStatic) {
	s.selectableTileObjectsMutex.Lock()

	selected := resolv.NewConvexPolygon(
		value.Position.X, value.Position.Y,
		[]float64{
			float64(value.TileWidth/2) / 2.0, 0,
			float64(value.TileWidth / 2), float64(value.TileHeight/2) / 2.0,
			float64(value.TileWidth/2) / 2.0, float64(value.TileHeight / 2),
			0, float64(value.TileHeight/2) / 2.0,
		},
	)

	s.localStaticObjects[key] = selected

	s.selectableTileObjectsMutex.Unlock()
}

// AddSelectableStaticObject adds new selectable static object with the provided value.
func (s *Selected) RemoveSelectableStaticObject(key string) {
	s.selectableTileObjectsMutex.Lock()

	delete(s.localStaticObjects, key)

	s.selectableTileObjectsMutex.Unlock()
}

// Scan performs scan operation with the provided camera.
func (s *Selected) Scan(camera *kamera.Camera) (dto.SelectedObjectDetails, bool) {
	var cursorPositionX, cursorPositionY int

	if store.GetApplicationStateGamepadEnabled() == value.GAMEPAD_ENABLED_APPLICATION_TRUE_VALUE && ebiten.IsFocused() {
		rawGamepadPosition := store.GetApplicationStateGamepadPointerPosition()

		cursorPositionX = int(rawGamepadPosition.X)
		cursorPositionY = int(rawGamepadPosition.Y)
	} else {
		cursorPositionX, cursorPositionY = ebiten.CursorPosition()
	}

	worldCursorPositionX := (camera.X + float64(cursorPositionX))
	worldCursorPositionY := -(camera.Y + float64(cursorPositionY))

	s.mainTrackableObjectMutex.Lock()

	if s.mainTrackableObject != nil {
		s.localStaticObjectsMutex.Lock()

		for _, object := range s.localStaticObjects {
			width := object.Bounds().Width()
			height := object.Bounds().Height()

			s.cursorTrackableObject.SetPosition(worldCursorPositionX-(width/4.25), worldCursorPositionY+float64(height))

			if s.cursorTrackableObject.IsIntersecting(object) {

				if object.DistanceTo(s.mainTrackableObject) > MIN_DISTANCE_TO_LOCAL_STATIC {
					continue
				}

				s.mainTrackableObjectMutex.Unlock()

				s.localStaticObjectsMutex.Unlock()

				return dto.SelectedObjectDetails{
					Position: dto.Position{
						X: object.Position().X,
						Y: object.Position().Y,
					},
					Kind: dto.SELECTED_LOCAL_STATIC_OBJECT,
				}, true
			}
		}

		s.localStaticObjectsMutex.Unlock()
	}

	s.mainTrackableObjectMutex.Unlock()

	s.externalMovableObjectsMutex.Lock()

	for _, object := range s.externalMovableObjects {
		width := object.Bounds().Width()
		height := object.Bounds().Height()

		s.cursorTrackableObject.SetPosition(worldCursorPositionX-(width/4.25), worldCursorPositionY+float64(height))

		if s.cursorTrackableObject.IsIntersecting(object) {
			s.externalMovableObjectsMutex.Unlock()

			return dto.SelectedObjectDetails{
				Position: dto.Position{
					X: object.Position().X,
					Y: object.Position().Y,
				},
				Kind: dto.SELECTED_MOVABLE_OBJECT,
			}, true
		}
	}

	s.externalMovableObjectsMutex.Unlock()

	s.selectableTileObjectsMutex.Lock()

	for _, object := range s.selectableTileObjects {
		width := object.Bounds().Width()
		height := object.Bounds().Height()

		s.cursorTrackableObject.SetPosition(worldCursorPositionX-(width/4.25), worldCursorPositionY+float64(height*1.5))

		if s.cursorTrackableObject.IsIntersecting(object) {
			s.selectableTileObjectsMutex.Unlock()

			return dto.SelectedObjectDetails{
				Position: dto.Position{
					X: object.Position().X,
					Y: object.Position().Y,
				},
				Kind: dto.SELECTED_TILE_OBJECT,
			}, true
		}
	}

	s.selectableTileObjectsMutex.Unlock()

	return dto.SelectedObjectDetails{}, false
}

// Clean performs clean operation for the configured collision holders.
func (s *Selected) Clean() {
	s.mainTrackableObjectMutex.Lock()

	s.mainTrackableObject = nil

	s.mainTrackableObjectMutex.Unlock()

	s.externalMovableObjectsMutex.Lock()

	clear(s.externalMovableObjects)

	s.externalMovableObjectsMutex.Unlock()

	s.selectableTileObjectsMutex.Lock()

	s.selectableTileObjects = s.selectableTileObjects[:0]

	s.selectableTileObjectsMutex.Unlock()

	s.localStaticObjectsMutex.Lock()

	clear(s.localStaticObjects)

	s.localStaticObjectsMutex.Unlock()
}

// newSelected initializes Selected.
func newSelected() *Selected {
	return &Selected{
		externalMovableObjects: make(map[string]*resolv.ConvexPolygon),
		localStaticObjects:     make(map[string]*resolv.ConvexPolygon),
		cursorTrackableObject:  resolv.NewRectangle(0, 0, 10, 10),
	}
}
