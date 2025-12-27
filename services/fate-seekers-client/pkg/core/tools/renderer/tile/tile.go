package tile

import (
	"math"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/kamera/v2"
)

// Tile represets tile instance to be rendered at a certain renderer level.
type Tile struct {
	// Represents processed tile selected for this exact renderer.
	processed *dto.ProcessedTile

	// Represents accumulated draw image options used for camera processing.
	opts ebiten.DrawImageOptions
}

// GetPosition retrieves processed position.
func (t *Tile) GetPosition() dto.Position {
	return t.processed.Position
}

// GetShiftBounds retrieves tile shift bounds.
func (t *Tile) GetShiftBounds() (float64, float64) {
	shiftWidth := t.processed.Image.Bounds().Dx()
	shiftHeight := t.processed.Image.Bounds().Dy()

	return float64(shiftWidth), float64(shiftHeight)
}

// Draw performs draw operation for the tile with the provided camera.
func (t *Tile) Draw(screen *ebiten.Image, selected bool, camera *kamera.Camera) {
	t.opts.GeoM.Reset()

	t.opts.ColorM.Reset()

	t.opts.GeoM.Translate(t.processed.Position.X, -t.processed.Position.Y)

	if selected {
		ticks := float64(time.Now().UnixMilli()) / 200.0
		pulse := (math.Sin(ticks) + 1.0) / 2.0

		intensity := 0.3 * pulse

		t.opts.ColorM.Translate(intensity, intensity, intensity*0.5, 0)
	}

	camera.Draw(t.processed.Image, &t.opts, screen)
}

// NewTile creates new Tile instance.
func NewTile(processed *dto.ProcessedTile) *Tile {
	return &Tile{
		processed: processed,
	}
}
