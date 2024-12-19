package action

import "github.com/luisvinicius167/godux"

// Describes all the available state actions for screen reducer.
const (
	SET_ACTIVE_SCREEN_ACTION = "SET_ACTIVE_SCREEN"
)

// NewSetActiveScreenAction creates new set active screen action.
func NewSetActiveScreenAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_ACTIVE_SCREEN_ACTION,
		Value: value,
	}
}
