package action

import (
	"github.com/luisvinicius167/godux"
)

// Describes all the available state actions for screen reducer.
const (
	SET_ACTIVE_SCREEN_ACTION   = "SET_ACTIVE_SCREEN"
	SET_PREVIOUS_SCREEN_ACTION = "SET_PREVIOUS_SCREEN"
)

// Describes all the available state actions for application reducer.
const (
	SET_EXIT_APPLICATION_ACTION          = "SET_EXIT_APPLICATION_ACTION"
	INCREMENT_LOADING_APPLICATION_ACTION = "INCREMENT_LOADING_APPLICATION_ACTION"
	DECREMENT_LOADING_APPLICATION_ACTION = "DECREMENT_LOADING_APPLICATION_ACTION"
)

// Describes all the available state actions for repository reducer.
const (
	SET_UUID_CHECKED_REPOSITORY_ACTION  = "SET_UUID_CHECKED_REPOSITORY_ACTION"
	SET_INTRO_CHECKED_REPOSITORY_ACTION = "SET_INTRO_CHECKED_REPOSITORY_ACTION"
)

// Describes all the available state actions for networking reducer.
const (
	SET_ENTRY_HANDSHAKE_STARTED_NETWORKING_ACTION = "SET_ENTRY_HANDSHAKE_STARTED_NETWORKING_ACTION"
	SET_PING_CONNECTION_STARTED_NETWORKING_ACTION = "SET_PING_CONNECTION_STARTED_NETWORKING_ACTION"
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

// Describes all the available state actions for prompt reducer.
const (
	SET_UPDATED_PROMPT_ACTION         = "SET_UPDATED_PROMPT_ACTION"
	SET_TEXT_PROMPT_ACTION            = "SET_TEXT_PROMPT_ACTION"
	SET_SUBMIT_CALLBACK_PROMPT_ACTION = "SET_SUBMIT_CALLBACK_PROMPT_ACTION"
	SET_CANCEL_CALLBACK_PROMPT_ACTION = "SET_CANCEL_CALLBACK_PROMPT_ACTION"
)

// Describes all the available state actions for event reducer.
const (
	SET_NAME_EVENT_ACTION    = "SET_NAME_EVENT_ACTION"
	SET_STARTED_EVENT_ACTION = "SET_STARTED_EVENT_ACTION"
	SET_ENDING_EVENT_ACTION  = "SET_ENDING_EVENT_ACTION"
)

// Describes all the available state actions for sound reducer.
const (
	SET_FX_UPDATED_SOUND_ACTION    = "SET_FX_UPDATED_SOUND_ACTION"
	SET_MUSIC_UPDATED_SOUND_ACTION = "SET_MUSIC_UPDATED_SOUND_ACTION"
)

// NewSetActiveScreenAction creates new set active screen action.
func NewSetActiveScreenAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_ACTIVE_SCREEN_ACTION,
		Value: value,
	}
}

// NewSetPreviousScreenAction creates new set previous screen action.
func NewSetPreviousScreenAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_PREVIOUS_SCREEN_ACTION,
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

// NewSetUUIDCheckedRepositoryAction creates new set uuid checked repository action.
func NewSetUUIDCheckedRepositoryAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_UUID_CHECKED_REPOSITORY_ACTION,
		Value: value,
	}
}

// NewSetIntroCheckedRepositoryAction creates new set intro checked repository action.
func NewSetIntroCheckedRepositoryAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_INTRO_CHECKED_REPOSITORY_ACTION,
		Value: value,
	}
}

// NewIncrementLoadingApplicationAction creates new increment loading application action.
func NewIncrementLoadingApplicationAction() godux.Action {
	return godux.Action{
		Type: INCREMENT_LOADING_APPLICATION_ACTION,
	}
}

// NewDecrementLoadingApplicationAction creates new decrement loading application action.
func NewDecrementLoadingApplicationAction() godux.Action {
	return godux.Action{
		Type: DECREMENT_LOADING_APPLICATION_ACTION,
	}
}

// NewSetEntryHandshakeStartedNetworkingAction creates new set entry handshake started networking action.
func NewSetEntryHandshakeStartedNetworkingAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_ENTRY_HANDSHAKE_STARTED_NETWORKING_ACTION,
		Value: value,
	}
}

// NewSetPingConnectionStartedNetworkingAction creates new set ping connection started networking action.
func NewSetPingConnectionStartedNetworkingAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_PING_CONNECTION_STARTED_NETWORKING_ACTION,
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

// NewSetPromptUpdated creates new set prompt updated action.
func NewSetPromptUpdated(value string) godux.Action {
	return godux.Action{
		Type:  SET_UPDATED_PROMPT_ACTION,
		Value: value,
	}
}

// NewSetPromptText creates new set prompt text action.
func NewSetPromptText(value string) godux.Action {
	return godux.Action{
		Type:  SET_TEXT_PROMPT_ACTION,
		Value: value,
	}
}

// NewSetPromptSubmitCallback creates new set prompt submit callback action.
func NewSetPromptSubmitCallback(value func()) godux.Action {
	return godux.Action{
		Type:  SET_SUBMIT_CALLBACK_PROMPT_ACTION,
		Value: value,
	}
}

// NewSetPromptCancelCallback creates new set prompt cancel callback action.
func NewSetPromptCancelCallback(value func()) godux.Action {
	return godux.Action{
		Type:  SET_CANCEL_CALLBACK_PROMPT_ACTION,
		Value: value,
	}
}

// NewSetEventName creates new set event name action.
func NewSetEventName(value string) godux.Action {
	return godux.Action{
		Type:  SET_NAME_EVENT_ACTION,
		Value: value,
	}
}

// NewSetEventStarted creates new set event started action.
func NewSetEventStarted(value string) godux.Action {
	return godux.Action{
		Type:  SET_STARTED_EVENT_ACTION,
		Value: value,
	}
}

// NewSetEventEnding creates new set event ending action.
func NewSetEventEnding(value string) godux.Action {
	return godux.Action{
		Type:  SET_ENDING_EVENT_ACTION,
		Value: value,
	}
}

// NewSetSoundFXUpdated creates new set sound fx updated action.
func NewSetSoundFXUpdated(value string) godux.Action {
	return godux.Action{
		Type:  SET_FX_UPDATED_SOUND_ACTION,
		Value: value,
	}
}

// NewSetSoundMusicUpdated creates new set sound music updated action.
func NewSetSoundMusicUpdated(value string) godux.Action {
	return godux.Action{
		Type:  SET_MUSIC_UPDATED_SOUND_ACTION,
		Value: value,
	}
}
