package star

import (
	"image/color"
	"math/rand"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// StarParticleElement represents a star particle element.
type StarParticleElement struct {
	// Represents X coordinate of the vector start point.
	fromX float32

	// Represents Y coordinate of the vector start point.
	fromY float32

	// Represents X coordinate of the vector end point.
	toX float32

	// Represents Y coordinate of the vector end point.
	toY float32

	// Represents brightness value of the vector.
	brightness float32

	// Represents division coefficient, which speeds up
	divider float32
}

// GetDivider retrieves divider value.
func (spe *StarParticleElement) GetDivider() float32 {
	return spe.divider
}

// SetDivider sets divider value.
func (spe *StarParticleElement) SetDivider(value float32) {
	spe.divider = value
}

// Reset performs particle state reset operation.
func (spe *StarParticleElement) Reset() {
	spe.toX = rand.Float32() * float32(config.GetWorldWidth())
	spe.fromX = spe.toX

	spe.toY = rand.Float32() * float32(config.GetWorldHeight())
	spe.fromY = spe.toY

	spe.brightness = rand.Float32() * 0xff
}

// Update performs particle state update operation.
func (spe *StarParticleElement) Update() {
	spe.fromX = spe.toX
	spe.fromY = spe.toY
	spe.toX += (spe.toX - float32(config.GetWorldWidth()/2)) / spe.divider
	spe.toY += (spe.toY - float32(config.GetWorldHeight()/2)) / spe.divider

	if spe.brightness < 255 {
		spe.brightness++
	}

	if spe.fromX < 0 || float32(config.GetWorldWidth()) < spe.fromX ||
		spe.fromY < 0 || float32(config.GetWorldHeight()) < spe.fromY {
		spe.Reset()
	}
}

// Draw performs draw operation for the particle vector.
func (spe *StarParticleElement) Draw(screen *ebiten.Image) {
	vector.StrokeLine(
		screen,
		spe.fromX,
		spe.fromY,
		spe.toX,
		spe.toY,
		2,
		&color.RGBA{
			R: uint8(187 * spe.brightness / 255),
			G: uint8(221 * spe.brightness / 255),
			B: uint8(255 * spe.brightness / 255),
			A: 255},
		true)
}

// NewStarParticleElement creates new StarParticleElement.
func NewStarParticleElement() *StarParticleElement {
	return new(StarParticleElement)
}
