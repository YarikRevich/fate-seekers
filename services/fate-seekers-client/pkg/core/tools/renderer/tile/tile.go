package tile

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/kamera/v2"
)

// Tile represets tile instance to be rendered at a certain renderer level.
type Tile struct {
	processed *dto.ProcessedTile

	// Represents accumulated draw image options used for camera processing.
	opts ebiten.DrawImageOptions
}

// Draw performs draw operation for the tile with the provided camera.
func (t *Tile) Draw(screen *ebiten.Image, camera *kamera.Camera) {
	t.opts.GeoM.Reset()

	t.opts.GeoM.Translate(t.processed.Position.X, -t.processed.Position.Y)

	camera.Draw(t.processed.Image, &t.opts, screen)
}

// NewTile creates new Tile instance.
func NewTile(processed *dto.ProcessedTile) *Tile {
	return &Tile{
		processed: processed,
	}
}
