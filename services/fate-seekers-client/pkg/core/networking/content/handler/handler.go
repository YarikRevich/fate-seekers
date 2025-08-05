package handler

import (
	"errors"
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/content/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/content/middleware"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/balacode/udpt"
	"google.golang.org/protobuf/proto"
)

var (
	// GetInstance retrieves instance of the handler, performing initilization if needed.
	GetInstance = sync.OnceValue[*Handler](newHandler)
)

var (
	ErrSendRequest = errors.New("err happened during request send operation")
)

// Handler represents content data handler.
type Handler struct {
	// Represents global UDP configuration used for send requests.
	configuration *udpt.Configuration
}

// send performs low level UDP send operation using udpt wrapper.
func (h *Handler) send(key string, value []byte) error {
	return middleware.
		GetInstance().
		Run(func() error {
			err := udpt.Send(
				config.GetSettingsNetworkingServerHost(),
				key,
				value,
				config.GetSettingsParsedNetworkingEncryptionKey(),
				h.configuration)
			if err != nil {
				return ErrSendRequest
			}

			return nil
		})
}

// PerformGetUserMetadataPositions performs get user positions retrieval.
func (h *Handler) PerformGetUserMetadataPositions(callback func(err error)) {
	go func() {
		message, err := proto.Marshal(&api.GetUserMetadataPositionsRequest{
			Issuer:    store.GetRepositoryUUID(),
			SessionId: "",
		})
		if err != nil {
			callback(err)

			return
		}

		err = h.send(api.GET_USER_METADATA_POSITIONS, message)
		if err != nil {
			callback(err)

			return
		}

		callback(nil)
	}()
}

// PerformUpdateUserMetadataPositions performs user positions update.
func (h *Handler) PerformUpdateUserMetadataPositions(callback func(err error)) {
	go func() {
		message, err := proto.Marshal(&api.UpdateUserMetadataPositionsRequest{
			Issuer:    store.GetRepositoryUUID(),
			SessionId: "",
			Position:  &api.Position{},
		})
		if err != nil {
			callback(err)

			return
		}

		err = h.send(api.UPDATE_USER_METADATA_POSITIONS, message)
		if err != nil {
			callback(err)

			return
		}

		callback(nil)
	}()
}

// newHandler initializes Handler.
func newHandler() *Handler {
	configuration := udpt.NewDefaultConfig()

	configuration.ReplyTimeout = 100 * time.Millisecond

	configuration.SendRetries = 2
	configuration.WriteTimeout = 25 * time.Millisecond

	return &Handler{
		configuration: configuration,
	}
}
