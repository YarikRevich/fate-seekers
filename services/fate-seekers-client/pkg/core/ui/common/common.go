package common

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/ebitenui/ebitenui/image"
)

// func loadGraphicImages(idle string, disabled string) (*widget.ButtonImageImage, error) {
// 	idleImage, err := newImageFromFile(idle)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var disabledImage *ebiten.Image
// 	if disabled != "" {
// 		disabledImage, err = newImageFromFile(disabled)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	return &widget.ButtonImageImage{
// 		Idle:     idleImage,
// 		Disabled: disabledImage,
// 	}, nil
// }

// GetImageAsNineSlice retrieves image in a nine slice form.
func GetImageAsNineSlice(name string, centerWidth int, centerHeight int) *image.NineSlice {
	rawImage := loader.GetInstance().GetStatic(name)

	w := rawImage.Bounds().Dx()
	h := rawImage.Bounds().Dy()

	return image.NewNineSlice(rawImage,
		[3]int{(w - centerWidth) / 2, centerWidth, w - (w-centerWidth)/2 - centerWidth},
		[3]int{(h - centerHeight) / 2, centerHeight, h - (h-centerHeight)/2 - centerHeight})
}
