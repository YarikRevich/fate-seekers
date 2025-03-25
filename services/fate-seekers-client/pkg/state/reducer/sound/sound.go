package sound

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/luisvinicius167/godux"
)

// Describes all the available event reducer store states.
const (
	MUSIC_UPDATED_SOUND_STATE = "music_updated"
	FX_UPDATED_SOUND_STATE    = "fx_updated"
)

// EventStateReducer represents reducer used for event state management.
type EventStateReducer struct {
	// Represents of instance of state store.
	store *godux.Store
}

func (esr *EventStateReducer) Init() {
	esr.store.SetState(FX_UPDATED_SOUND_STATE, value.SOUND_FX_UPDATED_TRUE_VALUE)
	esr.store.SetState(MUSIC_UPDATED_SOUND_STATE, value.SOUND_MUSIC_UPDATED_TRUE_VALUE)
}

func (esr *EventStateReducer) GetProcessor() func(value godux.Action) interface{} {
	return func(value godux.Action) interface{} {
		switch value.Type {
		case action.SET_FX_UPDATED_SOUND_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: FX_UPDATED_SOUND_STATE, Value: value.Value})

		case action.SET_MUSIC_UPDATED_SOUND_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: MUSIC_UPDATED_SOUND_STATE, Value: value.Value})

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
