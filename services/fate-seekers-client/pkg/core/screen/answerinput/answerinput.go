package answerinput

import (
	"strconv"
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition/transparent"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/builder"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/answerinput"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/common"
	answerinputmanager "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/answerinput"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/notification"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/storage/shared"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// GetInstance retrieves instance of the answer input screen, performing initilization if needed.
	GetInstance = sync.OnceValue[screen.Screen](newAnswerInputScreen)
)

// AnswerInputScreen represents answer input screen implementation.
type AnswerInputScreen struct {
	// Represents attached user interface.
	ui *ebitenui.UI

	// Represents transparent transition effect.
	transparentTransitionEffect transition.TransitionEffect

	// Represents global world view.
	world *ebiten.Image

	// Represents interface world view.
	interfaceWorld *ebiten.Image
}

func (ais *AnswerInputScreen) HandleInput() error {
	if !shared.GetInstance().GetBlinkingScreenAnimation().OnEnd() {
		shared.GetInstance().GetBlinkingScreenAnimation().Update()
	} else {
		if !ais.transparentTransitionEffect.Done() {
			if !ais.transparentTransitionEffect.OnEnd() {
				ais.transparentTransitionEffect.Update()
			} else {
				ais.transparentTransitionEffect.Clean()
			}
		}
	}

	if store.GetAnswerInputQuestionUpdated() == value.ANSWER_INPUT_QUESTION_UPDATED_FALSE_VALUE {
		answerinputmanager.GetInstance().UpdateQuestion()

		answerinput.GetInstance().SetText(
			answerinputmanager.GetInstance().GetGeneratedQuestion().Question)

		dispatcher.GetInstance().Dispatch(
			action.NewSetAnswerInputQuestionUpdated(value.ANSWER_INPUT_QUESTION_UPDATED_TRUE_VALUE))
	}

	ais.ui.Update()

	return nil
}

func (ais *AnswerInputScreen) HandleNetworking() {

}

func (ais *AnswerInputScreen) HandleRender(screen *ebiten.Image) {
	ais.world.Clear()

	ais.interfaceWorld.Clear()

	var blinkingScreenAnimationGeometry ebiten.GeoM

	blinkingScreenAnimationGeometry.Scale(
		scaler.GetScaleFactor(
			shared.GetInstance().GetBlinkingScreenAnimation().GetFrameWidth(),
			config.GetWorldWidth()),
		scaler.GetScaleFactor(
			shared.GetInstance().GetBlinkingScreenAnimation().GetFrameHeight(),
			config.GetWorldHeight()))

	shared.GetInstance().GetBlinkingScreenAnimation().DrawTo(ais.world, &ebiten.DrawImageOptions{
		GeoM: blinkingScreenAnimationGeometry,
	})

	ais.ui.Draw(ais.interfaceWorld)

	screen.DrawImage(ais.world, &ebiten.DrawImageOptions{})

	screen.DrawImage(ais.interfaceWorld, &ebiten.DrawImageOptions{
		ColorM: ais.transparentTransitionEffect.GetOptions().ColorM})
}

func (ais *AnswerInputScreen) Clean() {

}

// newAnswerInputScreen initializes AnswerInputScreen.
func newAnswerInputScreen() screen.Screen {
	transparentTransitionEffect := transparent.NewTransparentTransitionEffect()

	answerinput.GetInstance().SetSubmitCallback(func(valueRaw string) {
		result, err := strconv.Atoi(valueRaw)
		if err != nil {
			notification.GetInstance().Push(
				"Answer is incorrect!", time.Second*2, common.NotificationErrorTextColor)

			return
		}

		if answerinputmanager.GetInstance().GetGeneratedQuestion().Answer == result {
			notification.GetInstance().Push(
				"Answer is correct!", time.Second*2, common.NotificationInfoTextColor)

			dispatcher.GetInstance().Dispatch(
				action.NewSetAnswerInputQuestionUpdated(value.ANSWER_INPUT_QUESTION_UPDATED_FALSE_VALUE))
		} else {
			notification.GetInstance().Push(
				"Answer is incorrect!", time.Second*2, common.NotificationErrorTextColor)
		}
	})

	answerinput.GetInstance().SetCloseCallback(func() {
		shared.GetInstance().GetBlinkingScreenAnimation().Reset()

		dispatcher.GetInstance().Dispatch(
			action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_SESSION_VALUE))

		dispatcher.GetInstance().Dispatch(
			action.NewSetAnswerInputQuestionUpdated(value.ANSWER_INPUT_QUESTION_UPDATED_FALSE_VALUE))

		transparentTransitionEffect.Reset()
	})

	return &AnswerInputScreen{
		ui: builder.Build(
			answerinput.GetInstance().GetContainer()),
		transparentTransitionEffect: transparentTransitionEffect,
		world:                       ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
		interfaceWorld:              ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
	}
}
