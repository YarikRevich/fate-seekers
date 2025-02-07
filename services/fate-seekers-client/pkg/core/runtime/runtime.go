package runtime

import (
	"image/color"

	"github.com/Frabjous-Studios/asebiten"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition/transparent"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen/answerinput"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen/entry"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen/intro"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen/menu"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen/session"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen/settings"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen/travel"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/imgui"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/mask"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/builder"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/letter"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/letterimage"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/notification"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/subtitles"
	notificationmanager "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/notification"
	subtitlesmanager "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/subtitles"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/application"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
)

// Runtime represents main runtime flow implementation.
type Runtime struct {
	// Represents attached subtitles user interface.
	subtitlesInterface *ebitenui.UI

	// Represents attached notification user interface.
	notificationInterface *ebitenui.UI

	// Represents attached letter image user interface.
	letterImageInterface *ebitenui.UI

	// Represents attached letter user interface.
	letterInterface *ebitenui.UI

	// Represents transparent transition effect used for notification component.
	notificationTransparentTransitionEffect transition.TransitionEffect

	// Represents transparent transition effect used for letter image component.
	letterImageTransparentTransitionEffect transition.TransitionEffect

	// Represents transparent transition effect used for letter component.
	letterTransparentTransitionEffect transition.TransitionEffect

	// Represents notification interface world view.
	notificationInterfaceWorld *ebiten.Image

	// Represents letter image interface world view.
	letterImageInterfaceWorld *ebiten.Image

	// Represents letter image interface mask.
	letterImageInterfaceMask *ebiten.Image

	// Represents letter interface world view.
	letterInterfaceWorld *ebiten.Image

	// Represents letter interface mask.
	letterInterfaceMask *ebiten.Image

	// Represents currently active screen.
	activeScreen screen.Screen

	// Represents global loader animation.
	loaderAnimation *asebiten.Animation
}

