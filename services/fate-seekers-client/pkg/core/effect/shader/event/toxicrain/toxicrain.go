package toxicrain

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/storage/shared"
	"github.com/hajimehoshi/ebiten/v2"
)

// ToxicRainEventEffect represents toxic rain event effect.
type ToxicRainEventEffect struct {
}

func (tree *ToxicRainEventEffect) Draw(screen *ebiten.Image, brightness float64) {
	screen.DrawRectShader(
		screen.Bounds().Dx(),
		screen.Bounds().Dy(),
		loader.GetInstance().GetShader(loader.ToxicRainShader),
		&ebiten.DrawRectShaderOptions{
			Uniforms: map[string]interface{}{
				"Time":       shared.GetInstance().GetShaderTime(),
				"Brightness": brightness,
			},
		})
}

// NewToxicRainEventEffect initializes ToxicRainEventEffect.
func NewToxicRainEventEffect() *ToxicRainEventEffect {
	return new(ToxicRainEventEffect)
}
