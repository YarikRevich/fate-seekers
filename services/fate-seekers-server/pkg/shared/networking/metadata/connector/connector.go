package connector

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/metadata/api"
	"google.golang.org/grpc"
)

// NetworkingMetadataConnector represents networking metadata connector.
type NetworkingMetadataConnector struct {
	// Represents establushed connection instance.
	conn *grpc.ClientConn

	// Represents established client connection instance.
	client api.MetadataClient
}

func (nmc *NetworkingMetadataConnector) Connect() error {
	// conn, err := grpc.NewClient(
	// 	config.GetSettingsNetworkingServerHost(),
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
	// 	grpc.WithPerRPCCredentials(&middleware.AuthenticationMiddleware{}))
	// if err != nil {
	// 	return errors.Wrap(err, networking.ErrConnectorHostIsInvalid.Error())
	// }

	// sigc := make(chan os.Signal, 1)
	// signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// go func() {
	// 	select {
	// 	case <-sigc:
	// 		if err := conn.Close(); err != nil {
	// 			logging.GetInstance().Fatal(err.Error())
	// 		}
	// 	}
	// }()

	// nmc.conn = conn

	// nmc.client = api.NewMetadataClient(conn)

	return nil
}

func (nmc *NetworkingMetadataConnector) Ping() bool {
	// nmc.client.PingConnection()

	return false
}

func (nmc *NetworkingMetadataConnector) Close() error {
	return nmc.conn.Close()
}

// NewNetworkingMetadataConnector initializes NetworkingMetadataConnector.
func NewNetworkingMetadataConnector() networking.NetworkingConnector {
	return new(NetworkingMetadataConnector)
}
