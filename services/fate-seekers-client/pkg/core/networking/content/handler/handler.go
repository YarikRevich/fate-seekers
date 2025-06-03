package handler

import (
	"errors"
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/content/api"
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
}

// PerformGetUserPosition performs get user positions retrieval
func (h *Handler) PerformGetUserPositions(callback func(err error)) {
	proto.Marshal(&api.GetUserPositionsRequest{
		Issuer:  "itworks",
		Session: "",
	})

	go func() {
		err := h.send(api.GET_USER_POSITIONS, []byte("gjfkgjfk"))
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
