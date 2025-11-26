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
	RESET_SESSION_STATE                    = "reset_session"
	STATIC_SESSION_STATE                   = "static"
	POSITION_SESSION_STATE                 = "position"
	STAGE_POSITION_SESSION_STATE           = "stage_position"
	PREVIOUS_POSITION_SESSION_STATE        = "previous_position"
	RETRIEVED_USERS_METADATA_SESSION_STATE = "retrieved_users_metadata"
)

// SessionStateReducer represents reducer used for session state management.
type SessionStateReducer struct {
	// Represents of instance of state store.
	store *godux.Store
}

func (ssr *SessionStateReducer) Init() {
	ssr.store.SetState(RESET_SESSION_STATE, value.RESET_SESSION_FALSE_VALUE)
	ssr.store.SetState(STATIC_SESSION_STATE, value.STATIC_SESSION_EMPTY_VALUE)
	ssr.store.SetState(POSITION_SESSION_STATE, value.POSITION_SESSION_EMPTY_VALUE)
	ssr.store.SetState(STAGE_POSITION_SESSION_STATE, value.STAGE_POSITION_SESSION_EMPTY_VALUE)
	ssr.store.SetState(PREVIOUS_POSITION_SESSION_STATE, value.POSITION_SESSION_EMPTY_VALUE)
	ssr.store.SetState(RETRIEVED_USERS_METADATA_SESSION_STATE, value.RETRIEVED_USERS_METADATA_EMPTY_VALUE)
}

func (ssr *SessionStateReducer) GetProcessor() func(value godux.Action) interface{} {
	return func(value godux.Action) interface{} {
		switch value.Type {
		case action.SET_RESET_SESSION_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: RESET_SESSION_STATE, Value: value.Value})

		case action.SET_STATIC_POSITION_SESSION_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: STATIC_SESSION_STATE, Value: value.Value})

		case action.SET_POSITION_SESSION_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: POSITION_SESSION_STATE, Value: value.Value},
				dto.ReducerResultUnit{Key: PREVIOUS_POSITION_SESSION_STATE, Value: value.Value},
				dto.ReducerResultUnit{Key: STAGE_POSITION_SESSION_STATE, Value: value.Value})

		case action.REVERT_STAGE_POSITION_X_SESSION_ACTION:
			positionValue := ssr.store.GetState(POSITION_SESSION_STATE).(dto.Position)

			stagePositionValue := ssr.store.GetState(STAGE_POSITION_SESSION_STATE).(dto.Position)

			stagePositionValue.X = positionValue.X

			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: STAGE_POSITION_SESSION_STATE, Value: stagePositionValue})

		case action.REVERT_STAGE_POSITION_Y_SESSION_ACTION:
			positionValue := ssr.store.GetState(POSITION_SESSION_STATE).(dto.Position)

			stagePositionValue := ssr.store.GetState(STAGE_POSITION_SESSION_STATE).(dto.Position)

			stagePositionValue.Y = positionValue.Y

			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: STAGE_POSITION_SESSION_STATE, Value: stagePositionValue})

		case action.SYNC_STAGE_POSITION_X_SESSION_ACTION:
			stagePositionValue := ssr.store.GetState(STAGE_POSITION_SESSION_STATE).(dto.Position)

			positionValue := ssr.store.GetState(POSITION_SESSION_STATE).(dto.Position)

			positionValue.X = stagePositionValue.X

			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: POSITION_SESSION_STATE, Value: positionValue})

		case action.SYNC_STAGE_POSITION_Y_SESSION_ACTION:
			stagePositionValue := ssr.store.GetState(STAGE_POSITION_SESSION_STATE).(dto.Position)

			positionValue := ssr.store.GetState(POSITION_SESSION_STATE).(dto.Position)

			positionValue.Y = stagePositionValue.Y

			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: POSITION_SESSION_STATE, Value: positionValue})

		case action.SYNC_PREVIOUS_POSITION_SESSION_ACTION:
			value := ssr.store.GetState(POSITION_SESSION_STATE).(dto.Position)

			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: PREVIOUS_POSITION_SESSION_STATE, Value: value})

		case action.INCREMENT_X_POSITION_SESSION_ACTION:
			valueRaw := ssr.store.GetState(POSITION_SESSION_STATE).(dto.Position)

			valueRaw.X += 1

			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: STAGE_POSITION_SESSION_STATE, Value: valueRaw})

		case action.INCREMENT_Y_POSITION_SESSION_ACTION:
			valueRaw := ssr.store.GetState(POSITION_SESSION_STATE).(dto.Position)

			valueRaw.Y += 1

			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: STAGE_POSITION_SESSION_STATE, Value: valueRaw})

		case action.DECREMENT_X_POSITION_SESSION_ACTION:
			valueRaw := ssr.store.GetState(POSITION_SESSION_STATE).(dto.Position)

			valueRaw.X -= 1

			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: STAGE_POSITION_SESSION_STATE, Value: valueRaw})

		case action.DECREMENT_Y_POSITION_SESSION_ACTION:
			valueRaw := ssr.store.GetState(POSITION_SESSION_STATE).(dto.Position)

			valueRaw.Y -= 1

			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: STAGE_POSITION_SESSION_STATE, Value: valueRaw})

		case action.DIAGONAL_UP_LEFT_POSITION_SESSION_ACTION:
			valueRaw := ssr.store.GetState(POSITION_SESSION_STATE).(dto.Position)

			valueRaw.X = valueRaw.X - 2.0/2.2360679775
			valueRaw.Y = valueRaw.Y + 1.0/2.2360679775

			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: STAGE_POSITION_SESSION_STATE, Value: valueRaw})

		case action.DIAGONAL_UP_RIGHT_POSITION_SESSION_ACTION:
			valueRaw := ssr.store.GetState(POSITION_SESSION_STATE).(dto.Position)

			valueRaw.X = valueRaw.X + 2.0/2.2360679775
			valueRaw.Y = valueRaw.Y + 1.0/2.2360679775

			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: STAGE_POSITION_SESSION_STATE, Value: valueRaw})

		case action.DIAGONAL_DOWN_LEFT_POSITION_SESSION_ACTION:
			valueRaw := ssr.store.GetState(POSITION_SESSION_STATE).(dto.Position)

			valueRaw.X = valueRaw.X - 2.0/2.2360679775
			valueRaw.Y = valueRaw.Y - 1.0/2.2360679775

			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: STAGE_POSITION_SESSION_STATE, Value: valueRaw})

		case action.DIAGONAL_DOWN_RIGHT_POSITION_SESSION_ACTION:
			valueRaw := ssr.store.GetState(POSITION_SESSION_STATE).(dto.Position)

			valueRaw.X = valueRaw.X + 2.0/2.2360679775
			valueRaw.Y = valueRaw.Y - 1.0/2.2360679775

			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: STAGE_POSITION_SESSION_STATE, Value: valueRaw})

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
