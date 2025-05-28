package connector

import (
	"os"
	"os/signal"
	"sync"
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

var (
	// GetInstance retrieves instance of the networking metadata connector, performing initilization if needed.
	GetInstance = sync.OnceValue[networking.NetworkingConnector](newNetworkingMetadataConnector)
)

// NetworkingMetadataConnector represents networking metadata connector.
type NetworkingMetadataConnector struct {
	// Represents established connection instance.
	conn *grpc.ClientConn

	// Represents established client connection instance.
	client api.MetadataClient
}

func (nmc *NetworkingMetadataConnector) Connect() error {
	conn, err := grpc.NewClient(
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
			if err := conn.Close(); err != nil {
				logging.GetInstance().Fatal(err.Error())
			}
		}
	}()

	nmc.conn = conn

	nmc.client = api.NewMetadataClient(conn)

	return nil
}

func (nmc *NetworkingMetadataConnector) Ping() bool {
	// nmc.client.PingConnection()

	return false
}

func (nmc *NetworkingMetadataConnector) Close() error {
	return nmc.conn.Close()
}

// newNetworkingMetadataConnector initializes NetworkingMetadataConnector.
func newNetworkingMetadataConnector() networking.NetworkingConnector {
	return new(NetworkingMetadataConnector)
}
