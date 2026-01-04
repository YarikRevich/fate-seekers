package renderer

import (
	"math"
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/renderer/movable"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/renderer/static"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/renderer/tile"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/elliotchance/orderedmap/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/kamera/v2"
	"github.com/tidwall/btree"
)

var (
	// GetInstance retrieves instance of the renderer, performing initial creation if needed.
	GetInstance = sync.OnceValue[*Renderer](newRenderer)
)

// Represents static rendering options.
const (
	CAMERA_RENDERING_OFFSET = 60
)

// Renderer represents object renderer. It has three levels of rendered objects.
// The first level is main objects to be rendered in the center of a screen.
// The second level is for external objects to be rendered above the basic map layer.
// The third level is for the basic map layer.
type Renderer struct {
	// Represents tertiary static objects mutex.
	tertiaryTileObjectMutex sync.Mutex

	// Represents tertiary tilemap object to be rendered in the background.
	tertiaryTileObjects *orderedmap.OrderedMap[string, *tile.Tile]

	// Represents secondary objects mutex.
	secondaryTileObjectMutex sync.RWMutex

	// Represents secondary objects to be rendered in the background.
	secondaryTileObjects *orderedmap.OrderedMap[string, *tile.Tile]

	// Represents secondary movable objects mutex.
	secondaryLocalMovableObjectsMutex sync.RWMutex

	// Represents secondary movable objects to be rendered in the background.
	secondaryLocalMovableObjects map[string]*movable.Movable

	// Represents secondary static objects mutex.
	secondaryLocalStaticObjectsMutex sync.RWMutex

	// Represents secondary sttaic objects to be rendered in the background.
	secondaryLocalStaticObjects map[string]*static.Static

	// Represents secondary objects mutex.
	secondaryExternalMovableObjectsMutex sync.RWMutex

	// Represents secondary objects to be rendered in the background.
	secondaryExternalMovableObjects map[string]*movable.Movable

	// Represents main objects mutex.
	mainCenteredMovableObjectsMutex sync.RWMutex

	// Represents main movable objects to be rendered in front
	// (in this case main means character and visible inventory)
	mainCenteredMovableObjects map[string]*movable.Movable

	// Represents objects position mutex.
	objectPositionMutex sync.RWMutex

	// Represents objects positions, which define rendering order.
	objectPosition *btree.Map[float64, []dto.RendererPositionItem]

	// Represents selected object position mutex.
	selectedObjectPositionMutex sync.Mutex

	// Represents if selected object is enabled.
	selectedObjectEnabled bool

	// Represents selected object position, which is used to draw some additional glowing layer.
	selectedObjectPosition dto.Position

	// Represents ignored external movable objects mutex.
	ignoredExternalMovableObjectsMutex sync.Mutex

	// Represents ignored external movable objects, which are then used to reset delayed positions to show
	// the objects immediately.
	ignoredExternalMovableObjects map[string]bool
}

// TertiaryTileObjectExists checks if tertiary tile object exists.
func (r *Renderer) TertiaryTileObjectExists(name string) bool {
	r.tertiaryTileObjectMutex.Lock()

	ok := r.tertiaryTileObjects.Has(name)

	r.tertiaryTileObjectMutex.Unlock()

	return ok
}

// AddTertiaryTileObject adds new tertiary external tile object with the provided value.
func (r *Renderer) AddTertiaryTileObject(name string, value *tile.Tile) {
	r.tertiaryTileObjectMutex.Lock()

	r.tertiaryTileObjects.Set(name, value)

	r.tertiaryTileObjectMutex.Unlock()
}

// SecondaryTileObjectExists checks if secondary tile object exists.
func (r *Renderer) SecondaryTileObjectExists(name string) bool {
	r.secondaryTileObjectMutex.Lock()

	ok := r.secondaryTileObjects.Has(name)

	r.secondaryTileObjectMutex.Unlock()

	return ok
}

// AddSecondaryTileObject adds new secondary external tile object with the provided value.
func (r *Renderer) AddSecondaryTileObject(name string, value *tile.Tile) {
	r.secondaryTileObjectMutex.Lock()

	r.secondaryTileObjects.Set(name, value)

	r.secondaryTileObjectMutex.Unlock()
}

