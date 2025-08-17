package lobby

import (
	"fmt"
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/effect/transition/transparent"
	metadatav1 "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/stream"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/options"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/scaler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/builder"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/storage/shared"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// GetInstance retrieves instance of the lobby screen, performing initilization if needed.
	GetInstance = sync.OnceValue[screen.Screen](newLobbyScreen)
)

// LobbyScreen represents lobby screen implementation.
type LobbyScreen struct {
	// Represents attached user interface.
	ui *ebitenui.UI

	// Represents transparent transition effect.
	transparentTransitionEffect transition.TransitionEffect

	// Represents global world view.
	world *ebiten.Image
}

func (ls *LobbyScreen) HandleInput() error {
	if store.GetSessionMetadataRetrievalStartedNetworking() == value.SESSION_METADATA_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE {
		dispatcher.GetInstance().Dispatch(
			action.NewSetSessionMetadataRetrievalStartedNetworkingAction(
				value.SESSION_METADATA_RETRIEVAL_STARTED_NETWORKING_TRUE_VALUE))

		var sessionID int64

		for _, session := range store.GetRetrievedSessionsMetadata() {
			if session.Name == store.GetSelectedSessionMetadata() {
				sessionID = session.SessionID

				break
			}
		}

		stream.GetGetSessionMetadataSubmitter().Clean(func() {
			stream.GetGetSessionMetadataSubmitter().Submit(
				sessionID, func(response *metadatav1.GetSessionMetadataResponse, err error) bool {
					fmt.Println(response.GetStarted(), err)

					return false
				})
		})
	}

	if !ls.transparentTransitionEffect.Done() {
		if !ls.transparentTransitionEffect.OnEnd() {
			ls.transparentTransitionEffect.Update()
		} else {
			ls.transparentTransitionEffect.Clean()
		}
	}

	shared.GetInstance().GetBackgroundAnimation().Update()

	ls.ui.Update()

	return nil
}

func (ls *LobbyScreen) HandleRender(screen *ebiten.Image) {
	var backgroundAnimationGeometry ebiten.GeoM

	backgroundAnimationGeometry.Scale(
		scaler.GetScaleFactor(config.GetMinStaticWidth(), config.GetWorldWidth()),
		scaler.GetScaleFactor(config.GetMinStaticHeight(), config.GetWorldHeight()))

	shared.GetInstance().GetBackgroundAnimation().DrawTo(ls.world, &ebiten.DrawImageOptions{
		GeoM: backgroundAnimationGeometry,
	})

	ls.ui.Draw(ls.world)

	screen.DrawImage(ls.world, &ebiten.DrawImageOptions{
		ColorM: options.GetTransparentDrawOptions(
			ls.transparentTransitionEffect.GetValue()).ColorM})
}

// newLobbyScreen initializes LobbyScreen.
func newLobbyScreen() screen.Screen {
	return &LobbyScreen{
		ui: builder.Build(
		// TODO: add lobby initialization
		),
		transparentTransitionEffect: transparent.NewTransparentTransitionEffect(true, 255, 0, 5, time.Microsecond*10),
		world:                       ebiten.NewImage(config.GetWorldWidth(), config.GetWorldHeight()),
	}
}
