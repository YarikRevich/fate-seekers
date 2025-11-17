package application

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/luisvinicius167/godux"
)

// Describes all the available application reducer store states.
const (
	EXIT_APPLICATION_STATE                     = "exit"
	LOADING_APPLICATION_STATE                  = "loading"
	STATE_RESET_APPLICATION_STATE              = "state_reset"
	GAMEPAD_ENABLED_APPLICATION_STATE          = "gamepad_enabled"
	GAMEPAD_POINTER_POSITION_APPLICATION_STATE = "gamepad_position"
)

// ApplicationStateReducer represents reducer used for application state management.
type ApplicationStateReducer struct {
	// Represents of instance of state store.
	store *godux.Store
}

func (asr *ApplicationStateReducer) Init() {
	asr.store.SetState(EXIT_APPLICATION_STATE, value.ACTIVE_SCREEN_MENU_VALUE)
	asr.store.SetState(LOADING_APPLICATION_STATE, value.LOADING_APPLICATION_EMPTY_VALUE)
	asr.store.SetState(STATE_RESET_APPLICATION_STATE, value.STATE_RESET_APPLICATION_FALSE_VALUE)
	asr.store.SetState(GAMEPAD_ENABLED_APPLICATION_STATE, value.GAMEPAD_ENABLED_APPLICATION_FALSE_VALUE)

	pointer := loader.GetInstance().GetStatic(loader.Pointer)

	shiftWidth := pointer.Bounds().Dx()
	shiftHeight := pointer.Bounds().Dy()

	asr.store.SetState(GAMEPAD_POINTER_POSITION_APPLICATION_STATE, dto.Position{
		X: (float64(config.GetWorldWidth()/2) - (float64(shiftWidth) / 2)),
		Y: (float64(config.GetWorldHeight()/2) - (float64(shiftHeight) / 2)),
	})
}

func (asr *ApplicationStateReducer) GetProcessor() func(value godux.Action) interface{} {
	return func(value godux.Action) interface{} {
		switch value.Type {
		case action.SET_EXIT_APPLICATION_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: EXIT_APPLICATION_STATE, Value: value.Value})

		case action.INCREMENT_LOADING_APPLICATION_ACTION:
			valueRaw := asr.store.GetState(LOADING_APPLICATION_STATE).(int)
			valueRaw += 1

			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: LOADING_APPLICATION_STATE, Value: valueRaw})

		case action.DECREMENT_LOADING_APPLICATION_ACTION:
			valueRaw := asr.store.GetState(LOADING_APPLICATION_STATE).(int)
			if valueRaw > 0 {
				valueRaw -= 1
			}

			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: LOADING_APPLICATION_STATE, Value: valueRaw})

		case action.SET_STATE_RESET_APPLICATION_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: STATE_RESET_APPLICATION_STATE, Value: value.Value})

		case action.SET_GAMEPAD_ENABLED_APPLICATION_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: GAMEPAD_ENABLED_APPLICATION_STATE, Value: value.Value})

		case action.INCREMENT_X_GAMEPAD_POINTER_POSITION_APPLICATION_ACTION:
			valueRaw := asr.store.GetState(GAMEPAD_POINTER_POSITION_APPLICATION_STATE).(dto.Position)

			valueRaw.X += 3

			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: GAMEPAD_POINTER_POSITION_APPLICATION_STATE, Value: valueRaw})

		case action.INCREMENT_Y_GAMEPAD_POINTER_POSITION_APPLICATION_ACTION:
			valueRaw := asr.store.GetState(GAMEPAD_POINTER_POSITION_APPLICATION_STATE).(dto.Position)

			valueRaw.Y += 3

			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: GAMEPAD_POINTER_POSITION_APPLICATION_STATE, Value: valueRaw})

		case action.DECREMENT_X_GAMEPAD_POINTER_POSITION_APPLICATION_ACTION:
			valueRaw := asr.store.GetState(GAMEPAD_POINTER_POSITION_APPLICATION_STATE).(dto.Position)

			valueRaw.X -= 3

			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: GAMEPAD_POINTER_POSITION_APPLICATION_STATE, Value: valueRaw})

		case action.DECREMENT_Y_GAMEPAD_POINTER_POSITION_APPLICATION_ACTION:
			valueRaw := asr.store.GetState(GAMEPAD_POINTER_POSITION_APPLICATION_STATE).(dto.Position)

			valueRaw.Y -= 3

			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: GAMEPAD_POINTER_POSITION_APPLICATION_STATE, Value: valueRaw})

		default:
			return nil
		}
	}
}

// NewApplicationStateReducer initializes new instance of ApplicationStateReducer.
func NewApplicationStateReducer(store *godux.Store) reducer.Reducer {
	return &ApplicationStateReducer{
		store: store,
	}
}
