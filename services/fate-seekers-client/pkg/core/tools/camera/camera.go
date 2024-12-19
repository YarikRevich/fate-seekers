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

// TranslatePositionX translates camera horizontal position by the given value.
func (c *Camera) TranslatePositionX(x float64) {
	c.position[0] += x
}

// TranslatePositionY translates camera vertical position by the given value.
func (c *Camera) TranslatePositionY(y float64) {
	c.position[1] += y
}

// ZoomIn zooms in camera.
func (c *Camera) ZoomIn() {
	c.zoom++
}

// ZoomInBy zooms in camera by the given value.
func (c *Camera) ZoomInBy(value float64) {
	c.zoom += value
}

// ZoomOut zooms out camera.
func (c *Camera) ZoomOut() {
	c.zoom--
}

// ZoomOutBy zooms out camera by the given value.
func (c *Camera) ZoomOutBy(value float64) {
	c.zoom -= value
}

// GetZoom retrieves camera zoom value.
func (c *Camera) GetZoom() float64 {
	return c.zoom
}

// RotateRight performs right-based direction rotation.
func (c *Camera) RotateRight() {
	c.rotation++
}

// RotateRightBy performs right-based direction rotation by the given value.
func (c *Camera) RotateRightBy(value float64) {
	c.rotation += value
}

// RotateRight performs left-based direction rotation.
func (c *Camera) RotateLeft() {
	c.rotation--
}

// RotateRightBy performs left-based direction rotation by the given value.
func (c *Camera) RotateLeftBy(value float64) {
	c.rotation -= value
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
