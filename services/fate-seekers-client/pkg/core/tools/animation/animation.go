package animation

import "github.com/hajimehoshi/ebiten/v2"

// AnimatorAnimationUnit represents animator animation unit interface.
type AnimatorAnimationUnit interface {
	// Clean performes forced memory cleanup for the animation only.
	Clean()

	// Update performs update operation for all animations.
	Update()

	// Draw performs draw operation for all animations.
	Draw(screen *ebiten.Image)
}
