package travel

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/luisvinicius167/godux"
)

// Describes all the available travel reducer store states.
const (
	RESET_TRAVEL_STATE = "reset"
)

// TravelStateReducer represents reducer used for travel state management.
type TravelStateReducer struct {
	// Represents of instance of state store.
	store *godux.Store
}

func (tsr *TravelStateReducer) Init() {
	tsr.store.SetState(RESET_TRAVEL_STATE, value.RESET_TRAVEL_TRUE_VALUE)
}

func (tsr *TravelStateReducer) GetProcessor() func(value godux.Action) interface{} {
	return func(value godux.Action) interface{} {
		switch value.Type {
		case action.SET_RESET_TRAVEL_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: RESET_TRAVEL_STATE, Value: value.Value})

		default:
			return nil
		}
	}
}

// NewTravelStateReducer initializes new instance of TravelStateReducer.
func NewTravelStateReducer(store *godux.Store) reducer.Reducer {
	return &TravelStateReducer{
		store: store,
	}
}
