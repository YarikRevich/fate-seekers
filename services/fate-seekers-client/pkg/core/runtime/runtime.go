package runtime

import (
	"github.com/Frabjous-Studios/asebiten"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/camera"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/hajimehoshi/ebiten/v2"
)

// const (
// 	screenWidth  = 480
// 	screenHeight = 320
// )

// const (
// 	tileSize   = 16
// 	tileXCount = 25
// )

const (
	worldWidth  = 480
	worldHeight = 320
	// worldSizeX  = worldWidth / tileSize
)

// var (
// 	tilesImage *ebiten.Image
// )

// func init() {
// 	// Decode an image from the image file's byte slice.
// 	img, _, err := image.Decode(bytes.NewReader(images.Tiles_png))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	tilesImage = ebiten.NewImageFromImage(img)
// }

// type Camera struct {
// 	ViewPort   f64.Vec2
// 	Position   f64.Vec2
// 	ZoomFactor int
// 	Rotation   int
// }

// func (c *Camera) String() string {
// 	return fmt.Sprintf(
// 		"T: %.1f, R: %d, S: %d",
// 		c.Position, c.Rotation, c.ZoomFactor,
// 	)
// }

// func (c *Camera) viewportCenter() f64.Vec2 {
// 	return f64.Vec2{
// 		c.ViewPort[0] * 0.5,
// 		c.ViewPort[1] * 0.5,
// 	}
// }

// func (c *Camera) worldMatrix() ebiten.GeoM {
// 	m := ebiten.GeoM{}
// 	m.Translate(-c.Position[0], -c.Position[1])
// 	// We want to scale and rotate around center of image / screen
// 	m.Translate(-c.viewportCenter()[0], -c.viewportCenter()[1])

// 	m.Scale(
// 		math.Pow(1.01, float64(c.ZoomFactor)),
// 		math.Pow(1.01, float64(c.ZoomFactor)),
// 	)
// 	m.Rotate(float64(c.Rotation) * 2 * math.Pi / 360)
// 	m.Translate(c.viewportCenter()[0], c.viewportCenter()[1])
// 	return m
// }

// func (c *Camera) Render(world, screen *ebiten.Image) {
// 	screen.DrawImage(world, &ebiten.DrawImageOptions{
// 		GeoM: c.worldMatrix(),
// 	})
// }

// func (c *Camera) ScreenToWorld(posX, posY int) (float64, float64) {
// 	inverseMatrix := c.worldMatrix()
// 	if inverseMatrix.IsInvertible() {
// 		inverseMatrix.Invert()
// 		return inverseMatrix.Apply(float64(posX), float64(posY))
// 	} else {
// 		// When scaling it can happened that matrix is not invertable
// 		return math.NaN(), math.NaN()
// 	}
// }

// func (c *Camera) Reset() {
// 	c.Position[0] = 0
// 	c.Position[1] = 0
// 	c.Rotation = 0
// 	c.ZoomFactor = 0
// }

type Runtime struct {
	layers [][]int
	world  *ebiten.Image
	camera *camera.Camera

	animation *asebiten.Animation
}

func (r *Runtime) Update() error {
	switch store.GetActiveScreen() {
	case value.ACTIVE_SCREEN_ENTRY_VALUE:
	case value.ACTIVE_SCREEN_MENU_VALUE:
	}

	if ebiten.IsKeyPressed(ebiten.KeyU) {
		r.animation.Update()
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		r.camera.TranslatePositionX(-1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		r.camera.TranslatePositionX(1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		r.camera.TranslatePositionY(-1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		r.camera.TranslatePositionY(1)
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		if r.camera.GetZoom() > -2400 {
			r.camera.ZoomOut()
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		if r.camera.GetZoom() < 2400 {
			r.camera.ZoomIn()
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		r.camera.RotateLeft()
	}

	if ebiten.IsKeyPressed(ebiten.KeyT) {
		r.camera.RotateRight()
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		r.camera.Reset()
	}

	return nil
}

func (r *Runtime) Draw(screen *ebiten.Image) {
	// for _, l := range r.layers {
	// 	for i, t := range l {
	// 		op := &ebiten.DrawImageOptions{}
	// 		op.GeoM.Translate(float64((i%worldSizeX)*tileSize), float64((i/worldSizeX)*tileSize))

	// 		sx := (t % tileXCount) * tileSize
	// 		sy := (t / tileXCount) * tileSize
	// 		r.world.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
	// 	}
	// 	}

	r.world.Clear()

	r.animation.DrawTo(r.world, &ebiten.DrawImageOptions{})

	// r.world.DrawImage(loader.GetInstance().GetStatic("test.png"), &ebiten.DrawImageOptions{})

	screen.DrawImage(r.world, &ebiten.DrawImageOptions{
		GeoM: r.camera.GetWorldMatrix(),
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

func (r *Runtime) Layout(outsideWidth, outsideHeight int) (int, int) {
	return worldWidth, worldHeight
}

func NewRuntime() *Runtime {
	ebiten.SetWindowSize(1920, 1080)
	ebiten.SetWindowTitle("Shader (Ebitengine Demo)")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetVsyncEnabled(true)

	f := loader.GetInstance().GetAnimation(loader.IntroSkullAnimation, false)

	g := &Runtime{
		layers: [][]int{
			{
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 218, 243, 243, 243, 243, 243,
				243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 218, 243, 243, 243, 243, 243, 243, 244, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,

				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 244, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 218, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 219, 243, 243, 243, 219, 243, 243, 243, 243, 243, 243, 243, 218, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,

				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 244, 243, 243, 243, 243, 243, 243, 243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,

				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 218, 243, 243, 243,
				243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 244, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
			},
			{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 26, 27, 28, 29, 30, 31, 0, 0, 0, 0, 0, 0, 0, 0, 0, 63, 64, 65, 66, 67, 68, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 51, 52, 53, 54, 55, 56, 0, 0, 0, 0, 0, 0, 0, 0, 0, 88, 89, 90, 91, 92, 93, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 76, 77, 78, 79, 80, 81, 0, 0, 0, 0, 0, 0, 0, 0, 0, 113, 114, 115, 116, 117, 118, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 101, 102, 103, 104, 105, 106, 0, 0, 0, 0, 0, 0, 0, 0, 0, 138, 139, 140, 141, 142, 143, 0, 0, 0, 0,

				0, 0, 0, 0, 0, 126, 127, 128, 129, 130, 131, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 245, 242, 288, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 303, 303, 245, 242, 303, 303, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0,

				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 45, 46, 47, 48,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 70, 71, 72, 73,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 95, 96, 97, 98,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 120, 121, 122, 123,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 145, 146, 147, 148,

				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 267, 268, 268, 268, 268, 268, 268, 268, 268, 268, 268, 268, 268, 268, 268, 270, 242, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 192, 193, 193, 193, 193, 193, 193, 193, 193, 193, 193, 193, 193, 193, 193, 193, 222, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			},
		},
		camera:    camera.NewCamera(worldWidth, worldHeight),
		animation: f,
	}
	g.world = ebiten.NewImage(worldWidth, worldHeight)

	return g
}
