package networking

import "errors"

var (
	ErrConnectorHostIsInvalid                 = errors.New("err happened during provided connector host validation")
	ErrConnectorConnectionEstablishmentFailed = errors.New("err happened during connection establishment")
)

// NetworkingConnector represents networking connector interface.
type NetworkingConnector interface {
	// Connect performs networking connector connection operation using latest config properties.
	// Can be used for reconnection as well.
	Connect() error

	// Ping performs health check operation for the established connection.
	Ping() bool

	// Close performs networking connector connection close operation.
	Close() error
}
