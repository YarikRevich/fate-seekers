package entry

import (
	"image/color"
	"math/rand"
	"sync"
	"time"

	"github.com/Frabjous-Studios/asebiten"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/camera"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	// GetInstance retrieves instance of the entry screen, performing initilization if needed.
	GetInstance = sync.OnceValue[screen.Screen](newEntryScreen)
)

const (
	scale      = 64
	starsCount = 256
)

type Star struct {
	fromx, fromy, tox, toy, brightness float32
}

func (s *Star) Init() {
	s.tox = rand.Float32() * float32(config.GetWorldWidth()) * scale
	s.fromx = s.tox
	s.toy = rand.Float32() * float32(config.GetWorldHeight()) * scale
	s.fromy = s.toy
	s.brightness = rand.Float32() * 0xff
}

func (s *Star) Update() {
	s.fromx = s.tox
	s.fromy = s.toy
	s.tox += (s.tox - float32(config.GetWorldWidth()/2*scale)) / 32
	s.toy += (s.toy - float32(config.GetWorldHeight()/2*scale)) / 32
	s.brightness += 1
	if 0xff < s.brightness {
		s.brightness = 0xff
	}
	if s.fromx < 0 || float32(config.GetWorldWidth())*scale < s.fromx || s.fromy < 0 || float32(config.GetWorldHeight())*scale < s.fromy {
		s.Init()
	}
}

func (s *Star) Draw(screen *ebiten.Image) {
	c := color.RGBA{
		R: uint8(0xbb * s.brightness / 0xff),
		G: uint8(0xdd * s.brightness / 0xff),
		B: uint8(0xff * s.brightness / 0xff),
		A: 0xff}
	vector.StrokeLine(screen, s.fromx/scale, s.fromy/scale, s.tox/scale, s.toy/scale, 1, c, true)
}

// EntryScreen represents entry screen implementation.
type EntryScreen struct {
	stars [starsCount]Star

	// Represents global world view.
	world *ebiten.Image

	//
	camera *camera.Camera

	//
	cameraTicker *time.Ticker

	//
	logoAnimation *asebiten.Animation
}

func (es *EntryScreen) HandleInput() error {
	es.logoAnimation.Update()

	for i := 0; i < starsCount; i++ {
		es.stars[i].Update()
	}

	select {
	case <-es.cameraTicker.C:
		if es.camera.GetZoom()-2 > 20 {
			es.camera.ZoomOutBy(2)
		} else {
			es.cameraTicker.Stop()
		}
	default:
	}

	return nil
}

func (es *EntryScreen) HandleNetworking() {
	// TODO: do initial connection
}

func (es *EntryScreen) HandleRender(screen *ebiten.Image) {
	es.world.Fill(color.RGBA{R: 31, G: 62, B: 90})

	for i := 0; i < starsCount; i++ {
		es.stars[i].Draw(es.world)
	}

	var logoAnimationGeometry ebiten.GeoM

	logoAnimationGeometry.Translate(
		float64(config.GetWorldWidth())/16.5, float64(config.GetWorldHeight())/4)

	es.logoAnimation.DrawTo(es.world, &ebiten.DrawImageOptions{
		GeoM: logoAnimationGeometry})

	screen.DrawImage(es.world, &ebiten.DrawImageOptions{
		GeoM: es.camera.GetWorldMatrix(),
	})
}

func (es *EntryScreen) Clean() {
}

// newEntryScreen initializes EntryScreen.
func newEntryScreen() screen.Screen {
	camera := camera.NewCamera(
		float64(config.GetWorldWidth()), float64(config.GetWorldHeight()))

	camera.ZoomInBy(515)

	r := &EntryScreen{
		world:         ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
		camera:        camera,
		cameraTicker:  time.NewTicker(time.Millisecond * 30),
		logoAnimation: loader.GetInstance().GetAnimation(loader.LogoAnimation, false),
	}

	for i := 0; i < starsCount; i++ {
		r.stars[i].Init()
	}

	return r
}
