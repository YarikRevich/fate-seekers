package settings

import (
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition/transparent"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/builder"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/settings"
	settingsmanager "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/settings"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/translation"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/storage/shared"
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

func (ss *SettingsScreen) HandleNetworking() {

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
		ColorM: ss.transparentTransitionEffect.GetOptions().ColorM})

	screen.DrawImage(ss.world, &ebiten.DrawImageOptions{})
}

func (ss *SettingsScreen) Clean() {

}

func newSettingsScreen() screen.Screen {
	transparentTransitionEffect := transparent.NewTransparentTransitionEffect()

	return &SettingsScreen{
		ui: builder.Build(
			settings.NewSettingsComponent(
				func(soundMusic, soundFX int, networkingHost, language string) {
					settingsmanager.ProcessChanges(soundMusic, soundFX, networkingHost, language)

					dispatcher.GetInstance().Dispatch(
						action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))
				},
				func(soundMusic, soundFX int, networkingHost, language string) {
					if settingsmanager.AnyProvidedChanges(soundMusic, soundFX, networkingHost, language) {
						dispatcher.GetInstance().Dispatch(
							action.NewSetPromptText(
								translation.GetInstance().GetTranslation("prompt.settings")))

						dispatcher.GetInstance().Dispatch(
							action.NewSetPromptSubmitCallback(func() {
								settingsmanager.ProcessChanges(soundMusic, soundFX, networkingHost, language)

								transparentTransitionEffect.Reset()

								dispatcher.GetInstance().Dispatch(
									action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))
							}))

						dispatcher.GetInstance().Dispatch(
							action.NewSetPromptCancelCallback(func() {
								transparentTransitionEffect.Reset()

								dispatcher.GetInstance().Dispatch(
									action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))
							}))
					} else {
						transparentTransitionEffect.Reset()

						dispatcher.GetInstance().Dispatch(
							action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))
					}
				})),
		transparentTransitionEffect: transparentTransitionEffect,
		world: ebiten.NewImage(
			config.GetWorldWidth(), config.GetWorldHeight()),
		interfaceWorld: ebiten.NewImage(
			config.GetWorldWidth(), config.GetWorldHeight()),
	}
}
