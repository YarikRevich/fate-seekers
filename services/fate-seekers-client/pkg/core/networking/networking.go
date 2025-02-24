package networking

// NetworkingConnector represents networking connector interface.
type NetworkingConnector interface {
	// Reconnect performs networking connector reconnection operation.
	Reconnect()
}
