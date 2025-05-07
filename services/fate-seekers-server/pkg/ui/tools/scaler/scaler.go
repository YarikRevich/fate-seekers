package scaler

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/config"
	"github.com/hajimehoshi/ebiten/v2"
)

// GetScaleFactor retrieves scale factor for the given sizes.
func GetScaleFactor(currentSize int, expectedSize int) float64 {
	return float64(expectedSize) / float64(currentSize)
}

// GetPercentageOf retrieves given percentage of the given size.
func GetPercentageOf(size int, percentage int) int {
	return (size * percentage) / 100
}

// GetCenteredGeometry retrieves geometry with centered position for the provided image.
func GetCenteredGeometry(widthScaleFactor, heightScaleFactor, width, height int) ebiten.GeoM {
	var result ebiten.GeoM

	widthScale := GetScaleFactor(width, GetPercentageOf(config.GetWorldWidth(), widthScaleFactor))
	heightScale := GetScaleFactor(height, GetPercentageOf(config.GetWorldHeight(), heightScaleFactor))

	result.Scale(widthScale, heightScale)

	result.Translate(
		float64(GetPercentageOf(config.GetWorldWidth(), 50))-((float64(width)*widthScale)/2),
		float64(GetPercentageOf(config.GetWorldHeight(), 50))-(float64(height)*heightScale)/2)

	return result
}
