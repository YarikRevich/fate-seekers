package creator

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer"
	"github.com/luisvinicius167/godux"
)

// Describes all the available creator reducer store states.
const ()

// CreatorStateReducer represents reducer used for networking state management.
type CreatorStateReducer struct {
	// Represents of instance of state store.
	store *godux.Store
}

func (csr *CreatorStateReducer) Init() {

}

func (csr *CreatorStateReducer) GetProcessor() func(value godux.Action) interface{} {
	return func(value godux.Action) interface{} {
		switch value.Type {
		default:
			return nil
		}
	}
}

// NewCreatorStateReducer initializes new instance of CreatorStateReducer.
func NewCreatorStateReducer(store *godux.Store) reducer.Reducer {
	return &CreatorStateReducer{
		store: store,
	}
}
