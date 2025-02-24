package answerinput

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/luisvinicius167/godux"
)

// Describes all the available answer input reducer store states.
const (
	ANSWER_INPUT_QUESTION_UPDATED_STATE = "answer_input_question_updated"
)

// AnswerInputStateReducer represents reducer used for answer input state management.
type AnswerInputStateReducer struct {
	// Represents of instance of state store.
	store *godux.Store
}

func (aisr *AnswerInputStateReducer) Init() {
	aisr.store.SetState(ANSWER_INPUT_QUESTION_UPDATED_STATE, value.ANSWER_INPUT_QUESTION_UPDATED_FALSE_VALUE)
}

func (aisr *AnswerInputStateReducer) GetProcessor() func(value godux.Action) interface{} {
	return func(value godux.Action) interface{} {
		switch value.Type {
		case action.SET_ANSWER_INPUT_QUESTION_UPDATED_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: ANSWER_INPUT_QUESTION_UPDATED_STATE, Value: value.Value})

		default:
			return nil
		}
	}
}

// NewAnswerInputStateReducer initializes new instance of AnswerInputStateReducer.
func NewAnswerInputStateReducer(store *godux.Store) reducer.Reducer {
	return &AnswerInputStateReducer{
		store: store,
	}
}
