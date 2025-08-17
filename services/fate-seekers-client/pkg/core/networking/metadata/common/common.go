package common

import "errors"

// Describes all the available errors used for metadata connector.
var (
	ErrConnectionLost = errors.New("err happened connection with server lost")
)
