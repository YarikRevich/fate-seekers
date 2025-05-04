package connector

import (
	"sync"

	contentconnector "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/content/connector"
	metadataconnector "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/connector"
)

var (
	// GetInstance retrieves instance of the networking content connector, performing initilization if needed.
	GetInstance = sync.OnceValue[*GlobalNetworkingConnector](newGlobalNetworkingConnector)
)

// GlobalNetworkingConnector represents global networking connector.
type GlobalNetworkingConnector struct {
}

// Connect performs a connection for all the API modules.
func (gnc *GlobalNetworkingConnector) Connect(callback func(err error)) {
	go func() {
		err := contentconnector.GetInstance().Connect()
		if err != nil {
			callback(err)

			return
		}

		err = metadataconnector.GetInstance().Connect()
		if err != nil {
			if err := contentconnector.GetInstance().Close(); err != nil {
				callback(err)

				return
			}

			callback(err)

			return
		}

		callback(nil)
	}()
}

func (gnc *GlobalNetworkingConnector) Close(callback func(err error)) {
	go func() {
		err := contentconnector.GetInstance().Close()
		if err != nil {
			callback(err)

			return
		}

		callback(metadataconnector.GetInstance().Close())
	}()
}

// newGlobalNetworkingConnector initializes GlobalNetworkingConnector.
func newGlobalNetworkingConnector() *GlobalNetworkingConnector {
	return new(GlobalNetworkingConnector)
}
