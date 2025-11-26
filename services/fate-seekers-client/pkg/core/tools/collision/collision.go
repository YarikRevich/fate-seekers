package collision

import (
	"math"
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
	collisionTileObjects []*dto.CollidableTile

	// Represents main trackable object mutex.
	mainTrackableObjectMutex sync.Mutex

	// Represents main trackable collision object.
	mainTrackableObject dto.Position

	// Represents main trackable object shift width.
	mainTrackableObjectShiftWidth float64

	// Represents main trackable object shift height.
	mainTrackableObjectShiftHeight float64
}

// SetMainTrackableObject sets main trackable object with the provided value.
func (c *Collision) SetMainTrackableObject(value dto.Position, shiftWidth, shiftHeight float64) {
	c.mainTrackableObjectMutex.Lock()

	c.mainTrackableObject = value

	c.mainTrackableObjectShiftWidth = shiftWidth

	c.mainTrackableObjectShiftHeight = shiftHeight

	c.mainTrackableObjectMutex.Unlock()
}

// AddCollidableTileObject adds new collidable tile object with the provided value.
func (c *Collision) AddCollidableTileObject(value *dto.CollidableTile) {
	c.collisionTileObjectMutex.Lock()

	c.collisionTileObjects = append(c.collisionTileObjects, value)

	c.collisionTileObjectMutex.Unlock()
}

// IsColliding checks if main trackable object is colliding with any other collidable object.
func (c *Collision) IsColliding() bool {
	c.collisionTileObjectMutex.Lock()

	defer c.collisionTileObjectMutex.Unlock()

	c.mainTrackableObjectMutex.Lock()

	defer c.mainTrackableObjectMutex.Unlock()

	// for _, value := range c.collisionTileObjects {
	// 	if c.mainTrackableObject.X < value.Position.X+float64(value.TileWidth/3) &&
	// 		c.mainTrackableObject.X+(c.mainTrackableObjectShiftWidth/2) > value.Position.X &&
	// 		c.mainTrackableObject.Y < value.Position.Y+float64(value.TileHeight/3) &&
	// 		c.mainTrackableObject.Y+(c.mainTrackableObjectShiftHeight/2) > value.Position.Y {
	// 		return true
	// 	}
	// }

	// Player Center Calculation
	playerCenterX := c.mainTrackableObject.X + (c.mainTrackableObjectShiftWidth / 2.0)
	playerCenterY := c.mainTrackableObject.Y + (c.mainTrackableObjectShiftHeight / 2.0)

	for _, tile := range c.collisionTileObjects {
		// --- PHASE 1: BROAD PHASE (RECTANGLE CHECK) ---
		// If the player isn't even touching the bounding box, skip the math.
		// This saves CPU cycles.
		if playerCenterX < tile.Position.X ||
			playerCenterX > tile.Position.X+float64(tile.TileWidth) ||
			playerCenterY < tile.Position.Y ||
			playerCenterY > tile.Position.Y+float64(tile.TileHeight) {
			continue // Skip to next tile
		}

		tileHalfW := float64(tile.TileWidth) / 2.0
		tileHalfH := float64(tile.TileHeight) / 2.0
		tileCenterX := tile.Position.X + tileHalfW
		tileCenterY := tile.Position.Y + tileHalfH

		dx := math.Abs(playerCenterX - tileCenterX)
		dy := math.Abs(playerCenterY - tileCenterY)

		// Check if inside the diamond using multiplication optimization
		if (dx*tileHalfH)+(dy*tileHalfW) <= (tileHalfW * tileHalfH) {
			return true
		}

		// // 1. Calculate Centers
		// // We need to work relative to the center of the tile, not the top-left.
		// tileHalfW := float64(tile.TileWidth) / 2.0
		// tileHalfH := float64(tile.TileHeight) / 2.0

		// tileCenterX := tile.Position.X + tileHalfW
		// tileCenterY := tile.Position.Y + tileHalfH

		// // Player center (The test point)
		// px := c.mainTrackableObject.X + (c.mainTrackableObjectShiftWidth / 2.0)
		// py := c.mainTrackableObject.Y + c.mainTrackableObjectShiftHeight // Checking feet position

		// // 2. Translate (Get point relative to tile center)
		// localX := px - tileCenterX
		// localY := py - tileCenterY

		// // 3. Un-Scale (Inverse of Scale(1, 0.5))
		// // We multiply Y by 2 to restore the aspect ratio to a perfect square.
		// localY = localY * 2.0

		// // 4. Un-Rotate (Inverse of Rotate(45 deg))
		// // We rotate by -45 degrees (negative Pi/4)
		// // Formula:
		// // x' = x * cos(θ) - y * sin(θ)
		// // y' = x * sin(θ) + y * cos(θ)

		// theta := -45.0 * (math.Pi / 180.0) // -0.785 radians
		// cosT := math.Cos(theta)
		// sinT := math.Sin(theta)

		// rotatedX := localX*cosT - localY*sinT
		// rotatedY := localX*sinT + localY*cosT

		// 5. The Check
		// Now that 'rotatedX' and 'rotatedY' are in "Square Space",
		// we calculate the size of that original square.
		// For a standard isometric tile (width W), the side length of the source square is (W / sqrt(2)).
		// But conceptually, we just check if the point is within the radius of the square's side.

		// Since we un-scaled, the "radius" of our local square is simply half the tile Width.
		// However, because of the 45-degree geometry, the boundary is actually:
		// originalSquareHalfSize := float64(tile.TileWidth) / 2.0 * math.Sqrt(0.5) // approx Width * 0.3535

		// If you strictly want the "Diamond" bounds we used before, the math simplifies to:
		// The point must be within the logical square bounds.

		// limit := tileHalfW * math.Sqrt(2) // This depends on how exactly you drew your assets.
		// // A safer, pure math limit for the bounding box of the unrotated square:
		// limit = tileHalfW

		// // Rectangular check on the rotated coordinates
		// if rotatedX >= -limit && rotatedX <= limit &&
		// 	rotatedY >= -limit && rotatedY <= limit {
		// 	return true
		// }
	}

	return false
}

// Clean performs clean operation for the configured collision holders.
func (c *Collision) Clean() {
	c.collisionTileObjectMutex.Lock()

	c.collisionTileObjects = c.collisionTileObjects[:0]

	c.collisionTileObjectMutex.Unlock()
}

// newCollision initializes Collision.
func newCollision() *Collision {
	return new(Collision)
}
