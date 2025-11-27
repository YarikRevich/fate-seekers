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
	collisionPolygonsMutex sync.Mutex

	// Represents collision polygons.
	collisionPolygons []*resolv.ConvexPolygon

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

// AddCollidableTileObject adds new collidable tile object with the provided value.
func (c *Collision) AddCollidableTileObject(value *dto.CollidableTile) {
	c.collisionPolygonsMutex.Lock()

	collider := resolv.NewConvexPolygon(
		value.Position.X-(float64(value.TileWidth)/2), value.Position.Y-(float64(value.TileHeight)/2),
		[]float64{
			float64(value.TileWidth) / 2.0, 0,
			float64(value.TileWidth), float64(value.TileHeight) / 2.0,
			float64(value.TileWidth) / 2.0, float64(value.TileHeight),
			0, float64(value.TileHeight) / 2.0,
		},
	)

	c.collisionPolygons = append(c.collisionPolygons, collider)

	c.collisionPolygonsMutex.Unlock()
}

// IsColliding checks if main trackable object is colliding with any other collidable object.
func (c *Collision) IsColliding() bool {
	c.collisionPolygonsMutex.Lock()

	for _, polygon := range c.collisionPolygons {
		if polygon.IsIntersecting(c.mainTrackableObject) {
			c.collisionPolygonsMutex.Unlock()

			return true
		}
	}

	c.collisionPolygonsMutex.Unlock()

	return false
}

// Clean performs clean operation for the configured collision holders.
func (c *Collision) Clean() {
	c.collisionPolygonsMutex.Lock()

	c.collisionPolygons = c.collisionPolygons[:0]

	c.collisionPolygonsMutex.Unlock()

	c.mainTrackableObjectMutex.Lock()

	c.mainTrackableObject = nil

	c.mainTrackableObjectMutex.Unlock()
}

// newCollision initializes Collision.
func newCollision() *Collision {
	return new(Collision)
}
