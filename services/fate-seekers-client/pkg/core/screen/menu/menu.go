package menu

import (
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition/transparent"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/connector"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/handler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/options"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/builder"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/menu"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/notification"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/translation"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/validator/encryptionkey"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/storage/shared"
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
								translation.GetInstance().GetTranslation("shared.prompt.networking.encryption-key")))

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
						dispatcher.GetInstance().Dispatch(action.NewIncrementLoadingApplicationAction())

						dispatcher.GetInstance().Dispatch(
							action.NewSetEntryHandshakeStartedNetworkingAction(value.ENTRY_HANDSHAKE_STARTED_NETWORKING_TRUE_VALUE))

						connector.GetInstance().Clean(func() {
							connector.GetInstance().Connect(
								func(err1 error) {
									if err1 != nil {
										notification.GetInstance().Push(
											common.ComposeMessage(
												translation.GetInstance().GetTranslation("client.networking.connection-failure"),
												err1.Error()),
											time.Second*2,
											common.NotificationErrorTextColor)

										return
									}

									handler.PerformPingConnection(func(err2 error) {
										dispatcher.GetInstance().Dispatch(
											action.NewDecrementLoadingApplicationAction())

										if err2 != nil {
											notification.GetInstance().Push(
												common.ComposeMessage(
													translation.GetInstance().GetTranslation("client.networking.ping-connection-failure"),
													err2.Error()),
												time.Second*2,
												common.NotificationErrorTextColor)

											dispatcher.GetInstance().Dispatch(
												action.NewSetEntryHandshakeStartedNetworkingAction(value.ENTRY_HANDSHAKE_STARTED_NETWORKING_FALSE_VALUE))

											return
										}

										transparentTransitionEffect.Reset()

										dispatcher.GetInstance().Dispatch(
											action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_SELECTOR_VALUE))
									})
								},
								func(err error) {
									notification.GetInstance().Push(
										err.Error(),
										time.Second*2,
										common.NotificationErrorTextColor)

									dispatcher.GetInstance().Dispatch(
										action.NewIncrementLoadingApplicationAction())

									connector.GetInstance().Close(func(err error) {
										if err != nil {
											notification.GetInstance().Push(
												translation.GetInstance().GetTranslation("client.networking.close-failure"),
												time.Second*2,
												common.NotificationErrorTextColor)
										}

										dispatcher.GetInstance().Dispatch(
											action.NewDecrementLoadingApplicationAction())

										dispatcher.GetInstance().Dispatch(
											action.NewSetEntryHandshakeStartedNetworkingAction(value.ENTRY_HANDSHAKE_STARTED_NETWORKING_FALSE_VALUE))

										dispatcher.GetInstance().Dispatch(
											action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))
									})
								})
						})
					} else if store.GetPingConnectionStartedNetworking() == value.PING_CONNECTION_STARTED_NETWORKING_FALSE_VALUE {
						dispatcher.GetInstance().Dispatch(action.NewIncrementLoadingApplicationAction())

						dispatcher.GetInstance().Dispatch(
							action.NewSetPingConnectionStartedNetworkingAction(value.PING_CONNECTION_STARTED_NETWORKING_TRUE_VALUE))

						handler.PerformPingConnection(func(err error) {
							dispatcher.GetInstance().Dispatch(
								action.NewDecrementLoadingApplicationAction())

							dispatcher.GetInstance().Dispatch(
								action.NewSetPingConnectionStartedNetworkingAction(value.PING_CONNECTION_STARTED_NETWORKING_FALSE_VALUE))

							if err != nil {
								notification.GetInstance().Push(
									common.ComposeMessage(
										translation.GetInstance().GetTranslation("client.networking.ping-connection-failure"),
										err.Error()),
									time.Second*2,
									common.NotificationErrorTextColor)

								dispatcher.GetInstance().Dispatch(
									action.NewSetEntryHandshakeStartedNetworkingAction(value.ENTRY_HANDSHAKE_STARTED_NETWORKING_FALSE_VALUE))

								return
							}

							dispatcher.GetInstance().Dispatch(
								action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_SELECTOR_VALUE))
						})
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
