package action

import "github.com/luisvinicius167/godux"

// Describes all the available state actions for screen reducer.
const (
	SET_ACTIVE_SCREEN_ACTION = "SET_ACTIVE_SCREEN"
)

// Describes all the available state actions for application reducer.
const (
	SET_TRANSLATION_UPDATED_APPLICATION_ACTION = "SET_TRANSLATION_UPDATED_APPLICATION_ACTION"
	SET_EXIT_APPLICATION_ACTION                = "SET_EXIT_APPLICATION_ACTION"
	SET_LOADING_APPLICATION_ACTION             = "SET_LOADING_APPLICATION_ACTION"
)

// Describes all the available state actions for networking reducer.
const (
	SET_ENTRY_HANDSHAKE_STARTED_NETWORKING_ACTION = "SET_ENTRY_HANDSHAKE_STARTED_NETWORKING_ACTION"
)

// Describes all the available state actions for letter reducer.
const (
	SET_LETTER_UPDATED_ACTION = "SET_LETTER_UPDATED_ACTION"
	SET_LETTER_NAME_ACTION    = "SET_LETTER_NAME_ACTION"
	SET_LETTER_IMAGE_ACTION   = "SET_LETTER_IMAGE_ACTION"
)

// Describes all the available state actions for answer input reducer.
const (
	SET_ANSWER_INPUT_SELECTED_CHEST_ACTION   = "SET_ANSWER_INPUT_SELECTED_CHEST_ACTION"
	SET_ANSWER_INPUT_QUESTION_UPDATED_ACTION = "SET_ANSWER_INPUT_QUESTION_UPDATED_ACTION"
)

// NewSetActiveScreenAction creates new set active screen action.
func NewSetActiveScreenAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_ACTIVE_SCREEN_ACTION,
		Value: value,
	}
}

// NewSetTranslationUpdatedApplicationAction creates new set translation updated application action.
func NewSetTranslationUpdatedApplicationAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_TRANSLATION_UPDATED_APPLICATION_ACTION,
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

// NewSetLetterUpdatedAction creates new set letter updated action.
func NewSetLetterUpdatedAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_LETTER_UPDATED_ACTION,
		Value: value,
	}
}

// NewSetLetterNameAction creates new set letter name action.
func NewSetLetterNameAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_LETTER_NAME_ACTION,
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

// NewSetAnswerInputSelectedChest creates new set answer input selected chest action.
func NewSetAnswerInputSelectedChest(value string) godux.Action {
	return godux.Action{
		Type:  SET_ANSWER_INPUT_SELECTED_CHEST_ACTION,
		Value: value,
	}
}

// NewSetAnswerInputQuestionUpdated creates new set answer input question updated action.
func NewSetAnswerInputQuestionUpdated(value string) godux.Action {
	return godux.Action{
		Type:  SET_ANSWER_INPUT_QUESTION_UPDATED_ACTION,
		Value: value,
	}
}
