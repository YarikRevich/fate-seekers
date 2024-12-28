package loader

import (
	"bytes"
	"image"
	"io/fs"
	"path/filepath"
	"sync"

	"github.com/YarikRevich/fate-seekers/assets"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/logging"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	_ "image/png"

	"github.com/Frabjous-Studios/asebiten"
)

var (
	ErrReadingFile      = errors.New("err happened during file read operation")
	ErrLoadingShader    = errors.New("err happened during shader loading operation")
	ErrLoadingFont      = errors.New("err happened during font loading operation")
	ErrLoadingStatic    = errors.New("err happened during image loading operation")
	ErrLoadingAnimation = errors.New("err happened during animation loading operation")
)

var (
	// GetInstance retrieves instance of the asset loader manager, performing initial creation if needed.
	GetInstance = sync.OnceValue[*Loader](newLoader)
)

// Describes all the available statics to be loaded.
const (
	ButtonIdleButton  = "button-idle.png"
	ButtonHoverButton = "button-hover.png"

	PanelIdlePanel = "panel-idle.png"

	ListDisabled = "list-disabled.png"
	ListIdle     = "list-idle.png"

	ListMask          = "list-mask.png"
	ListTrackDisabled = "list-track-disabled.png"
	ListTrackIdle     = "list-track-idle.png"

	SliderHandleHover = "slider-handle-hover.png"
	SliderHandleIdle  = "slider-handle-idle.png"
)

// Describes all the available shaders to be loaded.
const (
	BasicTransitionShader = "basic-transition.kage"
)

// Describes all the available fonts to be loaded.
const (
	KyivRegularFont = "kyiv-regular.ttf"
)

// Describes all the available templates to be loaded.
const ()

// Describes all the available animations to be loaded.
const (
	IntroSkullAnimation = "intro-skull/intro-skull.json"
	LogoAnimation       = "logo/logo.json"
	BackgroundAnimation = "background/background.json"
)

// Decsribes all the embedded files base pathes.
const (
	ShadersPath    = "dist/shaders"
	FontsPath      = "dist/fonts"
	ObjectsPath    = "dist/statics"
	TemplatesPath  = "dist/templates"
	AnimationsPath = "dist/animations"
)

// Loader represents low level asset loading manager, which operates in a lazy mode manner.
type Loader struct {
	// Represents cache map of embedded statics.
	statics sync.Map

	// Represents cache map of embedded shaders.
	shaders sync.Map

	// Represents cache map of embedded fonts.
	fonts sync.Map

	// Represents cache map of embedded templates.
	templates sync.Map

	// Represents cache map of embedded animations.
	animations sync.Map
}

// GetObject retrieves object content with the given name.
func (l *Loader) GetStatic(name string) *ebiten.Image {
	result, ok := l.statics.Load(name)
	if ok {
		return result.(*ebiten.Image)
	}

	file, err := fs.ReadFile(assets.Assets, filepath.Join(ObjectsPath, name))
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrReadingFile.Error()).Error())
	}

	source, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrLoadingStatic.Error()).Error())
	}

	image := ebiten.NewImageFromImage(source)

	l.statics.Store(name, image)

	logging.GetInstance().Debug("Static has been loaded", zap.String("name", name))

	return image
}

// GetShader retrieves shader content with the given name.
func (l *Loader) GetShader(name string) *ebiten.Shader {
	result, ok := l.shaders.Load(name)
	if ok {
		return result.(*ebiten.Shader)
	}

	file, err := fs.ReadFile(assets.Assets, filepath.Join(ShadersPath, name))
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrReadingFile.Error()).Error())
	}

	shader, err := ebiten.NewShader(file)
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrLoadingShader.Error()).Error())
	}

	l.shaders.Store(name, shader)

	logging.GetInstance().Debug("Shader has been loaded", zap.String("name", name))

	return shader
}

// GetFont retrieves font content with the given name.
func (l *Loader) GetFont(name string) *text.GoTextFaceSource {
	result, ok := l.fonts.Load(name)
	if ok {
		return result.(*text.GoTextFaceSource)
	}

	file, err := fs.ReadFile(assets.Assets, filepath.Join(FontsPath, name))
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrReadingFile.Error()).Error())
	}

	font, err := text.NewGoTextFaceSource(bytes.NewReader(file))
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrLoadingFont.Error()).Error())
	}

	l.fonts.Store(name, font)

	logging.GetInstance().Debug("Font has been loaded", zap.String("name", name))

	return font
}

// GetTemplate retrieves template content with the given name.
func (l *Loader) GetTemplate(name string) []byte {
	result, ok := l.templates.Load(name)
	if ok {
		return result.([]byte)
	}

	file, err := fs.ReadFile(assets.Assets, filepath.Join(TemplatesPath, name))
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrReadingFile.Error()).Error())
	}

	l.templates.Store(name, file)

	logging.GetInstance().Debug("Template has been loaded", zap.String("name", name))

	return file
}

// GetAnimation retrieves animation content with the given name. Allows to load new instance everytime.
// In order to load global object, it's required to set second argumeent to 'true'.
func (l *Loader) GetAnimation(name string, shared bool) *asebiten.Animation {
	if shared {
		result, ok := l.animations.Load(name)
		if ok {
			return result.(*asebiten.Animation)
		}
	}

	animation, err := asebiten.LoadAnimation(
		assets.Assets, filepath.Join(AnimationsPath, name))
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrLoadingAnimation.Error()).Error())
	}

	if shared {
		l.animations.Store(name, animation)

		logging.GetInstance().Debug("Animation has been loaded", zap.String("name", name))
	}

	return animation
}

// newLoader initializes Loader.
func newLoader() *Loader {
	return new(Loader)
}
