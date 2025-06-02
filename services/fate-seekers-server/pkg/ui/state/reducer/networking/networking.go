package networking

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/reducer"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/value"
	"github.com/luisvinicius167/godux"
)

// Describes all the available networking reducer store states.
const (
	LISTENER_STARTED_NETWORKING_STATE = "listener_started"
)

// NetworkingStateReducer represents reducer used for networking state management.
type NetworkingStateReducer struct {
	// Represents of instance of state store.
	store *godux.Store
}

func (nsr *NetworkingStateReducer) Init() {
	nsr.store.SetState(
		LISTENER_STARTED_NETWORKING_STATE, value.LISTENER_STARTED_NETWORKING_STATE_FALSE_VALUE)
}

func (nsr *NetworkingStateReducer) GetProcessor() func(value godux.Action) interface{} {
	return func(value godux.Action) interface{} {
		switch value.Type {
		case action.SET_LISTENER_STARTED_NETWORKING_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{
					Key: LISTENER_STARTED_NETWORKING_STATE, Value: value.Value})

		default:
			return nil
		}
	}
}

// NewNetworkingStateReducer initializes new instance of NetworkingStateReducer.
func NewNetworkingStateReducer(store *godux.Store) reducer.Reducer {
	return &NetworkingStateReducer{
		store: store,
	}
}