// Update performs logic update operations.
func (r *Runtime) Update() error {
	if store.GetExitApplication() == value.EXIT_APPLICATION_TRUE_VALUE {
		return ebiten.Termination
	}

	if store.GetLoadingApplication() == value.LOADING_APPLICATION_TRUE_VALUE {
		r.loaderAnimation.Update()
	}

	subtitlesmanager.GetInstance().Update()

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

	if store.GetLetterName() != value.LETTER_NAME_EMPTY_VALUE {
		if !r.letterTransparentTransitionEffect.Done() {
			if !r.letterTransparentTransitionEffect.OnEnd() {
				r.letterTransparentTransitionEffect.Update()
			} else {
				r.letterTransparentTransitionEffect.Clean()
			}
		}
	}

	if store.GetLetterImage() != value.LETTER_IMAGE_EMPTY_VALUE {
		if !r.letterImageTransparentTransitionEffect.Done() {
			if !r.letterImageTransparentTransitionEffect.OnEnd() {
				r.letterImageTransparentTransitionEffect.Update()
			} else {
				r.letterImageTransparentTransitionEffect.Clean()
			}
		}
	}

	switch store.GetActiveScreen() {
	case value.ACTIVE_SCREEN_INTRO_VALUE:
		r.activeScreen = intro.GetInstance()

	case value.ACTIVE_SCREEN_ENTRY_VALUE:
		r.activeScreen = entry.GetInstance()

	case value.ACTIVE_SCREEN_MENU_VALUE:
		r.activeScreen = menu.GetInstance()

	case value.ACTIVE_SCREEN_SETTINGS_VALUE:
		r.activeScreen = settings.GetInstance()

	case value.ACTIVE_SCREEN_SESSION_VALUE:
		r.activeScreen = session.GetInstance()

	case value.ACTIVE_SCREEN_TRAVEL_VALUE:
		r.activeScreen = travel.GetInstance()

	case value.ACTIVE_SCREEN_ANSWER_INPUT_VALUE:
		r.activeScreen = answerinput.GetInstance()
	}

	if store.GetLetterImage() != value.LETTER_IMAGE_EMPTY_VALUE {
		r.letterImageInterface.Update()
	}

	if store.GetLetterName() != value.LETTER_NAME_EMPTY_VALUE {
		if store.GetLetterUpdated() == value.LETTER_UPDATED_FALSE_VALUE {
			loadedLetter := loader.GetInstance().GetLetter(store.GetLetterName())

			letter.GetInstance().SetText(loadedLetter.Text)

			letter.GetInstance().SetAttachment(loadedLetter.Attachment)

			dispatcher.GetInstance().Dispatch(
				action.NewSetLetterUpdatedAction(value.LETTER_UPDATED_TRUE_VALUE))
		}

		r.letterInterface.Update()
	}

	err := r.activeScreen.HandleInput()
	if err != nil {
		return err
	}

	if !subtitlesmanager.GetInstance().IsEmpty() {
		r.subtitlesInterface.Update()
	}

	if !notificationmanager.GetInstance().IsEmpty() {
		r.notificationInterface.Update()
	}

	r.activeScreen.HandleNetworking()

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

	if store.GetLetterImage() != value.LETTER_IMAGE_EMPTY_VALUE {
		r.letterImageInterfaceWorld.Clear()
	}

	if store.GetLetterName() != value.LETTER_NAME_EMPTY_VALUE {
		r.letterInterfaceWorld.Clear()
	}

	r.activeScreen.HandleRender(screen)

	if store.GetInstance().GetState(application.LOADING_APPLICATION_STATE) ==
		value.LOADING_APPLICATION_TRUE_VALUE {
		var loadingAnimationGeometry ebiten.GeoM

		loadingAnimationGeometry.Translate(
			float64(scaler.GetPercentageOf(config.GetWorldWidth(), 2)),
			float64(scaler.GetPercentageOf(config.GetWorldHeight(), 91)))

		r.loaderAnimation.DrawTo(screen, &ebiten.DrawImageOptions{GeoM: loadingAnimationGeometry})
	}

	if !subtitlesmanager.GetInstance().IsEmpty() {
		r.subtitlesInterface.Draw(screen)
	}

	if !notificationmanager.GetInstance().IsEmpty() {
		r.notificationInterface.Draw(r.notificationInterfaceWorld)
	}

	screen.DrawImage(r.notificationInterfaceWorld, &ebiten.DrawImageOptions{
		ColorM: r.notificationTransparentTransitionEffect.GetOptions().ColorM})

	if store.GetLetterName() != value.LETTER_NAME_EMPTY_VALUE {
		if store.GetLetterImage() == value.LETTER_IMAGE_EMPTY_VALUE {
			screen.DrawImage(r.letterInterfaceMask, &ebiten.DrawImageOptions{
				ColorM: mask.GetMaskEffect(80).ColorM,
			})
		}

		r.letterInterface.Draw(r.letterInterfaceWorld)

		screen.DrawImage(r.letterInterfaceWorld, &ebiten.DrawImageOptions{
			ColorM: r.letterTransparentTransitionEffect.GetOptions().ColorM})
	}

	if store.GetLetterImage() != value.LETTER_IMAGE_EMPTY_VALUE {
		screen.DrawImage(r.letterInterfaceMask, &ebiten.DrawImageOptions{
			ColorM: mask.GetMaskEffect(80).ColorM,
		})

		letterImage := loader.GetInstance().GetStatic(store.GetLetterImage())

		r.letterImageInterfaceWorld.DrawImage(
			letterImage,
			&ebiten.DrawImageOptions{
				GeoM: scaler.GetCenteredGeometry(
					50,
					70,
					letterImage.Bounds().Dx(),
					letterImage.Bounds().Dy(),
				)})

		r.letterImageInterface.Draw(r.letterImageInterfaceWorld)

		screen.DrawImage(r.letterImageInterfaceWorld, &ebiten.DrawImageOptions{
			ColorM: r.letterImageTransparentTransitionEffect.GetOptions().ColorM})
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
	letterImageTransparentTransitionEffect := transparent.NewTransparentTransitionEffect()

	letterTransparentTransitionEffect := transparent.NewTransparentTransitionEffect()

	letterImageInterfaceMask := ebiten.NewImage(
		config.GetWorldWidth(), config.GetWorldHeight())

	letterImageInterfaceMask.Fill(color.Black)

	letterInterfaceMask := ebiten.NewImage(
		config.GetWorldWidth(), config.GetWorldHeight())

	letterInterfaceMask.Fill(color.Black)

	letter.GetInstance().SetAttachmentCallback(func(value string) {
		dispatcher.GetInstance().Dispatch(
			action.NewSetLetterImageAction(value))
	})

	letter.GetInstance().SetCloseCallback(func() {
		dispatcher.GetInstance().Dispatch(
			action.NewSetLetterUpdatedAction(value.LETTER_UPDATED_FALSE_VALUE))

		dispatcher.GetInstance().Dispatch(
			action.NewSetLetterNameAction(value.LETTER_NAME_EMPTY_VALUE))
	})

	return &Runtime{
		subtitlesInterface: builder.Build(
			subtitles.GetInstance().GetContainer()),
		notificationInterface: builder.Build(
			notification.GetInstance().GetContainer()),
		letterImageInterface: builder.Build(
			letterimage.NewLetterImageComponent(func() {
				dispatcher.GetInstance().Dispatch(
					action.NewSetLetterImageAction(value.LETTER_IMAGE_EMPTY_VALUE))

				letterImageTransparentTransitionEffect.Reset()
			}),
		),
		letterInterface: builder.Build(
			letter.GetInstance().GetContainer()),
		notificationTransparentTransitionEffect: transparent.NewTransparentTransitionEffect(),
		letterImageTransparentTransitionEffect:  letterImageTransparentTransitionEffect,
		letterTransparentTransitionEffect:       letterTransparentTransitionEffect,
		notificationInterfaceWorld: ebiten.NewImage(
			config.GetWorldWidth(), config.GetWorldHeight()),
		letterImageInterfaceWorld: ebiten.NewImage(
			config.GetWorldWidth(), config.GetWorldHeight()),
		letterImageInterfaceMask: letterImageInterfaceMask,
		letterInterfaceWorld: ebiten.NewImage(
			config.GetWorldWidth(), config.GetWorldHeight()),
		letterInterfaceMask: letterInterfaceMask,

		// Guarantees non blocking rendering, if state management fails.
		activeScreen:    entry.GetInstance(),
		loaderAnimation: loader.GetInstance().GetAnimation(loader.LoaderAnimation, true),
	}
}
