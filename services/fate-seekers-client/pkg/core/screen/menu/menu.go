package menu

import (
	"sync"

	"github.com/Frabjous-Studios/asebiten"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition/transparent"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/camera"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/builder"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/letter"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/menu"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// GetInstance retrieves instance of the menu screen, performing initilization if needed.
	GetInstance = sync.OnceValue[screen.Screen](newMenuScreen)
)

// MenuScreen represents entry screen implementation.
type MenuScreen struct {
	// Represents attached user interface.
	ui *ebitenui.UI

	// Represents transparent transition effect.
	transparentTransitionEffect transition.TransitionEffect

	// Represents global world view.
	world *ebiten.Image

	animation *asebiten.Animation

	camera *camera.Camera
}

func (ms *MenuScreen) HandleInput() error {
	if !ms.transparentTransitionEffect.Done() {
		if !ms.transparentTransitionEffect.OnEnd() {
			ms.transparentTransitionEffect.Update()
		} else {
			ms.transparentTransitionEffect.Clean()
		}
	}

	ms.ui.Update()

	return nil
}

func (ms *MenuScreen) HandleNetworking() {

}

func (ms *MenuScreen) HandleRender(screen *ebiten.Image) {
	ms.world.Clear()

	// ms.animation.DrawTo(ms.world, &ebiten.DrawImageOptions{})

	// var shaderOpts ebiten.DrawRectShaderOptions

	// var g ebiten.GeoM
	// // g.Translate(float64(config.GetWorldWidth())/1.5, 0)

	// shaderOpts.GeoM = g

	// shaderOpts.Uniforms = make(map[string]interface{})
	// shaderOpts.Uniforms["Center"] = []float32{
	// 	float32(screen.Bounds().Dx()) / 2,
	// 	float32(screen.Bounds().Dy()) / 2,
	// }

	// screen.DrawRectShader(
	// 	screen.Bounds().Dx()/2,
	// 	screen.Bounds().Dy(),
	// 	loader.GetInstance().GetShader(loader.BasicTransitionShader),
	// 	&shaderOpts)

	// draw shader
	// indices := []uint16{0, 1, 2, 2, 1, 3} // map vertices to triangles
	// screen.(self.vertices[:], indices, self.shader, &shaderOpts)

	ms.ui.Draw(ms.world)

	screen.DrawImage(ms.world, &ebiten.DrawImageOptions{
		ColorM: ms.transparentTransitionEffect.GetOptions().ColorM})

	// worldX, worldY := r.camera.ScreenToWorld(ebiten.CursorPosition())
	// ebitenutil.DebugPrint(
	// 	screen,
	// 	fmt.Sprintf("TPS: %0.2f\nMove (WASD/Arrows)\nZoom (QE)\nRotate (R)\nReset (Space)", ebiten.ActualTPS()),
	// )

	// ebitenutil.DebugPrintAt(
	// 	screen,
	// 	fmt.Sprintf("%s\nCursor World Pos: %.2f,%.2f",
	// 		r.camera.String(),
	// 		worldX, worldY),
	// 	0, screenHeight-32,
	// )
}

func (ms *MenuScreen) Clean() {

}

// func loadImageNineSlice(path string, centerWidth int, centerHeight int) (*image.NineSlice, error) {
// 	i, err := newImageFromFile(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	w := i.Bounds().Dx()
// 	h := i.Bounds().Dy()
// 	return image.NewNineSlice(i,
// 			[3]int{(w - centerWidth) / 2, centerWidth, w - (w-centerWidth)/2 - centerWidth},
// 			[3]int{(h - centerHeight) / 2, centerHeight, h - (h-centerHeight)/2 - centerHeight}),
// 		nil
// }

// newMenuScreen initializes MenuScreen.
func newMenuScreen() screen.Screen {
	return &MenuScreen{
		ui:                          builder.Build(menu.NewMenuComponent(), letter.NewLetterComponent()),
		transparentTransitionEffect: transparent.NewTransparentTransitionEffect(),
		world:                       ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
		camera: camera.NewCamera(
			float64(config.GetWorldWidth()),
			float64(config.GetWorldHeight())),
		animation: loader.GetInstance().GetAnimation(loader.IntroSkullAnimation, false),
	}
}
