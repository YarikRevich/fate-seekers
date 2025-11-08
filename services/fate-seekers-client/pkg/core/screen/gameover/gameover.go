package gameover

import (
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition/transparent"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/prompt"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/translation"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/storage/shared"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// GetInstance retrieves instance of the resume screen, performing initilization if needed.
	GetInstance = sync.OnceValue[screen.Screen](newDeathScreen)
)

// DeathScreen represents death screen implementation.
type DeathScreen struct {
	// Represents transparent transition effect.
	transparentTransitionEffect transition.TransitionEffect

	// Represents global world view.
	world *ebiten.Image
}

func (ds *DeathScreen) HandleInput() error {
	if store.GetResetDeath() == value.RESET_DEATH_TRUE_VALUE {
		dispatcher.GetInstance().Dispatch(
			action.NewSetResetDeath(value.RESET_DEATH_FALSE_VALUE))

		dispatcher.
			GetInstance().
			Dispatch(
				action.NewSetStateResetApplicationAction(
					value.STATE_RESET_APPLICATION_FALSE_VALUE))

		prompt.GetInstance().HideSubmitButton()

		dispatcher.GetInstance().Dispatch(
			action.NewSetPromptText(
				translation.GetInstance().GetTranslation("client.prompt.death")))

		dispatcher.GetInstance().Dispatch(
			action.NewSetPromptCancelCallback(func() {
				ds.transparentTransitionEffect.Reset()

				dispatcher.GetInstance().Dispatch(
					action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))
			}))
	}

	if !ds.transparentTransitionEffect.Done() {
		if !ds.transparentTransitionEffect.OnEnd() {
			ds.transparentTransitionEffect.Update()
		} else {
			ds.transparentTransitionEffect.Clean()
		}
	}

	shared.GetInstance().GetBackgroundAnimation().Update()

	return nil
}

func (ds *DeathScreen) HandleRender(screen *ebiten.Image) {
	ds.world.Clear()

	var backgroundAnimationGeometry ebiten.GeoM

	backgroundAnimationGeometry.Scale(
		scaler.GetScaleFactor(config.GetMinStaticWidth(), config.GetWorldWidth()),
		scaler.GetScaleFactor(config.GetMinStaticHeight(), config.GetWorldHeight()))

	shared.GetInstance().GetBackgroundAnimation().DrawTo(ds.world, &ebiten.DrawImageOptions{
		GeoM: backgroundAnimationGeometry,
	})

	screen.DrawImage(ds.world, &ebiten.DrawImageOptions{})
}

func newDeathScreen() screen.Screen {
	transparentTransitionEffect := transparent.NewTransparentTransitionEffect(true, 255, 0, 5, time.Microsecond*10)

	return &DeathScreen{
		transparentTransitionEffect: transparentTransitionEffect,
		world:                       ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
	}
}
