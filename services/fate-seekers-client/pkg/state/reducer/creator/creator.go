package creator

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/luisvinicius167/godux"
)

// Describes all the available creator reducer store states.
const (
	INPUTS_UPDATED_STATE = "inputs_updated"
)

// CreatorStateReducer represents reducer used for networking state management.
type CreatorStateReducer struct {
	// Represents of instance of state store.
	store *godux.Store
}

func (csr *CreatorStateReducer) Init() {
	nsr.store.SetState(
		ENTRY_HANDSHAKE_STARTED_NETWORKING_STATE, value.ENTRY_HANDSHAKE_STARTED_NETWORKING_FALSE_VALUE)
	nsr.store.SetState(
		PING_CONNECTION_STARTED_NETWORKING_STATE, value.PING_CONNECTION_STARTED_NETWORKING_FALSE_VALUE)
	nsr.store.SetState(
		SESSION_RETRIEVAL_STARTED_NETWORKING_STATE, value.SESSION_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE)
	nsr.store.SetState(
		SESSION_CREATION_STARTED_NETWORKING_STATE, value.SESSION_CREATION_STARTED_NETWORKING_FALSE_VALUE)
	nsr.store.SetState(
		SESSION_JOINING_STARTED_NETWORKING_STATE, value.SESSION_JOINING_STARTED_NETWORKING_FALSE_VALUE)
	nsr.store.SetState(
		SESSION_REMOVAL_STARTED_NETWORKING_STATE, value.ENTRY_HANDSHAKE_STARTED_NETWORKING_FALSE_VALUE)
}

func (csr *CreatorStateReducer) GetProcessor() func(value godux.Action) interface{} {
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

		case action.SET_SESSION_JOINING_STARTED_NETWORKING_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{
					Key: SESSION_JOINING_STARTED_NETWORKING_STATE, Value: value.Value})

		case action.SET_SESSION_REMOVAL_STARTED_NETWORKING_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{
					Key: SESSION_REMOVAL_STARTED_NETWORKING_STATE, Value: value.Value})

		default:
			return nil
		}
	}
}

// NewCreatorStateReducer initializes new instance of CreatorStateReducer.
func NewCreatorStateReducer(store *godux.Store) reducer.Reducer {
	return &CreatorStateReducer{
		store: store,
	}
}
