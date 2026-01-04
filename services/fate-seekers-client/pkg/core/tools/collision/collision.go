package collision

import (
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/solarlune/resolv"
)

var (
	// GetInstance retrieves instance of the sounder, performing initial creation if needed.
	GetInstance = sync.OnceValue[*Collision](newCollision)
)

// Collision represents active map collision manager.
type Collision struct {
	// Represents collision polygons mutex.
	collisionPolygonsMutex sync.RWMutex

	// Represents collision polygons.
	collisionPolygons map[string]*resolv.ConvexPolygon

	// Represents main trackable object mutex.
	mainTrackableObjectMutex sync.Mutex

	// Represents main trackable collision object.
	mainTrackableObject *resolv.ConvexPolygon
}

// SetMainTrackableObject sets main trackable object with the provided value.
func (c *Collision) SetMainTrackableObject(value dto.Position, shiftWidth, shiftHeight float64) {
	c.mainTrackableObjectMutex.Lock()

	c.mainTrackableObject = resolv.NewRectangle(
		value.X, value.Y, shiftWidth/2, shiftHeight/2)

	c.mainTrackableObjectMutex.Unlock()
}

// CollidablExists checks if collidable with the provided name exists.
func (c *Collision) CollidableExists(name string) bool {
	c.collisionPolygonsMutex.Lock()

	_, ok := c.collisionPolygons[name]

	c.collisionPolygonsMutex.Unlock()

	return ok
}

// addCollidable adds collidable object with the provided key.
func (c *Collision) addCollidable(key string, collider *resolv.ConvexPolygon) {
	c.collisionPolygonsMutex.Lock()

	c.collisionPolygons[key] = collider

	c.collisionPolygonsMutex.Unlock()
}

// AddCollidableTileObject adds new collidable tile object with the provided value.
func (c *Collision) AddCollidableTileObject(key string, value *dto.CollidableTile) {
	c.addCollidable(key, resolv.NewConvexPolygon(
		value.Position.X-(float64(value.TileWidth)/2), value.Position.Y-(float64(value.TileHeight)/2),
		[]float64{
			float64(value.TileWidth) / 2.0, 0,
			float64(value.TileWidth), float64(value.TileHeight) / 2.0,
			float64(value.TileWidth) / 2.0, float64(value.TileHeight),
			0, float64(value.TileHeight) / 2.0,
		},
	))
}

// AddCollidableStaticObject adds new collidable static object with the provided name and value.
func (c *Collision) AddCollidableStaticObject(name string, value *dto.CollidableStatic) {
	c.addCollidable(name, resolv.NewConvexPolygon(
		value.Position.X-(float64(value.TileWidth)/2), value.Position.Y-(float64(value.TileHeight)/2),
		[]float64{
			float64(value.TileWidth) / 2.0, 0,
			float64(value.TileWidth), float64(value.TileHeight) / 2.0,
			float64(value.TileWidth) / 2.0, float64(value.TileHeight),
			0, float64(value.TileHeight) / 2.0,
		},
	))
}

// GetCollidableStaticObject retrieves collidable static object with the provided name.
func (c *Collision) GetCollidableStaticObject(name string) *resolv.ConvexPolygon {
	c.collisionPolygonsMutex.RLock()

	result, _ := c.collisionPolygons[name]

	c.collisionPolygonsMutex.RUnlock()

	return result
}

// RemoveCollidableObject removes collidable object with the provided name and value.
func (c *Collision) RemoveCollidableObject(name string) {
	c.collisionPolygonsMutex.Lock()

	delete(c.collisionPolygons, name)

	c.collisionPolygonsMutex.Unlock()
}

// IsColliding checks if main trackable object is colliding with any other collidable object.
func (c *Collision) IsColliding() bool {
	c.collisionPolygonsMutex.Lock()

	c.mainTrackableObjectMutex.Lock()

	for _, polygon := range c.collisionPolygons {
		if polygon.IsIntersecting(c.mainTrackableObject) {
			c.collisionPolygonsMutex.Unlock()

			c.mainTrackableObjectMutex.Unlock()

			return true
		}
	}

	c.mainTrackableObjectMutex.Unlock()

	c.collisionPolygonsMutex.Unlock()

	return false
}

// Clean performs clean operation for the configured collision holders.
func (c *Collision) Clean() {
	c.collisionPolygonsMutex.Lock()

	clear(c.collisionPolygons)

	c.collisionPolygonsMutex.Unlock()

	c.mainTrackableObjectMutex.Lock()

	c.mainTrackableObject = nil

	c.mainTrackableObjectMutex.Unlock()
}

// newCollision initializes Collision.
func newCollision() *Collision {
	return &Collision{
		collisionPolygons: make(map[string]*resolv.ConvexPolygon),
	}
}
