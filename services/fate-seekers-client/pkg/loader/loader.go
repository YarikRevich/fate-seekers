package loader

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// Describes all the available shaders to be loaded.
const ()

// Decsribes all the embedded files base pathes.
const (
	ShadersPath   = "dist/shaders"
	ObjectsPath   = "dist/objects"
	TemplatesPath = "dist/templates"
)

var (
	// GetInstance retrieves instance of the asset loader manager, performing initial creation if needed.
	GetInstance = sync.OnceValue[*Loader](newLoader)
)

// Loader represents asset loader manager, which operates in a lazy mode manner.
type Loader struct {
	// Represents cache map of embedded objects.
	objects map[string]*ebiten.Image

	// Represents cache map of embedded shaders.
	shaders map[string]*ebiten.Shader

	// Represents cache map of embedded templates.
	templates map[string][]byte
}

// GetObject retrieves object content with the given name.
func (l *Loader) GetObject(name string) *ebiten.Shader {
	return nil
}

// GetShader retrieves shader content with the given name.
func (l *Loader) GetShader(name string) *ebiten.Shader {
	return nil
}

// GetTemplate retrieves template content with the given name.
func (l *Loader) GetTemplate(name string) *ebiten.Shader {
	return nil
}

// newLoader initializes Loader.
func newLoader() *Loader {
	return &Loader{
		objects:   make(map[string]*ebiten.Image),
		shaders:   make(map[string]*ebiten.Shader),
		templates: make(map[string][]byte),
	}
}
