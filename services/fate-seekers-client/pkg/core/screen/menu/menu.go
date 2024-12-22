package menu

import (
	"sync"

	"github.com/Frabjous-Studios/asebiten"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/camera"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/builder"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/menu"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/application"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
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

	// Represents global world view.
	world *ebiten.Image

	animation *asebiten.Animation

	camera *camera.Camera
}

func (ms *MenuScreen) HandleInput() error {
	ms.ui.Update()

	if store.GetInstance().GetState(application.EXIT_APPLICATION_STATE) ==
		value.EXIT_APPLICATION_TRUE_VALUE {
		return ebiten.Termination
	}

	if ebiten.IsKeyPressed(ebiten.KeyU) {
		ms.animation.Update()
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

	var g ebiten.GeoM

	screen.DrawImage(ms.world, &ebiten.DrawImageOptions{GeoM: g})

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
	// // This creates the root container for this UI.
	// // All other UI elements must be added to this container.
	// rootContainer := widget.NewContainer()

	// // This adds the root container to the UI, so that it will be rendered.
	// eui := &ebitenui.UI{
	// 	Container: rootContainer,
	// }

	// fontFace := &text.GoTextFace{
	// 	Source: loader.GetInstance().GetFont(loader.KyivRegularFont),
	// 	Size:   128,
	// }

	// // This creates a text widget that says "Hello World!"
	// helloWorldLabel := widget.NewText(
	// 	widget.TextOpts.Text("Вітаю!", fontFace, color.White),
	// 	widget.TextOpts.WidgetOpts(widget.WidgetOpts.CursorEnterHandler(func(args *widget.WidgetCursorEnterEventArgs) { fmt.Println("pressed") })),
	// )

	// rootContainer.AddChild(c)

	// // To display the text widget, we have to add it to the root container.
	// rootContainer.AddChild(helloWorldLabel)

	return &MenuScreen{
		ui:    builder.Build(menu.NewMenuComponent()),
		world: ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
		camera: camera.NewCamera(
			float64(config.GetWorldWidth()),
			float64(config.GetWorldHeight())),
		animation: loader.GetInstance().GetAnimation(loader.IntroSkullAnimation, false),
	}
}
