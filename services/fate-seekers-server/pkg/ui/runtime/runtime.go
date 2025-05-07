package runtime

import (
	"image/color"
	"time"

	"github.com/Frabjous-Studios/asebiten"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen/entry"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/effect/transition"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/effect/transition/transparent"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/loader"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/screen/menu"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/screen/settings"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/store"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/value"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/tools/imgui"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/tools/mask"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/tools/options"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/ui/builder"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/ui/component/notification"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/ui/component/prompt"
	notificationmanager "github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/ui/manager/notification"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
)

// Runtime represents main runtime flow implementation.
type Runtime struct {
	// Represents attached notification user interface.
	notificationInterface *ebitenui.UI

	// Represents attached prompt user interface.
	promptInterface *ebitenui.UI

	// Represents transparent transition effect used for notification component.
	notificationTransparentTransitionEffect transition.TransitionEffect

	// Represents transparent transition effect used for prompt component.
	promptTransparentTransitionEffect transition.TransitionEffect

	// Represents notification interface world view.
	notificationInterfaceWorld *ebiten.Image

	// Represents prompt interface world view.
	promptInterfaceWorld *ebiten.Image

	// Represents prompt interface mask.
	promptInterfaceMask *ebiten.Image

	// Represents currently active screen.
	activeScreen screen.Screen

	// Represents global loader animation.
	loaderAnimation *asebiten.Animation
}

// Update performs logic update operations.
func (r *Runtime) Update() error {
	if store.GetApplicationExit() == value.EXIT_APPLICATION_TRUE_VALUE {
		return ebiten.Termination
	}

	if store.GetApplicationLoading() == value.LOADING_APPLICATION_TRUE_VALUE {
		r.loaderAnimation.Update()
	}

	if !notificationmanager.GetInstance().GetVisible() {
		if !r.notificationTransparentTransitionEffect.Done() {
			if !r.notificationTransparentTransitionEffect.OnEnd() {
				r.notificationTransparentTransitionEffect.Update()
			} else {
				notificationmanager.GetInstance().ToggleVisible()

				r.notificationTransparentTransitionEffect.Clean()
			}
		}
	}

	notificationmanager.GetInstance().Update()

	if !notificationmanager.GetInstance().GetTextUpdated() {
		r.notificationTransparentTransitionEffect.Reset()
	}

	if store.GetPromptText() != value.PROMPT_TEXT_EMPTY_VALUE {
		if !r.promptTransparentTransitionEffect.Done() {
			if !r.promptTransparentTransitionEffect.OnEnd() {
				r.promptTransparentTransitionEffect.Update()
			} else {
				r.promptTransparentTransitionEffect.Clean()
			}
		}
	}

	switch store.GetActiveScreen() {
	case value.ACTIVE_SCREEN_MENU_VALUE:
		r.activeScreen = menu.GetInstance()

	case value.ACTIVE_SCREEN_SETTINGS_VALUE:
		r.activeScreen = settings.GetInstance()
	}

	if store.GetPromptText() != value.PROMPT_TEXT_EMPTY_VALUE {
		if store.GetPromptUpdated() == value.PROMPT_UPDATED_FALSE_VALUE {
			prompt.GetInstance().SetText(store.GetPromptText())

			dispatcher.GetInstance().Dispatch(
				action.NewSetPromptUpdated(value.PROMPT_UPDATED_TRUE_VALUE))
		}

		r.promptInterface.Update()
	}

	err := r.activeScreen.HandleInput()
	if err != nil {
		return err
	}

	if !notificationmanager.GetInstance().IsEmpty() {
		r.notificationInterface.Update()
	}

	if config.GetDebug() {
		imgui.GetInstance().Update()
	}

	return nil
}

