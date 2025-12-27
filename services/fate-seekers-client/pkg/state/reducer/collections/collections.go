package collections

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/luisvinicius167/godux"
)

// Describes all the available collections reducer store states.
const (
	RESET_COLLECTIONS_STATE = "reset_collections"
)

// CollectionsStateReducer represents reducer used for collections state management.
type CollectionsStateReducer struct {
	// Represents of instance of state store.
	store *godux.Store
}

func (csr *CollectionsStateReducer) Init() {
	csr.store.SetState(RESET_COLLECTIONS_STATE, value.RESET_COLLECTIONS_FALSE_VALUE)
}

func (csr *CollectionsStateReducer) GetProcessor() func(value godux.Action) interface{} {
	return func(value godux.Action) interface{} {
		switch value.Type {
		case action.SET_RESET_COLLECTIONS_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: RESET_COLLECTIONS_STATE, Value: value.Value})

		default:
			return nil
		}
	}
}

// NewCollectionsStateReducer initializes new instance of CollectionsStateReducer.
func NewCollectionsStateReducer(store *godux.Store) reducer.Reducer {
	return &CollectionsStateReducer{
		store: store,
	}
}
