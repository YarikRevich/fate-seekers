package common

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/ebitenui/ebitenui/image"
)

// GetImageAsNineSlice retrieves image in a nine slice form.
func GetImageAsNineSlice(name string, centerWidth int, centerHeight int) *image.NineSlice {
	rawImage := loader.GetInstance().GetStatic(name)

	w := rawImage.Bounds().Dx()
	h := rawImage.Bounds().Dy()

	return image.NewNineSlice(rawImage,
		[3]int{(w - centerWidth) / 2, centerWidth, w - (w-centerWidth)/2 - centerWidth},
		[3]int{(h - centerHeight) / 2, centerHeight, h - (h-centerHeight)/2 - centerHeight})
}
