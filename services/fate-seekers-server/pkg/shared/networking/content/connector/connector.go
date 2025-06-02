package connector

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/logging"
	"github.com/balacode/udpt"
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

	var ctx context.Context

	ctx, ncc.close = context.WithCancel(context.Background())

	go func(ctx context.Context, close context.CancelFunc) {
		err := udpt.Receive(
			ctx,
			networkingServerPortInt,
			config.GetSettingsParsedNetworkingEncryptionKey(),
			func(k string, v []byte) error {
				fmt.Println(k, string(v))

				return nil
			})
		if err != nil {
			close()

			logging.GetInstance().Fatal(err.Error())
		}
	}(ctx, ncc.close)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func(close context.CancelFunc) {
		select {
		case <-sigc:
			close()
		}
	}(ncc.close)

	return nil
}

func (ncc *NetworkingContentConnector) Close() error {
	ncc.close()

	return nil
}

// NewNetworkingContentConnector initializes NetworkingContentConnector.
func NewNetworkingContentConnector() *NetworkingContentConnector {
	return new(NetworkingContentConnector)
}
