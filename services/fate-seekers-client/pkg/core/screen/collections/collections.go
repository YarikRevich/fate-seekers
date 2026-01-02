package collections

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
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/collections"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/notification"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/translation"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/repository"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/repository/converter"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/storage/shared"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// GetInstance retrieves instance of the collections screen, performing initilization if needed.
	GetInstance = sync.OnceValue[screen.Screen](newCollectionsScreen)
)

// CollectionsScreen represents collections screen implementation.
type CollectionsScreen struct {
	// Represents attached user interface.
	ui *ebitenui.UI

	// Represents transparent transition effect.
	transparentTransitionEffect transition.TransitionEffect

	// Represents global world view.
	world *ebiten.Image

	// Represents interface world view.
	interfaceWorld *ebiten.Image
}

func (cs *CollectionsScreen) HandleInput() error {
	if store.GetResetCollections() == value.RESET_COLLECTIONS_FALSE_VALUE {
		dispatcher.GetInstance().Dispatch(
			action.NewSetResetCollections(value.RESET_COLLECTIONS_TRUE_VALUE))

		items, err := repository.GetCollectionsRepository().GetAll()
		if err != nil {
			notification.GetInstance().Push(
				common.ComposeMessage(
					translation.GetInstance().GetTranslation("client.repository.collections-retrieval-failure"),
					err.Error()),
				time.Second*3,
				common.NotificationErrorTextColor)

		} else {
			collections.GetInstance().SetListsEntries(
				converter.ConvertRetrievedCollectionsToListEntries(items))
		}
	}

	if !cs.transparentTransitionEffect.Done() {
		if !cs.transparentTransitionEffect.OnEnd() {
			cs.transparentTransitionEffect.Update()
		} else {
			cs.transparentTransitionEffect.Clean()
		}
	}

	shared.GetInstance().GetBackgroundAnimation().Update()

	cs.ui.Update()

	return nil
}

func (cs *CollectionsScreen) HandleRender(screen *ebiten.Image) {
	cs.world.Clear()

	cs.interfaceWorld.Clear()

	var backgroundAnimationGeometry ebiten.GeoM

	backgroundAnimationGeometry.Scale(
		scaler.GetScaleFactor(config.GetMinStaticWidth(), config.GetWorldWidth()),
		scaler.GetScaleFactor(config.GetMinStaticHeight(), config.GetWorldHeight()))

	shared.GetInstance().GetBackgroundAnimation().DrawTo(cs.world, &ebiten.DrawImageOptions{
		GeoM: backgroundAnimationGeometry,
	})

	cs.ui.Draw(cs.interfaceWorld)

	cs.world.DrawImage(cs.interfaceWorld, &ebiten.DrawImageOptions{
		ColorM: options.GetTransparentDrawOptions(
			cs.transparentTransitionEffect.GetValue()).ColorM})

	screen.DrawImage(cs.world, &ebiten.DrawImageOptions{})
}

func newCollectionsScreen() screen.Screen {
	transparentTransitionEffect := transparent.NewTransparentTransitionEffect(true, 255, 0, 5, time.Microsecond*10)

	collections.GetInstance().SetBackCallback(func() {
		transparentTransitionEffect.Reset()

		dispatcher.GetInstance().Dispatch(
			action.NewSetResetCollections(value.RESET_COLLECTIONS_FALSE_VALUE))

		dispatcher.GetInstance().Dispatch(
			action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))
	})

	collections.GetInstance().SetEntrySelectedCallback(func(path string) {
		dispatcher.GetInstance().Dispatch(
			action.NewSetLetterNameAction(path))
	})

	return &CollectionsScreen{
		ui: builder.Build(
			collections.GetInstance().GetContainer()),
		transparentTransitionEffect: transparentTransitionEffect,
		world: ebiten.NewImage(
			config.GetWorldWidth(), config.GetWorldHeight()),
		interfaceWorld: ebiten.NewImage(
			config.GetWorldWidth(), config.GetWorldHeight()),
	}
}
