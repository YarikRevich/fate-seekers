package death

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/luisvinicius167/godux"
)

// Describes all the available death reducer store states.
const (
	RESET_DEATH_STATE = "reset"
)

// DeathStateReducer represents reducer used for death state management.
type DeathStateReducer struct {
	// Represents of instance of state store.
	store *godux.Store
}

func (dsr *DeathStateReducer) Init() {
	dsr.store.SetState(RESET_DEATH_STATE, value.RESET_DEATH_TRUE_VALUE)
}

func (dsr *DeathStateReducer) GetProcessor() func(value godux.Action) interface{} {
	return func(value godux.Action) interface{} {
		switch value.Type {
		case action.SET_RESET_DEATH_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: RESET_DEATH_STATE, Value: value.Value})

		default:
			return nil
		}
	}
}

// NewDeathStateReducer initializes new instance of DeathStateReducer.
func NewDeathStateReducer(store *godux.Store) reducer.Reducer {
	return &DeathStateReducer{
		store: store,
	}
}
