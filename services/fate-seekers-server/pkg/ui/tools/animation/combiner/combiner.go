package combiner

import (
	"errors"

	"github.com/Frabjous-Studios/asebiten"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/logging"
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

// OnEnd checks if animation batch has been located on the end.
func (ac *AnimationCombiner) OnEnd() bool {
	if ac.currentAnimationIndex == len(ac.animations)-1 {
		return ac.animations[ac.currentAnimationIndex].FrameIdx() ==
			len(ac.animations[ac.currentAnimationIndex].Source.Frames)-1
	}

	return false
}

// GetFrameWidth retrieves current frame width
func (ac *AnimationCombiner) GetFrameWidth() int {
	return ac.animations[ac.currentAnimationIndex].Bounds().Dx()
}

// GetFrameHeight retrieves current frame height
func (ac *AnimationCombiner) GetFrameHeight() int {
	return ac.animations[ac.currentAnimationIndex].Bounds().Dy()
}

// Reset performs animation reset opereation.
func (ac *AnimationCombiner) Reset() {
	ac.animations[ac.currentAnimationIndex].Restart()

	ac.currentAnimationIndex = 0
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
