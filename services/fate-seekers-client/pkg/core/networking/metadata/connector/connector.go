package connector

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking"
	metadatav1 "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/metadata/middleware"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/logging"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	// GetInstance retrieves instance of the networking content connector, performing initilization if needed.
	GetInstance = sync.OnceValue[*NetworkingMetadataConnector](newNetworkingMetadataConnector)
)

// NetworkingMetadataConnector represents networking metadata connector.
type NetworkingMetadataConnector struct {
	// Represents established connection instance.
	conn *grpc.ClientConn

	// Represents metadata connector client.
	client metadatav1.MetadataServiceClient
}

// Connecto performs connect operation.
func (nmc *NetworkingMetadataConnector) Connect() error {
	var err error

	nmc.conn, err = grpc.Dial(
		config.GetSettingsNetworkingServerHost(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(&middleware.AuthenticationMiddleware{}),
		grpc.WithUnaryInterceptor(middleware.CheckValidationMiddleware))
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

	nmc.client = metadatav1.NewMetadataServiceClient(nmc.conn)

	return nil
}

// Close performs close operation.
func (nmc *NetworkingMetadataConnector) Close() error {
	if nmc.conn != nil {
		return nmc.conn.Close()
	}

	return nil
}

// GetConnection retrieves metadata client connection instance.
func (ncm *NetworkingMetadataConnector) GetClient() metadatav1.MetadataServiceClient {
	return ncm.client
}

// newNetworkingMetadataConnector initializes NetworkingMetadataConnector.
func newNetworkingMetadataConnector() *NetworkingMetadataConnector {
	return new(NetworkingMetadataConnector)
}
