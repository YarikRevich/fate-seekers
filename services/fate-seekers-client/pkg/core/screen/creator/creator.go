package creator

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
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/creator"
	creatormanager "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/creator"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/storage/shared"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// GetInstance retrieves instance of the creator screen, performing initilization if needed.
	GetInstance = sync.OnceValue[screen.Screen](newCreatorScreen)
)

// CreatorScreen represents creator screen implementation.
type CreatorScreen struct {
	// Represents attached user interface.
	ui *ebitenui.UI

	// Represents transparent transition effect.
	transparentTransitionEffect transition.TransitionEffect

	// Represents global world view.
	world *ebiten.Image

	// Represents interface world view.
	interfaceWorld *ebiten.Image
}

func (cs *CreatorScreen) HandleInput() error {
	// TODO: check if form has been updated

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

func (cs *CreatorScreen) HandleRender(screen *ebiten.Image) {
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

// newCreatorScreen initializes CreatorScreen.
func newCreatorScreen() screen.Screen {
	transparentTransitionEffect := transparent.NewTransparentTransitionEffect(true, 255, 0, 5, time.Microsecond*10)

	creator.GetInstance().SetSubmitCallback(func(name, seed string) {
		if creatormanager.ProcessChanges(name, seed) {
			// if store.GetSessionRetrievalStartedNetworking() == value.SESSION_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE {
			// 	dispatcher.GetInstance().Dispatch(
			// 		action.NewSetSessionRetrievalStartedNetworkingAction(value.SESSION_RETRIEVAL_STARTED_NETWORKING_TRUE_VALUE))

			// 	handler.PerformGetSessions(func(response *api.GetSessionsResponse, err error) {
			// 		fmt.Println(response.GetSessions(), err)
			// 	})
			// }

			transparentTransitionEffect.Reset()

			dispatcher.GetInstance().Dispatch(
				action.NewSetSessionRetrievalStartedNetworkingAction(value.SESSION_RETRIEVAL_STARTED_NETWORKING_TRUE_VALUE))

			dispatcher.GetInstance().Dispatch(
				action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_LOBBY_VALUE))
		}

		transparentTransitionEffect.Reset()

		dispatcher.GetInstance().Dispatch(
			action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_SELECTOR_VALUE))
	})

	creator.GetInstance().SetBackCallback(func() {
		transparentTransitionEffect.Reset()

		dispatcher.GetInstance().Dispatch(
			action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_SELECTOR_VALUE))
	})

	return &CreatorScreen{
		ui:                          builder.Build(creator.GetInstance().GetContainer()),
		transparentTransitionEffect: transparentTransitionEffect,
		world:                       ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
		interfaceWorld: ebiten.NewImage(
			config.GetWorldWidth(), config.GetWorldHeight()),
	}
}
