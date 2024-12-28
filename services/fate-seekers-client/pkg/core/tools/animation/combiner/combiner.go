package combiner

import (
	"errors"

	"github.com/Frabjous-Studios/asebiten"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/logging"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	ErrNotEnoughAnimations = errors.New("amount of given animations is not enough")
)

// AnimationCombiner represents holder of animations, which are executed as one.
type AnimationCombiner struct {
	// Represents all the available animations.
	animations []*asebiten.Animation

	// Represents currently selected animation.
	currentAnimationIndex int
}

// Update updates currently selected animation.
func (ac *AnimationCombiner) Update() {
	if ac.animations[ac.currentAnimationIndex].FrameIdx() ==
		len(ac.animations[ac.currentAnimationIndex].Source.Frames)-1 {
		ac.animations[ac.currentAnimationIndex].Restart()

		if ac.currentAnimationIndex == len(ac.animations)-1 {
			ac.currentAnimationIndex = 0
		} else {
			ac.currentAnimationIndex++
		}
	}

	ac.animations[ac.currentAnimationIndex].Update()
}

// DrawTo draws currently selected animation.
func (ac *AnimationCombiner) DrawTo(screen *ebiten.Image, options *ebiten.DrawImageOptions) {
	ac.animations[ac.currentAnimationIndex].DrawTo(screen, options)
}

// NewAnimationCombiner initializes AnimationCombiner.
func NewAnimationCombiner(animations ...*asebiten.Animation) *AnimationCombiner {
	if len(animations) == 0 {
		logging.GetInstance().Fatal(ErrNotEnoughAnimations.Error())
	}

	return &AnimationCombiner{
		animations:            animations,
		currentAnimationIndex: 0,
	}
}
