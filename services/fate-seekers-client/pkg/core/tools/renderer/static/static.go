package static

import (
	"math"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/kamera/v2"
)

// Static represets static instance to be rendered at a certain renderer level.
type Static struct {
	// Represents static image for this exact renderer.
	image *ebiten.Image

	// Represents position of the image for this exact renderer.
	position dto.Position

	// Represents accumulated draw image options used for camera processing.
	opts ebiten.DrawImageOptions

	// Represents flags used for static object rendering management.
	flags map[string]bool
}

// GetPosition retrieves processed position.
func (s *Static) GetPosition() dto.Position {
	return s.position
}

// GetShiftBounds retrieves tile shift bounds.
func (s *Static) GetShiftBounds() (float64, float64) {
	shiftWidth := s.image.Bounds().Dx()
	shiftHeight := s.image.Bounds().Dy()

	return float64(shiftWidth), float64(shiftHeight)
}

// Draw performs draw operation for the static with the provided camera.
func (s *Static) Draw(screen *ebiten.Image, selected bool, camera *kamera.Camera) {
	s.opts.GeoM.Reset()

	s.opts.ColorM.Reset()

	s.opts.GeoM.Translate(s.position.X, -s.position.Y)

	if selected {
		ticks := float64(time.Now().UnixMilli()) / 200.0
		pulse := (math.Sin(ticks) + 1.0) / 2.0

		intensity := 0.3 * pulse

		s.opts.ColorM.Translate(intensity, intensity, intensity*0.5, 0)
	}

	camera.Draw(s.image, &s.opts, screen)
}

// NewStatic creates new Static instance.
func NewStatic(image *ebiten.Image, position dto.Position) *Static {
	return &Static{
		image:    image,
		position: position,
		flags:    make(map[string]bool),
	}
}
