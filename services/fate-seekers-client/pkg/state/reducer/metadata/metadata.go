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
	RETRIEVED_SESSIONS_METADATA_STATE = "retrieved_sessions"
	SELECTED_SESSION_METADATA_STATE   = "selected_session"
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
