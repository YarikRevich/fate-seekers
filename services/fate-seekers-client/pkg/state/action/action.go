package action

import "github.com/luisvinicius167/godux"

// Describes all the available state actions for screen reducer.
const (
	SET_ACTIVE_SCREEN_ACTION = "SET_ACTIVE_SCREEN"
)

// Describes all the available state actions for application reducer.
const (
	SET_EXIT_APPLICATION_ACTION    = "SET_EXIT_APPLICATION_ACTION"
	SET_LOADING_APPLICATION_ACTION = "SET_LOADING_APPLICATION_ACTION"
)

// Describes all the available state actions for networking reducer.
const (
	SET_ENTRY_HANDSHAKE_STARTED_NETWORKING_ACTION = "SET_ENTRY_HANDSHAKE_STARTED_NETWORKING_ACTION"
)

// Describes all the available state actions for letter reducer.
const (
	SET_LETTER_IMAGE_ACTION = "SET_LETTER_IMAGE_ACTION"
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

// NewSetLoadingApplicationAction creates new set loading application action.
func NewSetLoadingApplicationAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_LOADING_APPLICATION_ACTION,
		Value: value,
	}
}

// NewSetEntryHandshakeStartedNetworkingAction creates new set entry handshake started networking action.
func NewSetEntryHandshakeStartedNetworkingAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_ENTRY_HANDSHAKE_STARTED_NETWORKING_ACTION,
		Value: value,
	}
}

// NewSetLetterImageAction creates new set letter image action.
func NewSetLetterImageAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_LETTER_IMAGE_ACTION,
		Value: value,
	}
}
