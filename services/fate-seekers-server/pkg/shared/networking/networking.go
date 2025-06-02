package networking

import "errors"

var (
	ErrConnectorHostIsInvalid                 = errors.New("err happened during provided connector host validation")
	ErrConnectorConnectionEstablishmentFailed = errors.New("err happened during connection establishment")
)
