package entry

import (
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition/transparent"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/options"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/builder"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/subtitles"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/translation"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/logging"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/repository"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/repository/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/storage/shared"
	"github.com/ebitenui/ebitenui"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// GetInstance retrieves instance of the entry screen, performing initilization if needed.
	GetInstance = sync.OnceValue[screen.Screen](newEntryScreen)
)

// EntryScreen represents entry screen implementation.
type EntryScreen struct {
	// Represents attached user interface.
	ui *ebitenui.UI

	// Represents transparent transition effect.
	transparentTransitionEffect transition.TransitionEffect

	// Represents global world view.
	world *ebiten.Image

	// Represents stub timer, which is used to emulate some delay before screne switch.
	stubTimer *time.Timer
}

func (es *EntryScreen) HandleInput() error {
	if store.GetRepositoryUUIDChecked() == value.UUID_CHECKED_REPOSITORY_FALSE_VALUE {
		flag, ok, err := repository.GetFlagsRepository().GetByName(common.UUID_FLAG_NAME)
		if err != nil {
			logging.GetInstance().Fatal(err.Error())
		}

		if ok {
			dispatcher.GetInstance().Dispatch(
				action.NewSetUUIDRepositoryAction(flag.Value))
		} else {
			uuidRaw := uuid.NewString()

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

	if !es.transparentTransitionEffect.Done() {
		if !es.transparentTransitionEffect.OnEnd() {
			es.transparentTransitionEffect.Update()
		} else {
			es.transparentTransitionEffect.Clean()
		}
	}

	shared.GetInstance().GetBackgroundAnimation().Update()

	subtitles.GetInstance().Push(
		translation.GetInstance().GetTranslation("client.entry.welcome"),
		time.Second*3)

	dispatcher.GetInstance().Dispatch(
		action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))

	es.ui.Update()

	return nil
}

func (es *EntryScreen) HandleRender(screen *ebiten.Image) {
	var backgroundAnimationGeometry ebiten.GeoM

	backgroundAnimationGeometry.Scale(
		scaler.GetScaleFactor(config.GetMinStaticWidth(), config.GetWorldWidth()),
		scaler.GetScaleFactor(config.GetMinStaticHeight(), config.GetWorldHeight()))

	shared.GetInstance().GetBackgroundAnimation().DrawTo(es.world, &ebiten.DrawImageOptions{
		GeoM: backgroundAnimationGeometry,
	})

	es.ui.Draw(es.world)

	screen.DrawImage(es.world, &ebiten.DrawImageOptions{
		ColorM: options.GetTransparentDrawOptions(
			es.transparentTransitionEffect.GetValue()).ColorM})
}

// newEntryScreen initializes EntryScreen.
func newEntryScreen() screen.Screen {
	stubTimer := time.NewTimer(time.Minute)

	stubTimer.Stop()

	return &EntryScreen{
		ui:                          builder.Build(),
		transparentTransitionEffect: transparent.NewTransparentTransitionEffect(true, 255, 0, 5, time.Microsecond*10),
		world:                       ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
		stubTimer:                   stubTimer,
	}
}
