package logo

import (
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition/transparent"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/animation/combiner"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/options"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/logging"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/repository"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/repository/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/storage/shared"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// GetInstance retrieves instance of the logo screen, performing initilization if needed.
	GetInstance = sync.OnceValue[screen.Screen](newLogoScreen)
)

// LogoScreen represents logo screen implementation.
type LogoScreen struct {
	// Represents transparent transition effect.
	transparentTransitionEffect transition.TransitionEffect

	// Represents local logo animation.
	logoAnimation *combiner.AnimationCombiner

	// Represents global world view.
	world *ebiten.Image
}

func (ls *LogoScreen) HandleInput() error {
	if store.GetRepositoryUUIDChecked() == value.UUID_CHECKED_REPOSITORY_FALSE_VALUE {
		_, ok, err := repository.GetFlagsRepository().GetByName(common.UUID_FLAG_NAME)
		if err != nil {
			logging.GetInstance().Fatal(err.Error())
		}

		if !ok {
			uuidRaw := uuid.New().String()

			err = repository.GetFlagsRepository().InsertOrUpdate(common.UUID_FLAG_NAME, uuidRaw)
			if err != nil {
				logging.GetInstance().Fatal(err.Error())
			}

			dispatcher.GetInstance().Dispatch(
				action.NewSetUUIDRepositoryAction(uuidRaw))
		}

		dispatcher.GetInstance().Dispatch(
			action.NewSetUUIDCheckedRepositoryAction(value.UUID_CHECKED_REPOSITORY_TRUE_VALUE))
	}

	// TODO: check when logo animation is finished
	// _, ok, err := repository.GetFlagsRepository().GetByName(common.INTRO_FLAG_NAME)
	// if err != nil {
	// 	logging.GetInstance().Fatal(err.Error())
	// }

	// if !ok {
	// 	dispatcher.GetInstance().Dispatch(
	// 		action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_INTRO_VALUE))

	// 	err = repository.GetFlagsRepository().InsertOrUpdate(common.INTRO_FLAG_NAME, common.INTRO_FLAG_TRUE_VALUE)
	// 	if err != nil {
	// 		logging.GetInstance().Fatal(err.Error())
	// 	}
	// } else {
	// 	dispatcher.GetInstance().Dispatch(
	// 		action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_ENTRY_VALUE))
	// }

	dispatcher.GetInstance().Dispatch(
		action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))

	if !ls.transparentTransitionEffect.Done() {
		if !ls.transparentTransitionEffect.OnEnd() {
			ls.transparentTransitionEffect.Update()
		} else {
			ls.transparentTransitionEffect.Clean()
		}
	}

	shared.GetInstance().GetBackgroundAnimation().Update()

	return nil
}

func (ls *LogoScreen) HandleRender(screen *ebiten.Image) {
	var backgroundAnimationGeometry ebiten.GeoM

	backgroundAnimationGeometry.Scale(
		scaler.GetScaleFactor(config.GetMinStaticWidth(), config.GetWorldWidth()),
		scaler.GetScaleFactor(config.GetMinStaticHeight(), config.GetWorldHeight()))

	ls.logoAnimation.DrawTo(ls.world, &ebiten.DrawImageOptions{
		GeoM: backgroundAnimationGeometry,
	})

	screen.DrawImage(ls.world, &ebiten.DrawImageOptions{
		ColorM: options.GetTransparentDrawOptions(
			ls.transparentTransitionEffect.GetValue()).ColorM})
}

// newLogoScreen initializes LogoScreen.
func newLogoScreen() screen.Screen {
	return &LogoScreen{
		logoAnimation: combiner.NewAnimationCombiner(
			loader.GetInstance().GetAnimation(loader.Background1Animation, false),
			loader.GetInstance().GetAnimation(loader.Background2Animation, false),
			loader.GetInstance().GetAnimation(loader.Background3Animation, false),
			loader.GetInstance().GetAnimation(loader.Background4Animation, false),
			loader.GetInstance().GetAnimation(loader.Background5Animation, false),
			loader.GetInstance().GetAnimation(loader.Background6Animation, false),
		),
		transparentTransitionEffect: transparent.NewTransparentTransitionEffect(true, 255, 0, 5, time.Microsecond*10),
		world:                       ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
	}
}
