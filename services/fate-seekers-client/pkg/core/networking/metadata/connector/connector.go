package connector

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/middleware"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/logging"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NetworkingMetadataConnector represents networking metadata connector.
type NetworkingMetadataConnector struct {
	// Represents established connection instance.
	conn *grpc.ClientConn

	// Represents established client connection instance.
	client api.MetadataClient
}

func (nmc *NetworkingMetadataConnector) Connect() error {
	var err error

	nmc.conn, err = grpc.NewClient(
		config.GetSettingsNetworkingServerHost(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(&middleware.AuthenticationMiddleware{}))
	if err != nil {
		return errors.Wrap(err, networking.ErrConnectorHostIsInvalid.Error())
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		select {
		case <-sigc:
			if err := nmc.conn.Close(); err != nil {
				logging.GetInstance().Fatal(err.Error())
			}
		}
	}()

	nmc.client = api.NewMetadataClient(nmc.conn)

	return nil
}

func (nmc *NetworkingMetadataConnector) Close() error {
	if nmc.conn != nil {
		return nmc.conn.Close()
	}

	return nil
}

// NewNetworkingMetadataConnector initializes NetworkingMetadataConnector.
func NewNetworkingMetadataConnector() *NetworkingMetadataConnector {
	return new(NetworkingMetadataConnector)
}
