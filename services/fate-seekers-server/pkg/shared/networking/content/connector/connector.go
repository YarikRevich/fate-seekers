package connector

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/logging"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking"
	"github.com/balacode/udpt"
)

var (
	// GetInstance retrieves instance of the networking content connector, performing initilization if needed.
	GetInstance = sync.OnceValue[networking.NetworkingConnector](newNetworkingContentConnector)
)

// NetworkingContentConnector represents networking content connector.
type NetworkingContentConnector struct {
	// Represents context for initialized receiver.
	close context.CancelFunc
}

func (ncc *NetworkingContentConnector) Connect() error {
	networkingServerPortInt, err := strconv.Atoi(config.GetSettingsNetworkingServerPort())
	if err != nil {
		return err
	}

	ctx, close := context.WithCancel(context.Background())

	go func(ctx context.Context, close context.CancelFunc) {
		err := udpt.Receive(
			ctx,
			networkingServerPortInt,
			[]byte(config.GetSettingsNetworkingEncryptionKey()),
			func(k string, v []byte) error {
				return nil
			})
		if err != nil {
			close()

			logging.GetInstance().Fatal(err.Error())
		}
	}(ctx, close)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func(close context.CancelFunc) {
		select {
		case <-sigc:
			close()
		}
	}(close)

	return nil
}

func (ncc *NetworkingContentConnector) Ping() bool {
	return false
}

func (ncc *NetworkingContentConnector) Close() error {
	ncc.close()

	return nil
}

// newNetworkingContentConnector initializes NetworkingContentConnector.
func newNetworkingContentConnector() networking.NetworkingConnector {
	return new(NetworkingContentConnector)
}
