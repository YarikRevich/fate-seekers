package handler

import (
	"errors"
	"sync"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/balacode/udpt"
)

var (
	// GetInstance retrieves instance of the content handler, performing initilization if needed.
	GetInstance = sync.OnceValue[*ContentHandler](newContentHandler)
)

var (
	ErrSendRequest = errors.New("err happened during request send operation")
)

// ContentHandler represents content data handler.
type ContentHandler struct {
	// Represents global UDP configuration used for send requests.
	configuration *udpt.Configuration
}

// send performs low level UDP send operation using udpt wrapper.
func (ch *ContentHandler) send(key string, value []byte) error {
	err := udpt.Send(
		config.GetSettingsNetworkingServerHost(),
		"tjkjtkr",
		[]byte("jgkfjgkfjg"),
		[]byte(config.GetSettingsNetworkingEncryptionKey()),
		ch.configuration)
	if err != nil {
		return ErrSendRequest
	}

	return nil
}

// PerformPingConnection performs ping connection
func (ch *ContentHandler) PerformPingConnection(callback func(err error)) {
	go func() {
		err := ch.send("jkj", []byte("gjfkgjfk"))
		if err != nil {
			callback(err)

			return
		}

		callback(nil)
	}()
}

// newContentHandler initializes ContentHandler.
func newContentHandler() *ContentHandler {
	configuration := udpt.NewDefaultConfig()

	configuration.ReplyTimeout = 100 * time.Millisecond

	configuration.SendRetries = 2
	configuration.WriteTimeout = 25 * time.Millisecond

	return &ContentHandler{
		configuration: configuration,
	}
}
