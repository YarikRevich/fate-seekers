package resume

import (
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition/transparent"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/handler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/options"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/builder"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/resume"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/notification"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/translation"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/storage/shared"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// GetInstance retrieves instance of the resume screen, performing initilization if needed.
	GetInstance = sync.OnceValue[screen.Screen](newResumeScreen)
)

// ResumeScreen represents resume screen implementation.
type ResumeScreen struct {
	// Represents attached user interface.
	ui *ebitenui.UI

	// Represents transparent transition effect.
	transparentTransitionEffect transition.TransitionEffect

	// Represents global world view.
	world *ebiten.Image

	// Represents interface world view.
	interfaceWorld *ebiten.Image
}

func (rs *ResumeScreen) HandleInput() error {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		rs.transparentTransitionEffect.Reset()

		dispatcher.GetInstance().Dispatch(
			action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_SESSION_VALUE))
	}

	if !rs.transparentTransitionEffect.Done() {
		if !rs.transparentTransitionEffect.OnEnd() {
			rs.transparentTransitionEffect.Update()
		} else {
			rs.transparentTransitionEffect.Clean()
		}
	}

	shared.GetInstance().GetBackgroundAnimation().Update()

	rs.ui.Update()

	return nil
}

func (rs *ResumeScreen) HandleRender(screen *ebiten.Image) {
	rs.world.Clear()

	rs.interfaceWorld.Clear()

	var backgroundAnimationGeometry ebiten.GeoM

	backgroundAnimationGeometry.Scale(
		scaler.GetScaleFactor(config.GetMinStaticWidth(), config.GetWorldWidth()),
		scaler.GetScaleFactor(config.GetMinStaticHeight(), config.GetWorldHeight()))

	shared.GetInstance().GetBackgroundAnimation().DrawTo(rs.world, &ebiten.DrawImageOptions{
		GeoM: backgroundAnimationGeometry,
	})

	rs.ui.Draw(rs.interfaceWorld)

	rs.world.DrawImage(rs.interfaceWorld, &ebiten.DrawImageOptions{
		ColorM: options.GetTransparentDrawOptions(
			rs.transparentTransitionEffect.GetValue()).ColorM,
	})

	screen.DrawImage(rs.world, &ebiten.DrawImageOptions{})
}

func newResumeScreen() screen.Screen {
	transparentTransitionEffect := transparent.NewTransparentTransitionEffect(true, 255, 0, 5, time.Microsecond*10)

	return &ResumeScreen{
		ui: builder.Build(
			resume.NewResumeComponent(
				func() {
					transparentTransitionEffect.Reset()

					dispatcher.GetInstance().Dispatch(
						action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_SESSION_VALUE))
				},
				func() {
					transparentTransitionEffect.Reset()

					dispatcher.GetInstance().Dispatch(
						action.NewSetPreviousScreenAction(value.PREVIOUS_SCREEN_RESUME_VALUE))

					dispatcher.GetInstance().Dispatch(
						action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_SETTINGS_VALUE))
				},
				func() {
					transparentTransitionEffect.Reset()

					handler.PerformLeaveLobby(store.GetSelectedSessionMetadata().ID, func(err error) {
						if err != nil {
							notification.GetInstance().Push(
								common.ComposeMessage(
									translation.GetInstance().GetTranslation("client.networking.leave-lobby-failure"),
									err.Error()),
								time.Second*3,
								common.NotificationErrorTextColor)
						}

						dispatcher.
							GetInstance().
							Dispatch(
								action.NewSetStateResetApplicationAction(
									value.STATE_RESET_APPLICATION_FALSE_VALUE))

						dispatcher.GetInstance().Dispatch(
							action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))
					})
				})),
		transparentTransitionEffect: transparentTransitionEffect,
		world:                       ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
		interfaceWorld:              ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
	}
}
