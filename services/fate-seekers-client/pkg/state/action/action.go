package action

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
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
	SET_STATE_RESET_APPLICATION_ACTION   = "SET_STATE_RESET_APPLICATION_ACTION"
)

// Describes all the available state actions for repository reducer.
const (
	SET_UUID_REPOSITORY_ACTION         = "SET_UUID_REPOSITORY_ACTION"
	SET_UUID_CHECKED_REPOSITORY_ACTION = "SET_UUID_CHECKED_REPOSITORY_ACTION"
)

// Describes all the available state actions for networking reducer.
const (
	SET_ENTRY_HANDSHAKE_STARTED_NETWORKING_ACTION                = "SET_ENTRY_HANDSHAKE_STARTED_NETWORKING_ACTION"
	SET_PING_CONNECTION_STARTED_NETWORKING_ACTION                = "SET_PING_CONNECTION_STARTED_NETWORKING_ACTION"
	SET_SESSION_RETRIEVAL_STARTED_NETWORKING_ACTION              = "SET_SESSION_RETRIEVAL_STARTED_NETWORKING_ACTION"
	SET_SESSION_CREATION_STARTED_NETWORKING_ACTION               = "SET_SESSION_CREATION_STARTED_NETWORKING_ACTION"
	SET_SESSION_REMOVAL_STARTED_NETWORKING_ACTION                = "SET_SESSION_REMOVAL_STARTED_NETWORKING_ACTION"
	SET_LOBBY_SET_RETRIEVAL_STARTED_NETWORKING_ACTION            = "SET_LOBBY_SET_RETRIEVAL_STARTED_NETWORKING_ACTION"
	SET_LOBBY_SET_RETRIEVAL_CYCLE_FINISHED_NETWORKING_ACTION     = "SET_LOBBY_SET_RETRIEVAL_CYCLE_FINISHED_NETWORKING_ACTION"
	SET_LOBBY_CREATION_STARTED_NETWORKING_ACTION                 = "SET_LOBBY_CREATION_STARTED_NETWORKING_ACTION"
	SET_LOBBY_REMOVAL_STARTED_NETWORKING_ACTION                  = "SET_LOBBY_REMOVAL_STARTED_NETWORKING_ACTION"
	SET_SESSION_METADATA_RETRIEVAL_STARTED_NETWORKING_ACTION     = "SET_SESSION_METADATA_RETRIEVAL_STARTED_NETWORKING_ACTION"
	SET_UPDATE_USER_METADATA_POSITIONS_STARTED_NETWORKING_ACTION = "SET_UPDATE_USER_METADATA_POSITIONS_STARTED_NETWORKING_ACTION"
	SET_EVENT_RETRIEVAL_STARTED_NETWORKING_ACTION                = "SET_EVENT_RETRIEVAL_STARTED_NETWORKING_ACTION"
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
	SET_MUSIC_UPDATED_SOUND_ACTION = "SET_MUSIC_UPDATED_SOUND_ACTION"
)

// Describes all the available state actions for statistics reducer.
const (
	SET_CONTENT_PING_STATISTICS_ACTION  = "SET_CONTENT_PING_STATISTICS_ACTION"
	SET_METADATA_PING_STATISTICS_ACTION = "SET_METADATA_PING_STATISTICS_ACTION"
)

// Describes all the available state actions for metadata reducer.
const (
	SET_RETRIEVED_SESSIONS_METADATA_ACTION      = "SET_RETRIEVED_SESSIONS_METADATA_ACTION"
	SET_SELECTED_SESSION_METADATA_ACTION        = "SET_SELECTED_SESSION_METADATA_ACTION"
	SET_RETRIEVED_LOBBY_SET_METADATA_ACTION     = "SET_RETRIEVED_LOBBY_SET_METADATA_ACTION"
	SET_SELECTED_LOBBY_SET_UNIT_METADATA_ACTION = "SET_SELECTED_LOBBY_SET_UNIT_METADATA_ACTION"
	SET_SESSION_ALREADY_STARTED_METADATA_ACTION = "SET_SESSION_ALREADY_STARTED_METADATA_ACTION"
)

// Describes all the available state actions for session reducer.
const (
	INCREMENT_X_POSITION_SESSION_ACTION = "INCREMENT_X_POSITION_SESSION_ACTION"
	INCREMENT_Y_POSITION_SESSION_ACTION = "INCREMENT_Y_POSITION_SESSION_ACTION"
	DECREMENT_X_POSITION_SESSION_ACTION = "DECREMENT_X_POSITION_SESSION_ACTION"
	DECREMENT_Y_POSITION_SESSION_ACTION = "DECREMENT_Y_POSITION_SESSION_ACTION"
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

// NewSetUUIDRepositoryAction creates new set uuid repository action.
func NewSetUUIDRepositoryAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_UUID_REPOSITORY_ACTION,
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

// NewSetStateResetApplicationAction creates new set state reset application action.
func NewSetStateResetApplicationAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_STATE_RESET_APPLICATION_ACTION,
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

// NewSetPingConnectionStartedNetworkingAction creates new set ping connection started networking action.
func NewSetPingConnectionStartedNetworkingAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_PING_CONNECTION_STARTED_NETWORKING_ACTION,
		Value: value,
	}
}

// NewSetSessionRetrievalStartedNetworkingAction creates new set session retrieval started networking action.
func NewSetSessionRetrievalStartedNetworkingAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_SESSION_RETRIEVAL_STARTED_NETWORKING_ACTION,
		Value: value,
	}
}

// NewSetSessionCreationStartedNetworkingAction creates new set session creation started networking action.
func NewSetSessionCreationStartedNetworkingAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_SESSION_CREATION_STARTED_NETWORKING_ACTION,
		Value: value,
	}
}

