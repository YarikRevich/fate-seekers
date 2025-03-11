package toxicrain

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/shader"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/storage/shared"
	"github.com/hajimehoshi/ebiten/v2"
)

// ToxicRainEventEffect represents toxic rain event effect.
type ToxicRainEventEffect struct {
}

func (tree *ToxicRainEventEffect) Draw(screen *ebiten.Image) {
	screen.DrawRectShader(
		screen.Bounds().Dx(),
		screen.Bounds().Dy(),
		loader.GetInstance().GetShader(loader.ToxicRainShader),
		&ebiten.DrawRectShaderOptions{
			Uniforms: map[string]interface{}{
				"Time": shared.GetInstance().GetShaderTime(),
			},
		})
}

// NewToxicRainEventEffect initializes ToxicRainEventEffect.
func NewToxicRainEventEffect() shader.ShaderEffect {
	return new(ToxicRainEventEffect)
}
