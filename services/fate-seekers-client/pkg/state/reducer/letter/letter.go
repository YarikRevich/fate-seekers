package letter

import (
	"fmt"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer"
	"github.com/luisvinicius167/godux"
)

// Describes all the available letter reducer store states.
const (
	LETTER_IMAGE_APPLICATION_STATE = "letter_image"
)

// LetterStateReducer represents reducer used for letter state management.
type LetterStateReducer struct {
	// Represents of instance of state store.
	store *godux.Store
}

func (lsr *LetterStateReducer) Init() {
	lsr.store.SetState(LETTER_IMAGE_APPLICATION_STATE, loader.Girls)
}

func (lsr *LetterStateReducer) GetProcessor() func(value godux.Action) interface{} {
	return func(value godux.Action) interface{} {
		switch value.Type {
		case action.SET_LETTER_IMAGE_ACTION:
			fmt.Println("IN REDUCER")

			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: LETTER_IMAGE_APPLICATION_STATE, Value: value.Value})

		default:
			return nil
		}
	}
}

// NewLetterStateReducer initializes new instance of LetterStateReducer.
func NewLetterStateReducer(store *godux.Store) reducer.Reducer {
	return &LetterStateReducer{
		store: store,
	}
}
