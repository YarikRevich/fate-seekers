package stream

import (
	"context"
	"errors"
	"sync"

	metadatav1 "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/connector"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// GetUpdateSessionsActivitySubmitter retrieves instance of the update sessions activity submitter, performing initial creation if needed.
	GetUpdateSessionsActivitySubmitter = sync.OnceValue[*updateSessionsActivitySubmitter](newUpdateSessionActivitySubmitter)

	// GetGetSessionMetadataSubmitter retrieves instance of the sessions metadata retrieval submitter, performing initial creation if needed.
	GetGetSessionMetadataSubmitter = sync.OnceValue[*getSessionMetadataSubmitter](newGetSessionMetadataSubmitter)

	// GetGetLobbySetSubmitter retrieves instance of the lobby set retrieval submitter, performing initial creation if needed.
	GetGetLobbySetSubmitter = sync.OnceValue[*getLobbySetMetadataSubmitter](newGetLobbySetSubmitter)
)

// updateSessionsActivitySubmitter represents update sessions activity submitter.
type updateSessionsActivitySubmitter struct {
	// Represents general context used to manage submitted context.
	ctx context.Context

	// Represents channel, which is used to close the submitted action.
	cancel context.CancelFunc
}

// // TODO: create stream handler. Return callback with response and close channel which can be used to close the connection.
// func SubmitUpdateSessionActivity(input func() *metadatav1.UpdateSessionActivityRequest, finish func(cancel func()), callback func(err error)) {
// 	go func() {
// 		stream, err := connector.
// 			GetInstance().
// 			GetClient().
// 			UpdateSessionActivity(context.Background())
// 		if err != nil {
// 			callback(err)

// 			return
// 		}

// 		context.WithCancel(context.Background())

// 		cancel := make(chan bool, 1)

// 		for {
// 			err = stream.Send(input())
// 			if err != nil {
// 				callback(err)

// 				return
// 			}

// 			select {
// 			case <-cancel:
// 				_, err := stream.CloseAndRecv()
// 				if err != nil {
// 					callback(nil, err)
// 				}

// 				close(cancel)
// 			}
// 		}
// 	}()
// }

// newUpdateSessionActivitySubmitter initializes updateSessionsActivitySubmitter.
func newUpdateSessionActivitySubmitter() *updateSessionsActivitySubmitter {
	return new(updateSessionsActivitySubmitter)
}

// getSessionMetadataSubmitter represents session metadata retrieval submitter.
type getSessionMetadataSubmitter struct {
	// Represents general context used to manage submitted context.
	ctx context.Context

	// Represents channel, which is used to close the submitted action.
	cancel context.CancelFunc
}

// close performs stream submitter close operation.
func (gsms *getSessionMetadataSubmitter) close() {
	if gsms.ctx != nil {
		select {
		case <-gsms.ctx.Done():
		default:
			gsms.cancel()
		}
	}
}

// Submit performs a submittion of session metadata retrieval action. Callback is
// required to return boolean value, which defines whether submitter should be closed
// or not.
func (gsms *getSessionMetadataSubmitter) Submit(sessionID int64, callback func(response *metadatav1.GetSessionMetadataResponse, err error) bool) {
	gsms.ctx, gsms.cancel = context.WithCancel(context.Background())

	go func() {
		stream, err := connector.
			GetInstance().
			GetClient().
			GetSessionMetadata(
				gsms.ctx,
				&metadatav1.GetSessionMetadataRequest{
					SessionId: sessionID,
					Issuer:    store.GetRepositoryUUID(),
				})
		if err != nil {
			if callback(nil, err) {
				gsms.close()
			}

			return
		}

		for {
			response, err := stream.Recv()
			if err != nil {
				if status.Code(err) == codes.Unavailable {
					dispatcher.
						GetInstance().
						Dispatch(
							action.NewSetStateResetApplicationAction(
								value.STATE_RESET_APPLICATION_TRUE_VALUE))

					dispatcher.GetInstance().Dispatch(
						action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))

					if callback(nil, common.ErrConnectionLost) {
						gsms.close()
					}

					return
				}

				errRaw, ok := status.FromError(err)
				if !ok {
					if callback(nil, err) {
						gsms.close()
					}

					return
				}

				if callback(nil, errors.New(errRaw.Message())) {
					gsms.close()
				}

				break
			}

			if callback(response, nil) {
				gsms.close()
			}
		}
	}()
}

// Clean perform delayed submitter close operation, which results in a called
// provided callback when operation is finished.
func (gsms *getSessionMetadataSubmitter) Clean(callback func()) {
	go func() {
		gsms.close()

		callback()
	}()
}

// newGetSessionMetadataSubmitter initializes getSessionMetadataSubmitter.
func newGetSessionMetadataSubmitter() *getSessionMetadataSubmitter {
	return new(getSessionMetadataSubmitter)
}

// getLobbySetMetadataSubmitter represents lobby set retrieval submitter.
type getLobbySetMetadataSubmitter struct {
	// Represents general context used to manage submitted context.
	ctx context.Context

	// Represents channel, which is used to close the submitted action.
	cancel context.CancelFunc
}

// close performs stream submitter close operation.
func (glsms *getLobbySetMetadataSubmitter) close() {
	if glsms.ctx != nil {
		select {
		case <-glsms.ctx.Done():
		default:
			glsms.cancel()
		}
	}
}

// Submit performs a submittion of lobby set retrieval action. Callback is
// required to return boolean value, which defines whether submitter should be closed
// or not.
func (glsms *getLobbySetMetadataSubmitter) Submit(sessionID int64, callback func(response *metadatav1.GetLobbySetResponse, err error) bool) {
	glsms.ctx, glsms.cancel = context.WithCancel(context.Background())

	go func() {
		stream, err := connector.
			GetInstance().
			GetClient().
			GetLobbySet(
				glsms.ctx,
				&metadatav1.GetLobbySetRequest{
					SessionId: sessionID,
					Issuer:    store.GetRepositoryUUID(),
				})
		if err != nil {
			if callback(nil, err) {
				glsms.close()
			}

			return
		}

		for {
			response, err := stream.Recv()

			if err != nil {
				if status.Code(err) == codes.Unavailable {
					dispatcher.
						GetInstance().
						Dispatch(
							action.NewSetStateResetApplicationAction(
								value.STATE_RESET_APPLICATION_TRUE_VALUE))

					dispatcher.GetInstance().Dispatch(
						action.NewSetActiveScreenAction(value.ACTIVE_SCREEN_MENU_VALUE))

					if callback(nil, common.ErrConnectionLost) {
						glsms.close()
					}

					return
				}

				errRaw, ok := status.FromError(err)
				if !ok {
					if callback(nil, err) {
						glsms.close()
					}

					return
				}

				if callback(nil, errors.New(errRaw.Message())) {
					glsms.close()
				}

				break
			}

			if callback(response, nil) {
				glsms.close()
			}
		}
	}()
}

// Clean perform delayed submitter close operation, which results in a called
// provided callback when operation is finished.
func (glsms *getLobbySetMetadataSubmitter) Clean(callback func()) {
	go func() {
		glsms.close()

		callback()
	}()
}

// newGetLobbySetSubmitter initializes getLobbySetMetadataSubmitter.
func newGetLobbySetSubmitter() *getLobbySetMetadataSubmitter {
	return new(getLobbySetMetadataSubmitter)
}
