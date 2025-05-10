package application

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/reducer"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/value"
	"github.com/luisvinicius167/godux"
)

// Describes all the available application reducer store states.
const (
	EXIT_APPLICATION_STATE    = "exit"
	LOADING_APPLICATION_STATE = "loading"
)

// ApplicationStateReducer represents reducer used for application state management.
type ApplicationStateReducer struct {
	// Represents of instance of state store.
	store *godux.Store
}

func (asr *ApplicationStateReducer) Init() {
	asr.store.SetState(EXIT_APPLICATION_STATE, value.ACTIVE_SCREEN_MENU_VALUE)
	asr.store.SetState(LOADING_APPLICATION_STATE, value.LOADING_APPLICATION_EMPTY_VALUE)
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
