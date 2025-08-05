package info

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/reducer"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/value"
	"github.com/luisvinicius167/godux"
)

// Describes all the available networking reducer store states.
const (
	UPDATED_INFO_STATE         = "info_updated"
	TEXT_INFO_STATE            = "info_text"
	CANCEL_CALLBACK_INFO_STATE = "info_cancel_callback"
)

// InfoStateReducer represents reducer used for prompt state management.
type InfoStateReducer struct {
	// Represents of instance of state store.
	store *godux.Store
}

func (isr *InfoStateReducer) Init() {
	isr.store.SetState(UPDATED_INFO_STATE, value.UPDATED_INFO_FALSE_VALUE)
	isr.store.SetState(TEXT_INFO_STATE, value.TEXT_INFO_EMPTY_VALUE)
	isr.store.SetState(CANCEL_CALLBACK_INFO_STATE, value.CANCEL_INFO_CALLBACK_EMPTY_VALUE)
}

func (isr *InfoStateReducer) GetProcessor() func(value godux.Action) interface{} {
	return func(value godux.Action) interface{} {
		switch value.Type {
		case action.SET_UPDATED_INFO_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{
					Key: UPDATED_INFO_STATE, Value: value.Value})

		case action.SET_TEXT_INFO_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{
					Key: TEXT_INFO_STATE, Value: value.Value})

		case action.SET_CANCEL_CALLBACK_INFO_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{
					Key: CANCEL_CALLBACK_INFO_STATE, Value: value.Value})

		default:
			return nil
		}
	}
}

// NewInfoStateReducer initializes new instance of InfoStateReducer.
func NewInfoStateReducer(store *godux.Store) reducer.Reducer {
	return &InfoStateReducer{
		store: store,
	}
}
