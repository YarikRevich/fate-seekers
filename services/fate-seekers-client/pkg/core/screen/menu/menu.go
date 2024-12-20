package menu

import (
	"sync"

	"github.com/Frabjous-Studios/asebiten"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/camera"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// GetInstance retrieves instance of the menu screen, performing initilization if needed.
	GetInstance = sync.OnceValue[screen.Screen](newMenuScreen)
)

// MenuScreen represents entry screen implementation.
type MenuScreen struct {
	// Represents global world view.
	world *ebiten.Image

	animation *asebiten.Animation

	camera *camera.Camera

	r float64
}

func (ms *MenuScreen) HandleInput() {
	if ebiten.IsKeyPressed(ebiten.KeyU) {
		ms.animation.Update()
		ms.r++
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		ms.camera.TranslatePositionX(-1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		ms.camera.TranslatePositionX(1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		ms.camera.TranslatePositionY(-1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		ms.camera.TranslatePositionY(1)
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		if ms.camera.GetZoom() > -2400 {
			ms.camera.ZoomOutBy(10)
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		if ms.camera.GetZoom() < 2400 {
			ms.camera.ZoomInBy(10)
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		ms.camera.RotateLeft()
	}

	if ebiten.IsKeyPressed(ebiten.KeyT) {
		ms.camera.RotateRight()
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		ms.camera.Reset()
	}
}

func (ms *MenuScreen) HandleNetworking() {

}

func (ms *MenuScreen) HandleRender(screen *ebiten.Image) {
	ms.world.Clear()

	ms.animation.DrawTo(ms.world, &ebiten.DrawImageOptions{})

	var shaderOpts ebiten.DrawRectShaderOptions
	shaderOpts.Uniforms = make(map[string]interface{})
	shaderOpts.Uniforms["Center"] = []float32{
		float32(screen.Bounds().Dx()) / 2,
		float32(screen.Bounds().Dy()) / 2,
	}
	shaderOpts.Uniforms["Radius"] = float32(ms.r)

	screen.DrawRectShader(
		screen.Bounds().Dx(),
		screen.Bounds().Dy(),
		loader.GetInstance().GetShader(loader.BasicTransitionShader),
		&shaderOpts)

	// draw shader
	// indices := []uint16{0, 1, 2, 2, 1, 3} // map vertices to triangles
	// screen.(self.vertices[:], indices, self.shader, &shaderOpts)

	screen.DrawImage(ms.world, &ebiten.DrawImageOptions{
		GeoM: ms.camera.GetWorldMatrix(),
	})

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

// newMenuScreen initializes MenuScreen.
func newMenuScreen() screen.Screen {
	return &MenuScreen{
		world: ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
		camera: camera.NewCamera(
			float64(config.GetWorldWidth()),
			float64(config.GetWorldHeight())),
		animation: loader.GetInstance().GetAnimation(loader.IntroSkullAnimation, false),
		r:         80,
	}
}