// Draw performs render operation.
func (r *Runtime) Draw(screen *ebiten.Image) {
	if !notificationmanager.GetInstance().IsEmpty() {
		r.notificationInterfaceWorld.Clear()
	}

	if store.GetPromptText() != value.PROMPT_TEXT_EMPTY_VALUE {
		r.promptInterfaceWorld.Clear()
	}

	r.activeScreen.HandleRender(screen)

	if store.GetApplicationLoading() == value.LOADING_APPLICATION_TRUE_VALUE {
		var loadingAnimationGeometry ebiten.GeoM

		loadingAnimationGeometry.Translate(
			float64(scaler.GetPercentageOf(config.GetWorldWidth(), 2)),
			float64(scaler.GetPercentageOf(config.GetWorldHeight(), 91)))

		r.loaderAnimation.DrawTo(screen, &ebiten.DrawImageOptions{GeoM: loadingAnimationGeometry})
	}

	if !notificationmanager.GetInstance().IsEmpty() {
		r.notificationInterface.Draw(r.notificationInterfaceWorld)
	}

	screen.DrawImage(r.notificationInterfaceWorld, &ebiten.DrawImageOptions{
		ColorM: options.GetTransparentDrawOptions(
			r.notificationTransparentTransitionEffect.GetValue()).ColorM})

	if store.GetPromptText() != value.PROMPT_TEXT_EMPTY_VALUE {
		screen.DrawImage(r.promptInterfaceMask, &ebiten.DrawImageOptions{
			ColorM: mask.GetMaskEffect(80).ColorM,
		})

		r.promptInterface.Draw(r.promptInterfaceWorld)

		screen.DrawImage(r.promptInterfaceWorld, &ebiten.DrawImageOptions{
			ColorM: options.GetTransparentDrawOptions(
				r.promptTransparentTransitionEffect.GetValue()).ColorM})
	}

	if config.GetDebug() {
		imgui.GetInstance().Draw(screen)
	}
}

// Layout manages virtual world size.
func (r *Runtime) Layout(outsideWidth, outsideHeight int) (int, int) {
	if config.GetDebug() {
		imgui.GetInstance().Layout(outsideWidth, outsideHeight)
	}

	return config.GetWorldWidth(), config.GetWorldHeight()
}

// NewRuntime creates new instance of Runtime.
func NewRuntime() *Runtime {
	promptTransparentTransitionEffect := transparent.NewTransparentTransitionEffect(true, 255, 0, 5, time.Microsecond*10)

	promptInterfaceMask := ebiten.NewImage(
		config.GetWorldWidth(), config.GetWorldHeight())

	promptInterfaceMask.Fill(color.Black)

	prompt.GetInstance().SetSubmitCallback(func() {
		store.GetPromptSubmitCallback()()

		dispatcher.GetInstance().Dispatch(
			action.NewSetPromptSubmitCallback(value.PROMPT_SUBMIT_CALLBACK_EMPTY_VALUE))

		dispatcher.GetInstance().Dispatch(
			action.NewSetPromptCancelCallback(value.PROMPT_CANCEL_CALLBACK_EMPTY_VALUE))

		dispatcher.GetInstance().Dispatch(
			action.NewSetPromptText(value.PROMPT_TEXT_EMPTY_VALUE))

		dispatcher.GetInstance().Dispatch(
			action.NewSetPromptUpdated(value.PROMPT_UPDATED_FALSE_VALUE))
	})

	prompt.GetInstance().SetCloseCallback(func() {
		store.GetPromptCancelCallback()()

		dispatcher.GetInstance().Dispatch(
			action.NewSetPromptSubmitCallback(value.PROMPT_SUBMIT_CALLBACK_EMPTY_VALUE))

		dispatcher.GetInstance().Dispatch(
			action.NewSetPromptCancelCallback(value.PROMPT_CANCEL_CALLBACK_EMPTY_VALUE))

		dispatcher.GetInstance().Dispatch(
			action.NewSetPromptText(value.PROMPT_TEXT_EMPTY_VALUE))

		dispatcher.GetInstance().Dispatch(
			action.NewSetPromptUpdated(value.PROMPT_UPDATED_FALSE_VALUE))

		promptTransparentTransitionEffect.Reset()
	})

	return &Runtime{
		notificationInterface: builder.Build(
			notification.GetInstance().GetContainer()),
		promptInterface: builder.Build(
			prompt.GetInstance().GetContainer()),
		notificationTransparentTransitionEffect: transparent.NewTransparentTransitionEffect(true, 255, 0, 5, time.Microsecond*10),
		promptTransparentTransitionEffect:       promptTransparentTransitionEffect,
		notificationInterfaceWorld: ebiten.NewImage(
			config.GetWorldWidth(), config.GetWorldHeight()),
		promptInterfaceWorld: ebiten.NewImage(
			config.GetWorldWidth(), config.GetWorldHeight()),
		promptInterfaceMask: promptInterfaceMask,

		// Guarantees non blocking rendering, if state management fails.
		activeScreen:    entry.GetInstance(),
		loaderAnimation: loader.GetInstance().GetAnimation(loader.LoaderAnimation, true),
	}
}
