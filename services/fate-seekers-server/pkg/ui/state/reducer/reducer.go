package reducer

import "github.com/luisvinicius167/godux"

// Reducer represents common state reducer interface.
type Reducer interface {
	// Init performs initial reducer state setup.
	Init()

	// GetProcessor retrieves reducer processor logic.
	GetProcessor() func(value godux.Action) interface{}
}
