package handler

import (
	"errors"
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/content/middleware"
	"github.com/balacode/udpt"
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
func (h *Handler) Send(key string, value []byte) error {
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
