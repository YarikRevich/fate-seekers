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
	// Represents networking content connector.
	contentConnector *contentconnector.NetworkingContentConnector

	// Represents networking metadata connector.
	metadataConnector *metadataconnector.NetworkingMetadataConnector
}

// Connect performs a connection for all the API modules.
func (gnc *GlobalNetworkingConnector) Connect(callback func(err error), failover func(err error)) {
	go func() {
		err := gnc.contentConnector.Connect(failover)
		if err != nil {
			callback(err)

			return
		}

		err = gnc.metadataConnector.Connect()
		if err != nil {
			if err := gnc.contentConnector.Close(); err != nil {
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
		err := gnc.contentConnector.Close()
		if err != nil {
			callback(err)

			return
		}

		callback(gnc.metadataConnector.Close())
	}()
}

// newGlobalNetworkingConnector initializes GlobalNetworkingConnector.
func newGlobalNetworkingConnector() *GlobalNetworkingConnector {
	return &GlobalNetworkingConnector{
		contentConnector:  contentconnector.NewNetworkingContentConnector(),
		metadataConnector: metadataconnector.NewNetworkingMetadataConnector(),
	}
}