// PruneSecondaryExternalMovableObjects performs clean operation for abondoned secondary external movables.
func (r *Renderer) PruneSecondaryExternalMovableObjects(names map[string]bool) {
	r.secondaryExternalMovableObjectsMutex.Lock()

	for name := range r.secondaryExternalMovableObjects {
		if _, ok := names[name]; !ok {
			delete(r.secondaryExternalMovableObjects, name)
		}
	}

	r.secondaryExternalMovableObjectsMutex.Unlock()
}

// SecondaryExternalMovableObjectExists checks if secondary external movable object with the provided name exists.
func (r *Renderer) SecondaryExternalMovableObjectExists(name string) bool {
	r.secondaryExternalMovableObjectsMutex.RLock()

	_, ok := r.secondaryExternalMovableObjects[name]

	r.secondaryExternalMovableObjectsMutex.RUnlock()

	return ok
}

// AddSecondaryExternalMovableObject adds new secondary external movable object with the provided name and value.
func (r *Renderer) AddSecondaryExternalMovableObject(name string, value *movable.Movable) {
	r.secondaryExternalMovableObjectsMutex.Lock()

	r.secondaryExternalMovableObjects[name] = value

	r.secondaryExternalMovableObjectsMutex.Unlock()
}

// GetSecondaryExternalMovableObject retrieves secondary movable object with the provided name.
func (r *Renderer) GetSecondaryExternalMovableObject(name string) *movable.Movable {
	r.secondaryExternalMovableObjectsMutex.RLock()

	result, _ := r.secondaryExternalMovableObjects[name]

	r.secondaryExternalMovableObjectsMutex.RUnlock()

	return result
}

// SecondaryLocalStaticObjectExists checks if secondary local static object with the provided name exists.
func (r *Renderer) SecondaryLocalStaticObjectExists(name string) bool {
	r.secondaryLocalStaticObjectsMutex.RLock()

	_, ok := r.secondaryLocalStaticObjects[name]

	r.secondaryLocalStaticObjectsMutex.RUnlock()

	return ok
}

// AddSecondaryLocalStaticObject adds new secondary local static object with the provided name and value.
func (r *Renderer) AddSecondaryLocalStaticObject(name string, value *static.Static) {
	r.secondaryLocalStaticObjectsMutex.Lock()

	r.secondaryLocalStaticObjects[name] = value

	r.secondaryLocalStaticObjectsMutex.Unlock()
}

// GetSecondaryLocalStaticObject retrieves secondary local static object with the provided name.
func (r *Renderer) GetSecondaryLocalStaticObject(name string) *static.Static {
	r.secondaryLocalStaticObjectsMutex.RLock()

	result, _ := r.secondaryLocalStaticObjects[name]

	r.secondaryLocalStaticObjectsMutex.RUnlock()

	return result
}

// RemoveSecondaryLocalStaticObject removes secondary local static object with the provided name and value.
func (r *Renderer) RemoveSecondaryLocalStaticObject(name string) {
	r.secondaryLocalStaticObjectsMutex.Lock()

	delete(r.secondaryLocalStaticObjects, name)

	r.secondaryLocalStaticObjectsMutex.Unlock()
}

// MainCenteredMovableObjectExists checks if main centered movable object with the provided name exists.
func (r *Renderer) MainCenteredMovableObjectExists(name string) bool {
	r.mainCenteredMovableObjectsMutex.RLock()

	_, ok := r.mainCenteredMovableObjects[name]

	r.mainCenteredMovableObjectsMutex.RUnlock()

	return ok
}

// AddMainCenteredMovableObject adds new main centered movable object with the provided name and value.
func (r *Renderer) AddMainCenteredMovableObject(name string, value *movable.Movable) {
	r.mainCenteredMovableObjectsMutex.Lock()

	r.mainCenteredMovableObjects[name] = value

	r.mainCenteredMovableObjectsMutex.Unlock()
}

