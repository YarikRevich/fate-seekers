package settings

import (
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
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
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/ui/component/settings"
	settingsmanager "github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/ui/manager/settings"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/ui/manager/translation"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// GetInstance retrieves instance of the settings screen, performing initilization if needed.
	GetInstance = sync.OnceValue[screen.Screen](newSettingsScreen)
)

// SettingsScreen represents settings screen implementation.
type SettingsScreen struct {
	// Represents attached user interface.
	ui *ebitenui.UI

	// Represents transparent transition effect.
	transparentTransitionEffect transition.TransitionEffect

	// Represents global world view.
	world *ebiten.Image

	// Represents interface world view.
	interfaceWorld *ebiten.Image
}

func (ss *SettingsScreen) HandleInput() error {
	if !ss.transparentTransitionEffect.Done() {
		if !ss.transparentTransitionEffect.OnEnd() {
			ss.transparentTransitionEffect.Update()
		} else {
			ss.transparentTransitionEffect.Clean()
		}
	}

	shared.GetInstance().GetBackgroundAnimation().Update()

	ss.ui.Update()

	return nil
}

func (ss *SettingsScreen) HandleRender(screen *ebiten.Image) {
	ss.world.Clear()

	ss.interfaceWorld.Clear()

	var backgroundAnimationGeometry ebiten.GeoM

	backgroundAnimationGeometry.Scale(
		scaler.GetScaleFactor(config.GetMinStaticWidth(), config.GetWorldWidth()),
		scaler.GetScaleFactor(config.GetMinStaticHeight(), config.GetWorldHeight()))

	shared.GetInstance().GetBackgroundAnimation().DrawTo(ss.world, &ebiten.DrawImageOptions{
		GeoM: backgroundAnimationGeometry,
	})

	ss.ui.Draw(ss.interfaceWorld)

	ss.world.DrawImage(ss.interfaceWorld, &ebiten.DrawImageOptions{
		ColorM: options.GetTransparentDrawOptions(
			ss.transparentTransitionEffect.GetValue()).ColorM})

	screen.DrawImage(ss.world, &ebiten.DrawImageOptions{})
}

func newSettingsScreen() screen.Screen {
	transparentTransitionEffect := transparent.NewTransparentTransitionEffect(true, 255, 0, 5, time.Microsecond*10)

	return &SettingsScreen{
		ui: builder.Build(
			settings.NewSettingsComponent(
				func(soundFX int, networkingServerPort, networkingEncryptionKey, language string) {
					if settingsmanager.ProcessChanges(soundFX, networkingServerPort, networkingEncryptionKey, language) {
						dispatcher.GetInstance().Dispatch(
							action.NewSetActiveScreenAction(store.GetPreviousScreen()))

						dispatcher.GetInstance().Dispatch(
							action.NewSetPreviousScreenAction(value.PREVIOUS_SCREEN_EMPTY_VALUE))
					}
				},
				func(soundFX int, networkingServerPort, networkingEncryptionKey, language string) {
					if settingsmanager.AnyProvidedChanges(soundFX, networkingServerPort, networkingEncryptionKey, language) {
						dispatcher.GetInstance().Dispatch(
							action.NewSetPromptText(
								translation.GetInstance().GetTranslation("shared.prompt.settings")))

						dispatcher.GetInstance().Dispatch(
							action.NewSetPromptSubmitCallback(func() {
								settingsmanager.ProcessChanges(soundFX, networkingServerPort, networkingEncryptionKey, language)

								transparentTransitionEffect.Reset()

								dispatcher.GetInstance().Dispatch(
									action.NewSetActiveScreenAction(store.GetPreviousScreen()))

								dispatcher.GetInstance().Dispatch(
									action.NewSetPreviousScreenAction(value.PREVIOUS_SCREEN_EMPTY_VALUE))
							}))

						dispatcher.GetInstance().Dispatch(
							action.NewSetPromptCancelCallback(func() {
								transparentTransitionEffect.Reset()

								dispatcher.GetInstance().Dispatch(
									action.NewSetActiveScreenAction(store.GetPreviousScreen()))

								dispatcher.GetInstance().Dispatch(
									action.NewSetPreviousScreenAction(value.PREVIOUS_SCREEN_EMPTY_VALUE))
							}))
					} else {
						transparentTransitionEffect.Reset()

						dispatcher.GetInstance().Dispatch(
							action.NewSetActiveScreenAction(store.GetPreviousScreen()))

						dispatcher.GetInstance().Dispatch(
							action.NewSetPreviousScreenAction(value.PREVIOUS_SCREEN_EMPTY_VALUE))
					}
				})),
		transparentTransitionEffect: transparentTransitionEffect,
		world: ebiten.NewImage(
			config.GetWorldWidth(), config.GetWorldHeight()),
		interfaceWorld: ebiten.NewImage(
			config.GetWorldWidth(), config.GetWorldHeight()),
	}
}
