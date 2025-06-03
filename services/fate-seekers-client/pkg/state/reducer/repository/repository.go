package repository

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
	"github.com/luisvinicius167/godux"
)

// Describes all the available repository reducer store states.
const (
	UUID_CHECKED_REPOSITORY_STATE  = "uuid_checked"
	INTRO_CHECKED_REPOSITORY_STATE = "intro_checked"
)

// RepositoryStateReducer represents reducer used for repository state management.
type RepositoryStateReducer struct {
	// Represents of instance of state store.
	store *godux.Store
}

func (rsr *RepositoryStateReducer) Init() {
	rsr.store.SetState(UUID_CHECKED_REPOSITORY_STATE, value.UUID_CHECKED_REPOSITORY_FALSE_VALUE)
	rsr.store.SetState(INTRO_CHECKED_REPOSITORY_STATE, value.INTRO_CHECKED_REPOSITORY_FALSE_VALUE)
}

func (rsr *RepositoryStateReducer) GetProcessor() func(value godux.Action) interface{} {
	return func(value godux.Action) interface{} {
		switch value.Type {
		case action.SET_UUID_CHECKED_REPOSITORY_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: UUID_CHECKED_REPOSITORY_STATE, Value: value.Value})

		case action.SET_INTRO_CHECKED_REPOSITORY_ACTION:
			return dto.ComposeReducerResult(
				dto.ReducerResultUnit{Key: INTRO_CHECKED_REPOSITORY_STATE, Value: value.Value})

		default:
			return nil
		}
	}
}

// NewRepositoryStateReducer initializes new instance of RepositoryStateReducer.
func NewRepositoryStateReducer(store *godux.Store) reducer.Reducer {
	return &RepositoryStateReducer{
		store: store,
	}
}