// GetMainCenteredMovableObject retrieves main centered movable object with the provided name.
func (r *Renderer) GetMainCenteredMovableObject(name string) *movable.Movable {
	r.mainCenteredMovableObjectsMutex.RLock()

	result, _ := r.mainCenteredMovableObjects[name]

	r.mainCenteredMovableObjectsMutex.RUnlock()

	return result
}

// SetSelectedObject sets selected object position.
func (r *Renderer) SetSelectedObject(value dto.Position) {
	r.selectedObjectEnabled = true
	r.selectedObjectPosition = value
}

// DisableSelectedObject disables selected object position.
func (r *Renderer) DisableSelectedObject() {
	r.selectedObjectEnabled = false
}

// Clean performs clean operation for the configured animator holders.
func (r *Renderer) Clean() {
	r.tertiaryTileObjectMutex.Lock()

	r.tertiaryTileObjects = orderedmap.NewOrderedMap[string, *tile.Tile]()

	r.tertiaryTileObjectMutex.Unlock()

	r.secondaryTileObjectMutex.Lock()

	r.secondaryTileObjects = orderedmap.NewOrderedMap[string, *tile.Tile]()

	r.secondaryTileObjectMutex.Unlock()

	r.secondaryExternalMovableObjectsMutex.Lock()

	clear(r.secondaryExternalMovableObjects)

	r.secondaryExternalMovableObjectsMutex.Unlock()

	r.mainCenteredMovableObjectsMutex.Lock()

	clear(r.mainCenteredMovableObjects)

	r.mainCenteredMovableObjectsMutex.Unlock()

	r.objectPositionMutex.Lock()

	r.objectPosition.Clear()

	r.objectPositionMutex.Unlock()

	r.selectedObjectPositionMutex.Lock()

	r.selectedObjectEnabled = false

	r.selectedObjectPositionMutex.Unlock()
}

