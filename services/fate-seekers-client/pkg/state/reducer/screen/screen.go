package screen

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/luisvinicius167/godux"
)

// Describes all the available screen reducer store states.
const (
	ACTIVE_SCREEN_STATE   = "active"
	PREVIOUS_SCREEN_STATE = "previous"
)

// ScreenStateReducer represents reducer used for screen state management.
type ScreenStateReducer struct {
	// Represents of instance of state store.
	store *godux.Store
}

func (ssr *ScreenStateReducer) Init() {
	ssr.store.SetState(ACTIVE_SCREEN_STATE, value.ACTIVE_SCREEN_INTRO_VALUE)
	ssr.store.SetState(PREVIOUS_SCREEN_STATE, value.PREVIOUS_SCREEN_EMPTY_VALUE)
}

func (ssr *ScreenStateReducer) GetProcessor() func(value godux.Action) interface{} {
	return func(value godux.Action) interface{} {
		switch value.Type {
		case action.SET_ACTIVE_SCREEN_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: ACTIVE_SCREEN_STATE, Value: value.Value})

		case action.SET_PREVIOUS_SCREEN_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: PREVIOUS_SCREEN_STATE, Value: value.Value})

		default:
			return nil
		}
	}
}

// NewScreenStateReducer initializes new instance of ScreenStateReducer.
func NewScreenStateReducer(store *godux.Store) reducer.Reducer {
	return &ScreenStateReducer{
		store: store,
	}
}
