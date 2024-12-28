package application

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/luisvinicius167/godux"
)

// Describes all the available screen reducer store states.
const (
	EXIT_APPLICATION_STATE = "exit"
)

// ApplicationStateReducer represents reducer used for application state management.
type ApplicationStateReducer struct {
	// Represents of instance of state store.
	store *godux.Store
}

func (asr *ApplicationStateReducer) Init() {
	asr.store.SetState(EXIT_APPLICATION_STATE, value.ACTIVE_SCREEN_MENU_VALUE)
}

func (asr *ApplicationStateReducer) GetProcessor() func(value godux.Action) interface{} {
	return func(value godux.Action) interface{} {
		switch value.Type {
		case action.SET_EXIT_APPLICATION_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: EXIT_APPLICATION_STATE, Value: value.Value})

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
