package metadata

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/luisvinicius167/godux"
)

// Describes all the available metadata reducer store states.
const (
	RETRIEVED_SESSIONS_METADATA_STATE      = "retrieved_sessions"
	SELECTED_SESSION_METADATA_STATE        = "selected_session"
	RETRIEVED_LOBBY_SET_METADATA_STATE     = "retrieved_lobby_set"
	SELECTED_LOBBY_SET_UNIT_METADATA_STATE = "selected_lobby_set_unit"
	SESSION_ALREADY_STARTED_METADATA_STATE = "session_already_started"
)

// MetadataStateReducer represents reducer used for metadata state management.
type MetadataStateReducer struct {
	// Represents of instance of state store.
	store *godux.Store
}

func (msr *MetadataStateReducer) Init() {
	msr.store.SetState(
		RETRIEVED_SESSIONS_METADATA_STATE, value.RETRIEVED_SESSIONS_METADATA_EMPTY_VALUE)
	msr.store.SetState(
		SELECTED_SESSION_METADATA_STATE, value.SELECTED_SESSION_METADATA_EMPTY_VALUE)
	msr.store.SetState(
		RETRIEVED_LOBBY_SET_METADATA_STATE, value.RETRIEVED_LOBBY_SET_METADATA_EMPTY_VALUE)
	msr.store.SetState(
		SELECTED_LOBBY_SET_UNIT_METADATA_STATE, value.SELECTED_LOBBY_SET_UNIT_METADATA_EMPTY_VALUE)
	msr.store.SetState(
		SESSION_ALREADY_STARTED_METADATA_STATE, value.SESSION_ALREADY_STARTED_METADATA_STATE_FALSE_VALUE)
}

func (msr *MetadataStateReducer) GetProcessor() func(value godux.Action) interface{} {
	return func(value godux.Action) interface{} {
		switch value.Type {
		case action.SET_RETRIEVED_SESSIONS_METADATA_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{
					Key: RETRIEVED_SESSIONS_METADATA_STATE, Value: value.Value})

		case action.SET_SELECTED_SESSION_METADATA_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{
					Key: SELECTED_SESSION_METADATA_STATE, Value: value.Value})

		case action.SET_RETRIEVED_LOBBY_SET_METADATA_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{
					Key: RETRIEVED_LOBBY_SET_METADATA_STATE, Value: value.Value})

		case action.SET_SELECTED_LOBBY_SET_UNIT_METADATA_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{
					Key: SELECTED_LOBBY_SET_UNIT_METADATA_STATE, Value: value.Value})

		case action.SET_SESSION_ALREADY_STARTED_METADATA_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{
					Key: SESSION_ALREADY_STARTED_METADATA_STATE, Value: value.Value})

		default:
			return nil
		}
	}
}

// NewMetadataStateReducer initializes new instance of MetadataStateReducer.
func NewMetadataStateReducer(store *godux.Store) reducer.Reducer {
	return &MetadataStateReducer{
		store: store,
	}
}
