package prompt

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/luisvinicius167/godux"
)

// Describes all the available networking reducer store states.
const (
	UPDATED_PROMPT_STATE         = "prompt_updated"
	TEXT_PROMPT_STATE            = "prompt_text"
	SUBMIT_CALLBACK_PROMPT_STATE = "prompt_submit_callback"
	CANCEL_CALLBACK_PROMPT_STATE = "prompt_cancel_callback"
)

// PromptStateReducer represents reducer used for prompt state management.
type PromptStateReducer struct {
	// Represents of instance of state store.
	store *godux.Store
}

func (psr *PromptStateReducer) Init() {
	psr.store.SetState(UPDATED_PROMPT_STATE, value.PROMPT_UPDATED_FALSE_VALUE)
	psr.store.SetState(TEXT_PROMPT_STATE, value.PROMPT_TEXT_EMPTY_VALUE)
	psr.store.SetState(SUBMIT_CALLBACK_PROMPT_STATE, value.PROMPT_SUBMIT_CALLBACK_EMPTY_VALUE)
	psr.store.SetState(CANCEL_CALLBACK_PROMPT_STATE, value.PROMPT_CANCEL_CALLBACK_EMPTY_VALUE)
}

func (psr *PromptStateReducer) GetProcessor() func(value godux.Action) interface{} {
	return func(value godux.Action) interface{} {
		switch value.Type {
		case action.SET_UPDATED_PROMPT_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{
					Key: UPDATED_PROMPT_STATE, Value: value.Value})

		case action.SET_TEXT_PROMPT_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{
					Key: TEXT_PROMPT_STATE, Value: value.Value})

		case action.SET_SUBMIT_CALLBACK_PROMPT_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{
					Key: SUBMIT_CALLBACK_PROMPT_STATE, Value: value.Value})

		case action.SET_CANCEL_CALLBACK_PROMPT_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{
					Key: CANCEL_CALLBACK_PROMPT_STATE, Value: value.Value})

		default:
			return nil
		}
	}
}

// NewPromptStateReducer initializes new instance of PromptStateReducer.
func NewPromptStateReducer(store *godux.Store) reducer.Reducer {
	return &PromptStateReducer{
		store: store,
	}
}
