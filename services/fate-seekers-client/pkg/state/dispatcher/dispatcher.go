package dispatcher

import (
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"github.com/luisvinicius167/godux"
)

var (
	// GetInstance retrieves instance of the dispatcher, performing initilization if needed.
	GetInstance = sync.OnceValue[*Dispatcher](newDispatcher)
)

// Dispatcher represents dispatcher wrapper initialization.
type Dispatcher struct {
}

// Dispatch performs state dispatch operation.
func (d *Dispatcher) Dispatch(action godux.Action) {
	valueRaw := store.GetInstance().Dispatch(action)

	reducerResult := valueRaw.(dto.ReducerResult)

	for key, value := range reducerResult {
		store.GetInstance().SetState(key, value)
	}
}

// newDispatcher creates new instance of Dispatcher.
func newDispatcher() *Dispatcher {
	return new(Dispatcher)
}
