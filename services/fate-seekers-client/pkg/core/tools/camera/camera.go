package camera

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/math/f64"
)

// Camera represents implementation of camera.
type Camera struct {
	// Represents center position for the camera viewport.
	viewPortCenter f64.Vec2

	// Represents global position for the camera viewport.
	position f64.Vec2

	// Represents camera zoom value.
	zoom float64

	// Represents comera rotation value.
	rotation float64
}

// GetWorldMatrix retrieves calculated world matrix.
func (c *Camera) GetWorldMatrix() ebiten.GeoM {
	var matrix ebiten.GeoM

	matrix.Translate(-c.position[0], -c.position[1])

	matrix.Translate(-c.viewPortCenter[0], -c.viewPortCenter[1])

	matrix.Scale(math.Pow(1.01, c.zoom), math.Pow(1.01, c.zoom))

	matrix.Rotate(c.rotation * 2 * math.Pi / 360)

	matrix.Translate(c.viewPortCenter[0], c.viewPortCenter[1])

	return matrix
}

// ProjectPositionToWorld performs world matrix based projection using provided position.
func (c *Camera) ProjectPositionToWorld(x, y int) (float64, float64) {
	matrix := c.GetWorldMatrix()

	if matrix.IsInvertible() {
		matrix.Invert()

		return matrix.Apply(float64(x), float64(y))
	}

	return math.NaN(), math.NaN()
}

// SetPositionY sets camera horizontal position with the given value.
func (c *Camera) SetPositionX(x int) {
	c.position[0] = float64(x)
}

// SetPositionY sets camera vertical position with the given value.
func (c *Camera) SetPositionY(y int) {
	c.position[1] = float64(y)
}

// ZoomIn zooms in camera.
func (c *Camera) ZoomIn() {
	c.zoom++
}

// ZoomOut zooms out camera.
func (c *Camera) ZoomOut() {
	c.zoom--
}

// RotateRight performs right-based direction rotation.
func (c *Camera) RotateRight() {
	c.rotation++
}

// RotateRight performs left-based direction rotation.
func (c *Camera) RotateLeft() {
	c.rotation--
}

// Reset performs camera properties reset.
func (c *Camera) Reset() {
	c.position[0] = 0
	c.position[1] = 0
	c.rotation = 0
	c.zoom = 0
}

// NewCamera creates new instance of Camera. Requires target object width and height.
func NewCamera(width, height float64) *Camera {
	return &Camera{
		viewPortCenter: f64.Vec2{width * 0.5, height * 0.5},
	}
}
