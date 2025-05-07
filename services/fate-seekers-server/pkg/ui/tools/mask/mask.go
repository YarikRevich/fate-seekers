package mask

import (
	"errors"
	"image/color"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/core/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/logging"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	ErrValueOutOfRange = errors.New("err value is out of range")
)

// GetMaskEffect retrieves image mask draw options.
func GetMaskEffect(value int) *ebiten.DrawImageOptions {
	if value > 100 && value < 0 {
		logging.GetInstance().Fatal(ErrValueOutOfRange.Error())
	}

	var c ebiten.ColorM

	c.ScaleWithColor(color.RGBA{R: 255, G: 255, B: 255, A: uint8(scaler.GetPercentageOf(255, value))})

	return &ebiten.DrawImageOptions{ColorM: c}
}
