package event

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/luisvinicius167/godux"
)

// Describes all the available event reducer store states.
const (
	NAME_EVENT_STATE    = "name"
	STARTED_EVENT_STATE = "started"
	ENDING_EVENT_STATE  = "ending"
)

// EventStateReducer represents reducer used for event state management.
type EventStateReducer struct {
	// Represents of instance of state store.
	store *godux.Store
}

func (esr *EventStateReducer) Init() {
	esr.store.SetState(NAME_EVENT_STATE, value.EVENT_NAME_EMPTY_VALUE)
	esr.store.SetState(STARTED_EVENT_STATE, value.EVENT_STARTED_FALSE_VALUE)
	esr.store.SetState(ENDING_EVENT_STATE, value.EVENT_ENDING_FALSE_VALUE)
}

func (esr *EventStateReducer) GetProcessor() func(value godux.Action) interface{} {
	return func(value godux.Action) interface{} {
		switch value.Type {
		case action.SET_NAME_EVENT_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: NAME_EVENT_STATE, Value: value.Value})

		case action.SET_STARTED_EVENT_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: STARTED_EVENT_STATE, Value: value.Value})

		case action.SET_ENDING_EVENT_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: ENDING_EVENT_STATE, Value: value.Value})

		default:
			return nil
		}
	}
}

// NewEventStateReducer initializes new instance of EventStateReducer.
func NewEventStateReducer(store *godux.Store) reducer.Reducer {
	return &EventStateReducer{
		store: store,
	}
}