// Update performs update operation and position rearangemenet for all the configured objects.
func (r *Renderer) Update(camera *kamera.Camera) {
	minCameraViewportWidth := store.GetPositionSession().X - (math.Abs(camera.CenterOffsetX) + CAMERA_RENDERING_OFFSET)
	maxCameraViewportWidth := store.GetPositionSession().X + (math.Abs(camera.CenterOffsetX) + CAMERA_RENDERING_OFFSET)

	minCameraViewportHeight := store.GetPositionSession().Y - (math.Abs(camera.CenterOffsetY) + CAMERA_RENDERING_OFFSET)
	maxCameraViewportHeight := store.GetPositionSession().Y + (math.Abs(camera.CenterOffsetY) + CAMERA_RENDERING_OFFSET)

	r.secondaryExternalMovableObjectsMutex.RLock()

	r.objectPosition.Clear()

	var (
		presentObjectPositions []dto.RendererPositionItem
		ok                     bool
	)

	for name, movable := range r.secondaryExternalMovableObjects {
		finalPosition := movable.GetFinalPositions()

		if (finalPosition.X < minCameraViewportWidth || finalPosition.X > maxCameraViewportWidth) ||
			(finalPosition.Y < minCameraViewportHeight || finalPosition.Y > maxCameraViewportHeight) {
			r.ignoredExternalMovableObjectsMutex.Lock()

			if _, ok := r.ignoredExternalMovableObjects[name]; !ok {
				r.ignoredExternalMovableObjects[name] = true
			}

			r.ignoredExternalMovableObjectsMutex.Unlock()

			continue
		}

		movable.Update()

		r.objectPositionMutex.RLock()

		shiftWidth, shiftHeight := movable.GetShiftBounds()

		position := movable.GetPosition()

		presentObjectPositions, ok = r.objectPosition.Get((position.X + (shiftWidth / 2)) + (position.Y + (shiftHeight)))
		if ok {
			presentObjectPositions = append(
				presentObjectPositions,
				dto.RendererPositionItem{
					Name: name,
					Type: dto.RendererPositionItemSecondaryExternalMovable})
		} else {
			presentObjectPositions = []dto.RendererPositionItem{
				dto.RendererPositionItem{
					Name: name,
					Type: dto.RendererPositionItemSecondaryExternalMovable}}
		}

		r.objectPositionMutex.RUnlock()

		r.objectPositionMutex.Lock()

		r.objectPosition.Set((position.X+(shiftWidth/2))+(position.Y+(shiftHeight)), presentObjectPositions)

		r.objectPositionMutex.Unlock()
	}

	r.secondaryExternalMovableObjectsMutex.RUnlock()

	r.secondaryLocalStaticObjectsMutex.Lock()

	for name, static := range r.secondaryLocalStaticObjects {
		if (static.GetPosition().X < minCameraViewportWidth || static.GetPosition().X > maxCameraViewportWidth) ||
			(static.GetPosition().Y < minCameraViewportHeight || static.GetPosition().Y > maxCameraViewportHeight) {
			continue
		}

		r.objectPositionMutex.RLock()

		position := static.GetPosition()

		shiftWidth, shiftHeight := static.GetShiftBounds()

		presentObjectPositions, ok = r.objectPosition.Get((position.X + (shiftWidth / 2)) + (position.Y + (shiftHeight)))
		if ok {
			presentObjectPositions = append(
				presentObjectPositions,
				dto.RendererPositionItem{
					Name: name,
					Type: dto.RendererPositionItemSecondaryStatic})
		} else {
			presentObjectPositions = []dto.RendererPositionItem{
				dto.RendererPositionItem{
					Name: name,
					Type: dto.RendererPositionItemSecondaryStatic}}
		}

		r.objectPositionMutex.RUnlock()

		r.objectPositionMutex.Lock()

		r.objectPosition.Set((position.X+(shiftWidth/2))+(position.Y+(shiftHeight)), presentObjectPositions)

		r.objectPositionMutex.Unlock()
	}

	r.secondaryLocalStaticObjectsMutex.Unlock()

	r.secondaryTileObjectMutex.RLock()

	for iter := r.secondaryTileObjects.Front(); iter != nil; iter = iter.Next() {
		if (iter.Value.GetPosition().X < minCameraViewportWidth || iter.Value.GetPosition().X > maxCameraViewportWidth) ||
			(iter.Value.GetPosition().Y < minCameraViewportHeight || iter.Value.GetPosition().Y > maxCameraViewportHeight) {
			continue
		}

		r.objectPositionMutex.RLock()

		position := iter.Value.GetPosition()

		shiftWidth, shiftHeight := iter.Value.GetShiftBounds()

		presentObjectPositions, ok = r.objectPosition.Get((position.X + (shiftWidth / 2)) + (position.Y + (shiftHeight)))
		if ok {
			presentObjectPositions = append(
				presentObjectPositions,
				dto.RendererPositionItem{
					Name: iter.Key,
					Type: dto.RendererPositionItemSecondaryTile})
		} else {
			presentObjectPositions = []dto.RendererPositionItem{
				dto.RendererPositionItem{
					Name: iter.Key,
					Type: dto.RendererPositionItemSecondaryTile}}
		}

		r.objectPositionMutex.RUnlock()

		r.objectPositionMutex.Lock()

		r.objectPosition.Set((position.X+(shiftWidth/2))+(position.Y+(shiftHeight)), presentObjectPositions)

		r.objectPositionMutex.Unlock()
	}

	r.secondaryTileObjectMutex.RUnlock()

	r.mainCenteredMovableObjectsMutex.RLock()

	for name, movable := range r.mainCenteredMovableObjects {
		movable.Update()

		r.objectPositionMutex.RLock()

		shiftWidth, shiftHeight := movable.GetShiftBounds()

		position := movable.GetPosition()

		presentObjectPositions, ok = r.objectPosition.Get((position.X + (shiftWidth / 2)) + (position.Y + (shiftHeight)))
		if ok {
			presentObjectPositions = append(
				presentObjectPositions,
				dto.RendererPositionItem{
					Name: name,
					Type: dto.RendererPositionItemMainCenteredMovable})
		} else {
			presentObjectPositions = []dto.RendererPositionItem{
				dto.RendererPositionItem{
					Name: name,
					Type: dto.RendererPositionItemMainCenteredMovable}}
		}

		r.objectPositionMutex.RUnlock()

		r.objectPositionMutex.Lock()

		r.objectPosition.Set((position.X+(shiftWidth/2))+(position.Y+(shiftHeight)), presentObjectPositions)

		r.objectPositionMutex.Unlock()
	}

	r.mainCenteredMovableObjectsMutex.RUnlock()
}

