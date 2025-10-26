package movable

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/animation"
	"github.com/hajimehoshi/ebiten"
)

// TODO: create a struct with animations created per some id.

type Movables struct {
	objects map[string]animation.AnimatorAnimationUnit
}

func (m *Movables) Clean() {
	clear(m.objects)
}

func (m *Movables) Get(name string) {

}

func (m *Movables) Add(name string, object animation.AnimatorAnimationUnit) {

}

func (m *Movables) Update() {
	for _, movable := range a.movables {
		movable.Update()
	}
}

func (m *Movables) Draw(screen *ebiten.Image) {
	for _, movable := range a.movables {
		movable.Draw(screen)
	}
}

func NewMovables() *Movables {
	return &Movables{
		objects: make(map[string]animation.AnimatorAnimationUnit),
	}
}
