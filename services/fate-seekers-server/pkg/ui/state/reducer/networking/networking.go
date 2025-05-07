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
	SESSION_HANDSHAKE_STARTED_NETWORKING_STATE = "session_handshake_started"
	ENTRY_HANDSHAKE_STARTED_NETWORKING_STATE   = "entry_handshake_started"
)

// NetworkingStateReducer represents reducer used for networking state management.
type NetworkingStateReducer struct {
	// Represents of instance of state store.
	store *godux.Store
}

func (nsr *NetworkingStateReducer) Init() {
	nsr.store.SetState(
		ENTRY_HANDSHAKE_STARTED_NETWORKING_STATE, value.ENTRY_HANDSHAKE_STARTED_NETWORKING_FALSE_VALUE)
}

func (nsr *NetworkingStateReducer) GetProcessor() func(value godux.Action) interface{} {
	return func(value godux.Action) interface{} {
		switch value.Type {
		case action.SET_ENTRY_HANDSHAKE_STARTED_NETWORKING_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{
					Key: ENTRY_HANDSHAKE_STARTED_NETWORKING_STATE, Value: value.Value})

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
