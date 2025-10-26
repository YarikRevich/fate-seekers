package session

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/luisvinicius167/godux"
)

// Describes all the available session reducer store states.
const (
	POSITION_SESSION_STATE                 = "position"
	RETRIEVED_USERS_METADATA_SESSION_STATE = "retrieved_users_metadata"
	// SESSION_ALREADY_STARTED_METADATA_STATE = "session_already_started"
)

// SessionStateReducer represents reducer used for session state management.
type SessionStateReducer struct {
	// Represents of instance of state store.
	store *godux.Store
}

func (ssr *SessionStateReducer) Init() {
	ssr.store.SetState(POSITION_SESSION_STATE, value.POSITION_SESSION_EMPTY_VALUE)
	ssr.store.SetState(RETRIEVED_USERS_METADATA_SESSION_STATE, value.RETRIEVED_USERS_METADATA_EMPTY_VALUE)
}

func (ssr *SessionStateReducer) GetProcessor() func(value godux.Action) interface{} {
	return func(value godux.Action) interface{} {
		switch value.Type {
		case action.INCREMENT_X_POSITION_SESSION_ACTION:
			valueRaw := ssr.store.GetState(POSITION_SESSION_STATE).(dto.Position)
			valueRaw.X += 1

			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: POSITION_SESSION_STATE, Value: valueRaw})

		case action.INCREMENT_Y_POSITION_SESSION_ACTION:
			valueRaw := ssr.store.GetState(POSITION_SESSION_STATE).(dto.Position)
			valueRaw.Y += 1

			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: POSITION_SESSION_STATE, Value: valueRaw})

		case action.DECREMENT_X_POSITION_SESSION_ACTION:
			valueRaw := ssr.store.GetState(POSITION_SESSION_STATE).(dto.Position)
			valueRaw.X -= 1

			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: POSITION_SESSION_STATE, Value: valueRaw})

		case action.DECREMENT_Y_POSITION_SESSION_ACTION:
			valueRaw := ssr.store.GetState(POSITION_SESSION_STATE).(dto.Position)
			valueRaw.Y -= 1

			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: POSITION_SESSION_STATE, Value: valueRaw})

		case action.SET_RETRIEVED_USERS_METADATA_SESSION_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: RETRIEVED_USERS_METADATA_SESSION_STATE, Value: value.Value})

		default:
			return nil
		}
	}
}

// NewSessionStateReducer initializes new instance of SessionStateReducer.
func NewSessionStateReducer(store *godux.Store) reducer.Reducer {
	return &SessionStateReducer{
		store: store,
	}
}