// Draw performs draw operation for all the configured objects.
func (r *Renderer) Draw(screen *ebiten.Image, camera *kamera.Camera) {
	minCameraViewportWidth := store.GetPositionSession().X - (math.Abs(camera.CenterOffsetX) + CAMERA_RENDERING_OFFSET)
	maxCameraViewportWidth := store.GetPositionSession().X + (math.Abs(camera.CenterOffsetX) + CAMERA_RENDERING_OFFSET)

	minCameraViewportHeight := store.GetPositionSession().Y - (math.Abs(camera.CenterOffsetY) + CAMERA_RENDERING_OFFSET)
	maxCameraViewportHeight := store.GetPositionSession().Y + (math.Abs(camera.CenterOffsetY) + CAMERA_RENDERING_OFFSET)

	for iter := r.tertiaryTileObjects.Front(); iter != nil; iter = iter.Next() {
		if (iter.Value.GetPosition().X >= minCameraViewportWidth && iter.Value.GetPosition().X <= maxCameraViewportWidth) &&
			(iter.Value.GetPosition().Y >= minCameraViewportHeight && iter.Value.GetPosition().Y <= maxCameraViewportHeight) {
			iter.Value.Draw(screen, false, camera)
		}
	}

	r.objectPositionMutex.RLock()

	r.objectPosition.Reverse(func(key float64, value []dto.RendererPositionItem) bool {
		for _, item := range value {
			var (
				selected              bool
				resetDelayedPositions bool
			)

			switch item.Type {
			case dto.RendererPositionItemSecondaryTile:
				value, ok := r.secondaryTileObjects.Get(item.Name)
				if !ok {
					continue
				}

				if r.selectedObjectEnabled && value.GetPosition() == r.selectedObjectPosition {
					selected = true
				}

				value.Draw(screen, selected, camera)

			case dto.RendererPositionItemSecondaryStatic:
				static, ok := r.secondaryLocalStaticObjects[item.Name]
				if !ok {
					continue
				}

				if r.selectedObjectEnabled && static.GetPosition() == r.selectedObjectPosition {
					selected = true
				}

				static.Draw(screen, selected, camera)

			case dto.RendererPositionItemSecondaryExternalMovable:
				movable, ok := r.secondaryExternalMovableObjects[item.Name]
				if !ok {
					continue
				}

				if r.selectedObjectEnabled && movable.GetPosition() == r.selectedObjectPosition {
					selected = true
				}

				r.ignoredExternalMovableObjectsMutex.Lock()

				if _, ok := r.ignoredExternalMovableObjects[item.Name]; ok {
					resetDelayedPositions = true

					delete(r.ignoredExternalMovableObjects, item.Name)
				}

				r.ignoredExternalMovableObjectsMutex.Unlock()

				movable.Draw(screen, resetDelayedPositions, selected, false, camera)

			case dto.RendererPositionItemMainCenteredMovable:
				movable, ok := r.mainCenteredMovableObjects[item.Name]
				if !ok {
					continue
				}

				if r.selectedObjectEnabled && movable.GetPosition() == r.selectedObjectPosition {
					selected = true
				}

				movable.Draw(screen, false, selected, true, camera)
			}
		}

		return true
	})

	r.objectPositionMutex.RUnlock()
}

// newRenderer initializes Renderer.
func newRenderer() *Renderer {
	return &Renderer{
		tertiaryTileObjects:             orderedmap.NewOrderedMap[string, *tile.Tile](),
		secondaryTileObjects:            orderedmap.NewOrderedMap[string, *tile.Tile](),
		secondaryLocalMovableObjects:    make(map[string]*movable.Movable),
		secondaryLocalStaticObjects:     make(map[string]*static.Static),
		secondaryExternalMovableObjects: make(map[string]*movable.Movable),
		mainCenteredMovableObjects:      make(map[string]*movable.Movable),
		objectPosition:                  btree.NewMap[float64, []dto.RendererPositionItem](32),
		ignoredExternalMovableObjects:   make(map[string]bool),
	}
}