// NewSetSessionRemovalStartedNetworkingAction creates new set session removal started networking action.
func NewSetSessionRemovalStartedNetworkingAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_SESSION_REMOVAL_STARTED_NETWORKING_ACTION,
		Value: value,
	}
}

// NewSetLobbySetRetrievalStartedNetworkingAction creates new lobby set retrieval started networking action.
func NewSetLobbySetRetrievalStartedNetworkingAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_LOBBY_SET_RETRIEVAL_STARTED_NETWORKING_ACTION,
		Value: value,
	}
}

// NewSetLobbySetRetrievalCycleFinishedNetworkingAction creates new lobby set retrieval cycle finished networking action.
func NewSetLobbySetRetrievalCycleFinishedNetworkingAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_LOBBY_SET_RETRIEVAL_CYCLE_FINISHED_NETWORKING_ACTION,
		Value: value,
	}
}

// NewSetLobbyCreationStartedNetworkingAction creates new set lobby creation started networking action.
func NewSetLobbyCreationStartedNetworkingAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_LOBBY_CREATION_STARTED_NETWORKING_ACTION,
		Value: value,
	}
}

// NewSetLobbyRemovalStartedNetworkingAction creates new set lobby removal started networking action.
func NewSetLobbyRemovalStartedNetworkingAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_LOBBY_REMOVAL_STARTED_NETWORKING_ACTION,
		Value: value,
	}
}

// NewSetSessionMetadataRetrievalStartedNetworkingAction creates new set session metadata retrieval started networking action.
func NewSetSessionMetadataRetrievalStartedNetworkingAction(value string) godux.Action {
	return godux.Action{
		Type:  SET_SESSION_METADATA_RETRIEVAL_STARTED_NETWORKING_ACTION,
		Value: value,
	}
}

// NewSetUpdateUserMetadataPositionsStartedNetworking creates new set update user metadata positions started networking action.
func NewSetUpdateUserMetadataPositionsStartedNetworking(value string) godux.Action {
	return godux.Action{
		Type:  SET_UPDATE_USER_METADATA_POSITIONS_STARTED_NETWORKING_ACTION,
		Value: value,
	}
}

// NewSetEventRetrievalStartedNetworking creates new set event retrieval started networking action.
func NewSetEventRetrievalStartedNetworking(value string) godux.Action {
	return godux.Action{
		Type:  SET_EVENT_RETRIEVAL_STARTED_NETWORKING_ACTION,
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

// NewSetSoundMusicUpdated creates new set sound music updated action.
func NewSetSoundMusicUpdated(value string) godux.Action {
	return godux.Action{
		Type:  SET_MUSIC_UPDATED_SOUND_ACTION,
		Value: value,
	}
}

// NewSetStatisticsContentPing creates new set statistics content ping action.
func NewSetStatisticsContentPing(value int64) godux.Action {
	return godux.Action{
		Type:  SET_CONTENT_PING_STATISTICS_ACTION,
		Value: value,
	}
}

// NewSetStatisticsMetadataPing creates new set statistics metadata ping action.
func NewSetStatisticsMetadataPing(value int64) godux.Action {
	return godux.Action{
		Type:  SET_METADATA_PING_STATISTICS_ACTION,
		Value: value,
	}
}

// NewSetRetrievedSessionsMetadata creates new set retrieved sessions metadata action.
func NewSetRetrievedSessionsMetadata(value []dto.RetrievedSessionMetadata) godux.Action {
	return godux.Action{
		Type:  SET_RETRIEVED_SESSIONS_METADATA_ACTION,
		Value: value,
	}
}

// NewSetSelectedSessionMetadata creates new set selected session metadata action.
func NewSetSelectedSessionMetadata(value *dto.SelectedSessionMetadata) godux.Action {
	return godux.Action{
		Type:  SET_SELECTED_SESSION_METADATA_ACTION,
		Value: value,
	}
}

// NewSetRetrievedLobbySetMetadata creates new set retrieved lobby set metadata action.
func NewSetRetrievedLobbySetMetadata(value []dto.RetrievedLobbySetMetadata) godux.Action {
	return godux.Action{
		Type:  SET_RETRIEVED_LOBBY_SET_METADATA_ACTION,
		Value: value,
	}
}

// NewSetSelectedLobbySetUnitMetadata creates new set selected lobby set unit metadata action.
func NewSetSelectedLobbySetUnitMetadata(value *dto.SelectedLobbySetUnitMetadata) godux.Action {
	return godux.Action{
		Type:  SET_SELECTED_LOBBY_SET_UNIT_METADATA_ACTION,
		Value: value,
	}
}

// NewSetSessionAlreadyStartedMetadata creates new set session already started metadata action.
func NewSetSessionAlreadyStartedMetadata(value string) godux.Action {
	return godux.Action{
		Type:  SET_SESSION_ALREADY_STARTED_METADATA_ACTION,
		Value: value,
	}
}

// NewIncrementXPositionSession creates new increment x position session action.
func NewIncrementXPositionSession() godux.Action {
	return godux.Action{
		Type: INCREMENT_X_POSITION_SESSION_ACTION,
	}
}

// NewIncrementYPositionSession creates new increment y position session action.
func NewIncrementYPositionSession() godux.Action {
	return godux.Action{
		Type: INCREMENT_Y_POSITION_SESSION_ACTION,
	}
}

// NewDecrementXPositionSession creates new decrement x position session action.
func NewDecrementXPositionSession() godux.Action {
	return godux.Action{
		Type: DECREMENT_X_POSITION_SESSION_ACTION,
	}
}

// NewDecrementYPositionSession creates new decrement y position session action.
func NewDecrementYPositionSession() godux.Action {
	return godux.Action{
		Type: DECREMENT_Y_POSITION_SESSION_ACTION,
	}
}
