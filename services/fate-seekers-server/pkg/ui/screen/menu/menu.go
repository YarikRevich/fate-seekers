package menu

import (
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/connector"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/storage/shared"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/effect/transition"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/effect/transition/transparent"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/store"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/value"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/tools/options"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/ui/builder"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/ui/component/menu"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/ui/manager/translation"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/ui/validator/encryptionkey"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// GetInstance retrieves instance of the menu screen, performing initilization if needed.
	GetInstance = sync.OnceValue[screen.Screen](newMenuScreen)
)

// MenuScreen represents menu screen implementation.
type MenuScreen struct {
	// Represents attached user interface.
	ui *ebitenui.UI

	// Represents transparent transition effect.
	transparentTransitionEffect transition.TransitionEffect

	// Represents global world view.
	world *ebiten.Image

	// Represents interface world view.
	interfaceWorld *ebiten.Image
}

func (ms *MenuScreen) HandleInput() error {
	if !ms.transparentTransitionEffect.Done() {
		if !ms.transparentTransitionEffect.OnEnd() {
			ms.transparentTransitionEffect.Update()
		} else {
			ms.transparentTransitionEffect.Clean()
		}
	}

	shared.GetInstance().GetBackgroundAnimation().Update()

	ms.ui.Update()

	return nil
}

func (ms *MenuScreen) HandleRender(screen *ebiten.Image) {
	ms.world.Clear()

	ms.interfaceWorld.Clear()

	var backgroundAnimationGeometry ebiten.GeoM

	backgroundAnimationGeometry.Scale(
		scaler.GetScaleFactor(config.GetMinStaticWidth(), config.GetWorldWidth()),
		scaler.GetScaleFactor(config.GetMinStaticHeight(), config.GetWorldHeight()))

	shared.GetInstance().GetBackgroundAnimation().DrawTo(ms.world, &ebiten.DrawImageOptions{
		GeoM: backgroundAnimationGeometry,
	})

	ms.ui.Draw(ms.interfaceWorld)

	ms.world.DrawImage(ms.interfaceWorld, &ebiten.DrawImageOptions{
		ColorM: options.GetTransparentDrawOptions(
			ms.transparentTransitionEffect.GetValue()).ColorM,
	})

	screen.DrawImage(ms.world, &ebiten.DrawImageOptions{})
}

func newMenuScreen() screen.Screen {
	transparentTransitionEffect := transparent.NewTransparentTransitionEffect(true, 255, 0, 5, time.Microsecond*10)

	return &MenuScreen{
		ui: builder.Build(
			menu.NewMenuComponent(
				func() {
					if !encryptionkey.Validate(config.GetSettingsNetworkingEncryptionKey()) {
						dispatcher.GetInstance().Dispatch(
							action.NewSetPromptText(
								translation.GetInstance().GetTranslation("prompt.networking.encryption-key")))

						dispatcher.GetInstance().Dispatch(
							action.NewSetPromptSubmitCallback(func() {
								transparentTransitionEffect.Reset()

								dispatcher.GetInstance().Dispatch(
									action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_SETTINGS_VALUE))

								dispatcher.GetInstance().Dispatch(
									action.NewSetPreviousScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))
							}))

						dispatcher.GetInstance().Dispatch(
							action.NewSetPromptCancelCallback(func() {
								dispatcher.GetInstance().Dispatch(
									action.NewSetActiveScreenAction(store.GetPreviousScreen()))

								dispatcher.GetInstance().Dispatch(
									action.NewSetPreviousScreenAction(value.PREVIOUS_SCREEN_EMPTY_VALUE))
							}))

						return
					}

					if store.GetEntryHandshakeStartedNetworking() == value.ENTRY_HANDSHAKE_STARTED_NETWORKING_FALSE_VALUE {
						connector.GetInstance().Connect(func(err error) {
							// TODO: check if error is related to encryption key.

							dispatcher.GetInstance().Dispatch(
								action.NewSetLoadingApplicationAction(value.LOADING_APPLICATION_FALSE_VALUE))

							if err == nil {
								transparentTransitionEffect.Reset()

								dispatcher.GetInstance().Dispatch(
									action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_ANSWER_INPUT_VALUE))
							}
						})

						dispatcher.GetInstance().Dispatch(
							action.NewSetLoadingApplicationAction(value.LOADING_APPLICATION_TRUE_VALUE))

						dispatcher.GetInstance().Dispatch(
							action.NewSetEntryHandshakeStartedNetworkingAction(value.ENTRY_HANDSHAKE_STARTED_NETWORKING_TRUE_VALUE))
					}
				},
				func() {

				},
				func() {

				},
				func() {
					transparentTransitionEffect.Reset()

					dispatcher.GetInstance().Dispatch(
						action.NewSetPreviousScreenAction(value.PREVIOUS_SCREEN_MENU_VALUE))

					dispatcher.GetInstance().Dispatch(
						action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_SETTINGS_VALUE))
				},
				func() {
					dispatcher.GetInstance().Dispatch(
						action.NewSetExitApplicationAction(value.EXIT_APPLICATION_TRUE_VALUE))
				})),
		transparentTransitionEffect: transparentTransitionEffect,
		world:                       ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
		interfaceWorld:              ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
	}
}
