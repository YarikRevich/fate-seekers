package connector

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/balacode/udpt"
	"golang.org/x/crypto/blake2b"
)

// NetworkingContentConnector represents networking content connector.
type NetworkingContentConnector struct {
	// Represents context for initialized receiver.
	close context.CancelFunc
}

// Connect performs connection attempt. Requires failover callback which is used for the
// case when initialized connection is interrupted by some error.
func (ncc *NetworkingContentConnector) Connect(failover func(err error)) error {
	networkingReceiverPortInt, err := strconv.Atoi(config.GetSettingsNetworkingReceiverPort())
	if err != nil {
		return err
	}

	networkingEncryptionKeyHash, err := blake2b.New256([]byte(config.GetSettingsNetworkingEncryptionKey()))
	if err != nil {
		return err
	}

	var ctx context.Context

	ctx, ncc.close = context.WithCancel(context.Background())

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
			ncc.close()

			failover(err)
		}
	}()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		select {
		case <-sigc:
			ncc.close()
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

// NewNetworkingContentConnector initializes NetworkingContentConnector.
func NewNetworkingContentConnector() *NetworkingContentConnector {
	return new(NetworkingContentConnector)
}
