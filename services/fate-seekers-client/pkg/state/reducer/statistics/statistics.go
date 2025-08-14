package statistics

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/luisvinicius167/godux"
)

// Describes all the available statistics reducer store states.
const (
	CONTENT_PING_STATISTICS_STATE  = "content_ping"
	METADATA_PING_STATISTICS_STATE = "metadata_ping"
)

// StatisticsStateReducer represents reducer used for statistics state management.
type StatisticsStateReducer struct {
	// Represents of instance of state store.
	store *godux.Store
}

func (ssr *StatisticsStateReducer) Init() {
	ssr.store.SetState(CONTENT_PING_STATISTICS_STATE, value.CONTENT_PING_STATISTICS_EMPTY_VALUE)
	ssr.store.SetState(METADATA_PING_STATISTICS_STATE, value.METADATA_PING_STATISTICS_EMPTY_VALUE)
}

func (ssr *StatisticsStateReducer) GetProcessor() func(value godux.Action) interface{} {
	return func(value godux.Action) interface{} {
		switch value.Type {
		case action.SET_CONTENT_PING_STATISTICS_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: CONTENT_PING_STATISTICS_STATE, Value: value.Value})

		case action.SET_METADATA_PING_STATISTICS_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: METADATA_PING_STATISTICS_STATE, Value: value.Value})

		default:
			return nil
		}
	}
}

// NewStatisticsStateReducer initializes new instance of StatisticsStateReducer.
func NewStatisticsStateReducer(store *godux.Store) reducer.Reducer {
	return &StatisticsStateReducer{
		store: store,
	}
}
