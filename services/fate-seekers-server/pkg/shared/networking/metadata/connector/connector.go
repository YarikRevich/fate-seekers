package connector

import (
	"fmt"
	"net"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	metadatav1 "github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/metadata/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/metadata/handler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/metadata/middleware"
	"google.golang.org/grpc"
)

// NetworkingMetadataConnector represents networking metadata connector.
type NetworkingMetadataConnector struct {
	// Represents established connection instance.
	conn net.Listener
}

func (nmc *NetworkingMetadataConnector) Connect(callback func(err error)) {
	go func() {
		var err error

		nmc.conn, err = net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", config.GetSettingsNetworkingServerPort()))
		if err != nil {
			callback(err)

			return
		}

		grpcServer := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				middleware.CheckValidationMiddleware,
				middleware.CheckAuthenticationMiddleware,
			),
		)

		metadatav1.RegisterMetadataServiceServer(grpcServer, handler.NewHandler())

		callback(nil)

		grpcServer.Serve(nmc.conn)
	}()
}

func (nmc *NetworkingMetadataConnector) Close() error {
	return nmc.conn.Close()
}

// NewNetworkingMetadataConnector initializes NetworkingMetadataConnector.
func NewNetworkingMetadataConnector() *NetworkingMetadataConnector {
	return new(NetworkingMetadataConnector)
}
