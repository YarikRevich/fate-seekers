package connector

import (
	"net"
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking"
)

var (
	// GetInstance retrieves instance of the netowkring content connector, performing initilization if needed.
	GetInstance = sync.OnceValue[networking.NetworkingConnector](newNetworkingContentConnector)
)

// NetworkingContentConnector represents networking content connector.
type NetworkingContentConnector struct {
}

func (ncc *NetworkingContentConnector) Reconnect() {
	addr := &net.UDPAddr{
		Port: 1234,
		IP:   net.ParseIP("127.0.0.1"),
	}

	conn, err := net.ListenUDP("udp", addr)

	// notification.GetInstance().Push("Тестове повідомлення!", time.Second*6)
	// 	notification.GetInstance().Push("Друге повідомлення!", time.Second*6)
}

// newNetworkingContentConnector initializes NetworkingContentConnector.
func newNetworkingContentConnector() networking.NetworkingConnector {
	return new(NetworkingContentConnector)
}
