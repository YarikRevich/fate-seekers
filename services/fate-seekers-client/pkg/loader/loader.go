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
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	ErrReadingFile   = errors.New("err happened during file read operation")
	ErrLoadingShader = errors.New("err happened during shader loading operation")
	ErrLoadingImage  = errors.New("err happened during image loading operation")
)

var (
	// GetInstance retrieves instance of the asset loader manager, performing initial creation if needed.
	GetInstance = sync.OnceValue[*Loader](newLoader)
)

// Describes all the available objects to be loaded.
const ()

// Describes all the available shaders to be loaded.
const ()

// Describes all the available templates to be loaded.
const ()

// Decsribes all the embedded files base pathes.
const (
	ShadersPath   = "dist/shaders"
	ObjectsPath   = "dist/objects"
	TemplatesPath = "dist/templates"
)

// Loader represents asset loader manager, which operates in a lazy mode manner.
type Loader struct {
	// Represents cache map of embedded objects.
	objects sync.Map

	// Represents cache map of embedded shaders.
	shaders sync.Map

	// Represents cache map of embedded templates.
	templates sync.Map
}

// GetObject retrieves object content with the given name.
func (l *Loader) GetObject(name string) *ebiten.Image {
	result, ok := l.objects.Load(name)
	if ok {
		return result.(*ebiten.Image)
	}

	file, err := fs.ReadFile(assets.Assets, filepath.Join(ObjectsPath, name))
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrReadingFile.Error()).Error())
	}

	source, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrLoadingImage.Error()).Error())
	}

	image := ebiten.NewImageFromImage(source)

	l.objects.Store(name, image)

	logging.GetInstance().Debug("Image has been loaded", zap.String("name", name))

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

// newLoader initializes Loader.
func newLoader() *Loader {
	return new(Loader)
}
