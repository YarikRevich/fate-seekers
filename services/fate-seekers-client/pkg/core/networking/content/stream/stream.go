package stream

import (
	"context"
	"fmt"
	"sync"
	"time"

	contentv1 "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/content/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/content/handler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"google.golang.org/protobuf/proto"
)

// Describes constant values used for stream management.
const (
	updateUserMetadataPositionsFrequency = time.Millisecond * 100
)

var (
	// GetUpdateUserMetadataPositionsSubmitter retrieves instance of the update user metadata positions submitter, performing initial creation if needed.
	GetUpdateUserMetadataPositionsSubmitter = sync.OnceValue[*updateUserMetadataPositionsSubmitter](newUpdateUserMetadataPositionsSubmitter)
)

// updateUserMetadataPositionsSubmitter represents update user metadata positions submitter.
type updateUserMetadataPositionsSubmitter struct {
	// Represents general context used to manage submitted context.
	ctx context.Context

	// Represents channel, which is used to close the submitted action.
	cancel context.CancelFunc

	// Represents previously retrieved position value.
	previousPosition dto.Position
}

// close performs stream submitter close operation.
func (uumps *updateUserMetadataPositionsSubmitter) close() {
	if uumps.ctx != nil {
		select {
		case <-uumps.ctx.Done():
		default:
			uumps.cancel()
		}
	}
}

// Submit performs a submittion of lobby set retrieval action. Callback is
// required to return boolean value, which defines whether submitter should be closed
// or not.
func (uumps *updateUserMetadataPositionsSubmitter) Submit(
	lobbyID int64, callback func(err error) bool) {
	uumps.ctx, uumps.cancel = context.WithCancel(context.Background())

	go func() {
		ticker := time.NewTicker(updateUserMetadataPositionsFrequency)

		for {
			select {
			case <-ticker.C:
				position := store.GetPositionSession()

				if position == uumps.previousPosition {
					continue
				}

				fmt.Println(position)

				message, err := proto.Marshal(&contentv1.UpdateUserMetadataPositionsRequest{
					Issuer:  store.GetRepositoryUUID(),
					LobbyId: lobbyID,
					Position: &contentv1.Position{
						X: position.X,
						Y: position.Y,
					},
				})
				if err != nil {
					if callback(err) {
						uumps.close()
					}

					return
				}

				err = handler.GetInstance().Send(contentv1.UPDATE_USER_METADATA_POSITIONS, message)
				if err != nil {
					if callback(err) {
						uumps.close()
					}

					return
				}

				uumps.previousPosition = position

				if callback(nil) {
					uumps.close()
				}
			case <-uumps.ctx.Done():
				return
			}
		}
	}()
}

// Clean perform delayed submitter close operation, which results in a called
// provided callback when operation is finished.
func (uumps *updateUserMetadataPositionsSubmitter) Clean(callback func()) {
	go func() {
		uumps.close()

		callback()
	}()
}

// newUpdateUserMetadataPositionsSubmitter initializes updateUserMetadataPositionsSubmitter.
func newUpdateUserMetadataPositionsSubmitter() *updateUserMetadataPositionsSubmitter {
	return new(updateUserMetadataPositionsSubmitter)
}
