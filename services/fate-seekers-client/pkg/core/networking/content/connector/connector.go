package connector

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/logging"
	"github.com/balacode/udpt"
	"golang.org/x/crypto/blake2b"
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
	networkingReceiverPortInt, err := strconv.Atoi(config.GetSettingsNetworkingReceiverPort())
	if err != nil {
		return err
	}

	networkingEncryptionKeyHash, err := blake2b.New256([]byte(config.GetSettingsNetworkingEncryptionKey()))
	if err != nil {
		return err
	}

	ctx, close := context.WithCancel(context.Background())

	go func() {
		err := udpt.Receive(
			ctx,
			networkingReceiverPortInt,
			networkingEncryptionKeyHash.Sum(nil),
			func(k string, v []byte) error {
				fmt.Println(k, string(v))

				return nil
			})
		if err != nil {
			close()

			logging.GetInstance().Fatal(err.Error())
		}
	}()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		select {
		case <-sigc:
			close()
		}
	}()

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
