package action

import "github.com/luisvinicius167/godux"

// Describes all the available state actions for screen reducer.
const (
	SET_ACTIVE_SCREEN_ACTION = "SET_ACTIVE_SCREEN"
)

// Describes all the available state actions for application reducer.
const (
	SET_EXIT_APPLICATION_ACTION = "SET_EXIT_APPLICATION_ACTION"
)

// NewSetActiveScreenAction creates new set active screen action.
func NewSetActiveScreenAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_ACTIVE_SCREEN_ACTION,
		Value: value,
	}
}

// NewSetExitApplicationAction creates new set exit application action.
func NewSetExitApplicationAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_EXIT_APPLICATION_ACTION,
		Value: value,
	}
}
