package networking

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/luisvinicius167/godux"
)

// Describes all the available networking reducer store states.
const (
	ENTRY_HANDSHAKE_STARTED_NETWORKING_STATE            = "entry_handshake_started"
	PING_CONNECTION_STARTED_NETWORKING_STATE            = "ping_connection_started"
	SESSION_RETRIEVAL_STARTED_NETWORKING_STATE          = "session_retrieval_started"
	SESSION_CREATION_STARTED_NETWORKING_STATE           = "session_creation_started"
	SESSION_REMOVAL_STARTED_NETWORKING_STATE            = "session_removal_started"
	LOBBY_SET_RETRIEVAL_STARTED_NETWORKING_STATE        = "lobby_set_retrieval_started"
	LOBBY_SET_RETRIEVAL_CYCLE_FINISHED_NETWORKING_STATE = "lobby_set_retrieval_cycle_finished"
	LOBBY_CREATION_STARTED_NETWORKING_STATE             = "lobby_creation_started"
	LOBBY_REMOVAL_STARTED_NETWORKING_STATE              = "lobby_removal_started"
	SESSION_METADATA_RETRIEVAL_STARTED_NETWORKING_STATE = "session_metadata_retrieval_started"
)

// NetworkingStateReducer represents reducer used for networking state management.
type NetworkingStateReducer struct {
	// Represents of instance of state store.
	store *godux.Store
}

func (nsr *NetworkingStateReducer) Init() {
	nsr.store.SetState(
		ENTRY_HANDSHAKE_STARTED_NETWORKING_STATE, value.ENTRY_HANDSHAKE_STARTED_NETWORKING_FALSE_VALUE)
	nsr.store.SetState(
		PING_CONNECTION_STARTED_NETWORKING_STATE, value.PING_CONNECTION_STARTED_NETWORKING_FALSE_VALUE)
	nsr.store.SetState(
		SESSION_RETRIEVAL_STARTED_NETWORKING_STATE, value.SESSION_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE)
	nsr.store.SetState(
		SESSION_CREATION_STARTED_NETWORKING_STATE, value.SESSION_CREATION_STARTED_NETWORKING_FALSE_VALUE)
	nsr.store.SetState(
		SESSION_REMOVAL_STARTED_NETWORKING_STATE, value.ENTRY_HANDSHAKE_STARTED_NETWORKING_FALSE_VALUE)
	nsr.store.SetState(
		LOBBY_SET_RETRIEVAL_STARTED_NETWORKING_STATE, value.LOBBY_SET_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE)
	nsr.store.SetState(
		LOBBY_SET_RETRIEVAL_CYCLE_FINISHED_NETWORKING_STATE, value.LOBBY_SET_RETRIEVAL_CYCLE_FINISHED_NETWORKING_FALSE_VALUE)
	nsr.store.SetState(
		LOBBY_CREATION_STARTED_NETWORKING_STATE, value.LOBBY_CREATION_STARTED_NETWORKING_FALSE_VALUE)
	nsr.store.SetState(
		LOBBY_REMOVAL_STARTED_NETWORKING_STATE, value.LOBBY_REMOVAL_STARTED_NETWORKING_FALSE_VALUE)
	nsr.store.SetState(
		SESSION_METADATA_RETRIEVAL_STARTED_NETWORKING_STATE,
		value.SESSION_METADATA_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE)
}

func (nsr *NetworkingStateReducer) GetProcessor() func(value godux.Action) interface{} {
	return func(value godux.Action) interface{} {
		switch value.Type {
		case action.SET_ENTRY_HANDSHAKE_STARTED_NETWORKING_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{
					Key: ENTRY_HANDSHAKE_STARTED_NETWORKING_STATE, Value: value.Value})

		case action.SET_PING_CONNECTION_STARTED_NETWORKING_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{
					Key: PING_CONNECTION_STARTED_NETWORKING_STATE, Value: value.Value})

		case action.SET_SESSION_RETRIEVAL_STARTED_NETWORKING_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{
					Key: SESSION_RETRIEVAL_STARTED_NETWORKING_STATE, Value: value.Value})

		case action.SET_SESSION_CREATION_STARTED_NETWORKING_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{
					Key: SESSION_CREATION_STARTED_NETWORKING_STATE, Value: value.Value})

		case action.SET_SESSION_REMOVAL_STARTED_NETWORKING_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{
					Key: SESSION_REMOVAL_STARTED_NETWORKING_STATE, Value: value.Value})

		case action.SET_LOBBY_SET_RETRIEVAL_STARTED_NETWORKING_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{
					Key: LOBBY_SET_RETRIEVAL_STARTED_NETWORKING_STATE, Value: value.Value})

		case action.SET_LOBBY_SET_RETRIEVAL_CYCLE_FINISHED_NETWORKING_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{
					Key: LOBBY_SET_RETRIEVAL_CYCLE_FINISHED_NETWORKING_STATE, Value: value.Value})

		case action.SET_LOBBY_CREATION_STARTED_NETWORKING_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{
					Key: LOBBY_CREATION_STARTED_NETWORKING_STATE, Value: value.Value})

		case action.SET_LOBBY_REMOVAL_STARTED_NETWORKING_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{
					Key: LOBBY_REMOVAL_STARTED_NETWORKING_STATE, Value: value.Value})

		case action.SET_SESSION_METADATA_RETRIEVAL_STARTED_NETWORKING_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{
					Key: SESSION_METADATA_RETRIEVAL_STARTED_NETWORKING_STATE, Value: value.Value})

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
