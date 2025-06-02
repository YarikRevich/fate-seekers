package options

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// GetTransparentDrawOptions retrieves transparent draw options with the provided alpha value.
func GetTransparentDrawOptions(value float64) *ebiten.DrawImageOptions {
	var c ebiten.ColorM

	c.ScaleWithColor(color.RGBA{R: 255, G: 255, B: 255, A: uint8(value)})

	return &ebiten.DrawImageOptions{ColorM: c}
}
