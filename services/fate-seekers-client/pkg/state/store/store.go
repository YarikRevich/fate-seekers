package store

import (
	"fmt"
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/answerinput"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/application"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/creator"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/event"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/letter"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/metadata"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/networking"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/prompt"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/repository"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/session"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/sound"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/statistics"
	"github.com/luisvinicius167/godux"
)

var (
	// GetInstance retrieves instance of the store, performing initilization if needed.
	GetInstance = sync.OnceValue[*godux.Store](newStore)
)

// GetActiveScreen retrieves active screen state value.
func GetActiveScreen() string {
	instance := GetInstance()

	return instance.GetState(screen.ACTIVE_SCREEN_STATE).(string)
}

// GetPreviousScreen retrieves previous screen state value.
func GetPreviousScreen() string {
	instance := GetInstance()

	return instance.GetState(screen.PREVIOUS_SCREEN_STATE).(string)
}

// GetApplicationExit retrieves exit application state value.
func GetApplicationExit() string {
	instance := GetInstance()

	return instance.GetState(application.EXIT_APPLICATION_STATE).(string)
}

// GetApplicationLoading retrieves loading application state value.
func GetApplicationLoading() int {
	instance := GetInstance()

	return instance.GetState(application.LOADING_APPLICATION_STATE).(int)
}

func GetApplicationStateReset() string {
	// GetApplicationStateReset retrieves state reset application state value.
	instance := GetInstance()

	return instance.GetState(application.STATE_RESET_APPLICATION_STATE).(string)
}

// GetRepositoryUUID retrieves uuid repository state value.
func GetRepositoryUUID() string {
	instance := GetInstance()

	return instance.GetState(repository.UUID_REPOSITORY_STATE).(string)
}

// GetRepositoryUUIDChecked retrieves uuid checked repository state value.
func GetRepositoryUUIDChecked() string {
	instance := GetInstance()

	return instance.GetState(repository.UUID_CHECKED_REPOSITORY_STATE).(string)
}

// GetEntryHandshakeStartedNetworking retrieves entry handshake started networking state value.
func GetEntryHandshakeStartedNetworking() string {
	instance := GetInstance()

	return instance.GetState(networking.ENTRY_HANDSHAKE_STARTED_NETWORKING_STATE).(string)
}

// GetPingConnectionStartedNetworking retrieves ping connection started networking state value.
func GetPingConnectionStartedNetworking() string {
	instance := GetInstance()

	return instance.GetState(networking.PING_CONNECTION_STARTED_NETWORKING_STATE).(string)
}

// GetSessionRetrievalStartedNetworking retrieves session retrieval started networking state value.
func GetSessionRetrievalStartedNetworking() string {
	instance := GetInstance()

	return instance.GetState(networking.SESSION_RETRIEVAL_STARTED_NETWORKING_STATE).(string)
}

// GetSessionCreationStartedNetworking retrieves session creation started networking state value.
func GetSessionCreationStartedNetworking() string {
	instance := GetInstance()

	return instance.GetState(networking.SESSION_CREATION_STARTED_NETWORKING_STATE).(string)
}

// GetSessionRemovalStartedNetworking retrieves session removal started networking state value.
func GetSessionRemovalStartedNetworking() string {
	instance := GetInstance()

	return instance.GetState(networking.SESSION_REMOVAL_STARTED_NETWORKING_STATE).(string)
}

// GetLobbySetRetrievalStartedNetworking retrieves lobby set retrieval started networking state value.
func GetLobbySetRetrievalStartedNetworking() string {
	instance := GetInstance()

	return instance.GetState(networking.LOBBY_SET_RETRIEVAL_STARTED_NETWORKING_STATE).(string)
}

// GetLobbySetRetrievalCycleFinishedNetworking retrieves lobby set retrieval cycle finished networking state value.
func GetLobbySetRetrievalCycleFinishedNetworking() string {
	instance := GetInstance()

	return instance.GetState(networking.LOBBY_SET_RETRIEVAL_CYCLE_FINISHED_NETWORKING_STATE).(string)
}

// GetLobbyCreationStartedNetworking retrieves lobby creation started networking state value.
func GetLobbyCreationStartedNetworking() string {
	instance := GetInstance()

	return instance.GetState(networking.LOBBY_CREATION_STARTED_NETWORKING_STATE).(string)
}

// GetLobbyRemovalStartedNetworking retrieves lobby removal started networking state value.
func GetLobbyRemovalStartedNetworking() string {
	instance := GetInstance()

	return instance.GetState(networking.LOBBY_REMOVAL_STARTED_NETWORKING_STATE).(string)
}

// GetSessionMetadataRetrievalStartedNetworking retrieves session metadata retrieval started networking state value.
func GetSessionMetadataRetrievalStartedNetworking() string {
	instance := GetInstance()

	return instance.GetState(networking.SESSION_METADATA_RETRIEVAL_STARTED_NETWORKING_STATE).(string)
}

// GetUpdateUserMetadataPositionsStartedNetworking retrieves session metadata retrieval started networking state value.
func GetUpdateUserMetadataPositionsStartedNetworking() string {
	instance := GetInstance()

	return instance.GetState(networking.UPDATE_USER_METADATA_POSITIONS_STARTED_NETWORKING_STATE).(string)
}

// GetEventRetrievalStartedNetworking retrieves event retrieval started networking state value.
func GetEventRetrievalStartedNetworking() string {
	instance := GetInstance()

	return instance.GetState(networking.EVENT_RETRIEVAL_STARTED_NETWORKING_STATE).(string)
}

// GetUsersMetadataRetrievalStartedNetworking retrieves users metadata retrieval started networking state value.
func GetUsersMetadataRetrievalStartedNetworking() string {
	instance := GetInstance()

	return instance.GetState(networking.USERS_METADATA_RETRIEVAL_STARTED_NETWORKING_STATE).(string)
}

// GetLetterUpdated retrieves letter updated state value.
func GetLetterUpdated() string {
	instance := GetInstance()

	return instance.GetState(letter.LETTER_UPDATED_LETTER_STATE).(string)
}

// GetLetterName retrieves letter name state value.
func GetLetterName() string {
	instance := GetInstance()

	return instance.GetState(letter.LETTER_NAME_LETTER_STATE).(string)
}

// GetLetterImage retrieves letter image state value.
func GetLetterImage() string {
	instance := GetInstance()

	return instance.GetState(letter.LETTER_IMAGE_LETTER_STATE).(string)
}

// GetAnswerInputQuestionUpdated retrieves question updated state value.
func GetAnswerInputQuestionUpdated() string {
	instance := GetInstance()

	return instance.GetState(answerinput.ANSWER_INPUT_QUESTION_UPDATED_STATE).(string)
}

// GetPromptUpdated retrieves prompt updated state value.
func GetPromptUpdated() string {
	instance := GetInstance()

	return instance.GetState(prompt.UPDATED_PROMPT_STATE).(string)
}

// GetPromptText retrieves prompt text state value.
func GetPromptText() string {
	instance := GetInstance()

	return instance.GetState(prompt.TEXT_PROMPT_STATE).(string)
}

// GetPromptSubmitCallback retrieves prompt submit callback state value.
func GetPromptSubmitCallback() func() {
	instance := GetInstance()

	return instance.GetState(prompt.SUBMIT_CALLBACK_PROMPT_STATE).(func())
}

// GetPromptCancelCallback retrieves prompt cancel callback state value.
func GetPromptCancelCallback() func() {
	instance := GetInstance()

	return instance.GetState(prompt.CANCEL_CALLBACK_PROMPT_STATE).(func())
}

// GetEventName retrieves event name state value.
func GetEventName() string {
	instance := GetInstance()

	return instance.GetState(event.NAME_EVENT_STATE).(string)
}

// GetEventStarted retrieves event started state value.
func GetEventStarted() string {
	instance := GetInstance()

	return instance.GetState(event.STARTED_EVENT_STATE).(string)
}

// GetEventEnding retrieves event ending state value.
func GetEventEnding() string {
	instance := GetInstance()

	return instance.GetState(event.ENDING_EVENT_STATE).(string)
}

// GetSoundMusicUpdated retrieves sound music updated state value.
func GetSoundMusicUpdated() string {
	instance := GetInstance()

	return instance.GetState(sound.MUSIC_UPDATED_SOUND_STATE).(string)
}

// GetStatisticsContentPing retrieves statistics content ping state value.
func GetStatisticsContentPing() int64 {
	instance := GetInstance()

	return instance.GetState(statistics.CONTENT_PING_STATISTICS_STATE).(int64)
}

// GetStatisticsMetadataPing retrieves statistics metadata ping state value.
func GetStatisticsMetadataPing() int64 {
	instance := GetInstance()

	return instance.GetState(statistics.METADATA_PING_STATISTICS_STATE).(int64)
}

// GetRetrievedSessionsMetadata retrieves retrieved sessions metadata state value.
func GetRetrievedSessionsMetadata() []dto.RetrievedSessionMetadata {
	instance := GetInstance()

	return instance.GetState(metadata.RETRIEVED_SESSIONS_METADATA_STATE).([]dto.RetrievedSessionMetadata)
}

// GetSelectedSessionMetadata retrieves selected session metadata state value.
func GetSelectedSessionMetadata() *dto.SelectedSessionMetadata {
	instance := GetInstance()

	return instance.GetState(metadata.SELECTED_SESSION_METADATA_STATE).(*dto.SelectedSessionMetadata)
}

// GetRetrievedLobbySetMetadata retrieves retrieved lobby set metadata state value.
func GetRetrievedLobbySetMetadata() []dto.RetrievedLobbySetMetadata {
	instance := GetInstance()

	return instance.GetState(metadata.RETRIEVED_LOBBY_SET_METADATA_STATE).([]dto.RetrievedLobbySetMetadata)
}

// GetSelectedLobbySetUnitMetadata retrieves selected lobby set unit metadata state value.
func GetSelectedLobbySetUnitMetadata() *dto.SelectedLobbySetUnitMetadata {
	instance := GetInstance()

	return instance.GetState(metadata.SELECTED_LOBBY_SET_UNIT_METADATA_STATE).(*dto.SelectedLobbySetUnitMetadata)
}

// GetSessionAlreadyStartedMetadata retrieves session already started metadata state value.
func GetSessionAlreadyStartedMetadata() string {
	instance := GetInstance()

	return instance.GetState(metadata.SESSION_ALREADY_STARTED_METADATA_STATE).(string)
}

// GetPositionSession retrieves position session state value.
func GetPositionSession() dto.Position {
	instance := GetInstance()

	return instance.GetState(session.POSITION_SESSION_STATE).(dto.Position)
}

// GetRetrievedUsersMetadataSession retrieves retrieved users metadata session state value.
func GetRetrievedUsersMetadataSession() dto.RetrievedUsersMetadataSessionSet {
	instance := GetInstance()

	return instance.GetState(session.RETRIEVED_USERS_METADATA_SESSION_STATE).(dto.RetrievedUsersMetadataSessionSet)
}

// newStore creates new instance of application store.
func newStore() *godux.Store {
	store := godux.NewStore()

	screenStateReducer := screen.NewScreenStateReducer(store)
	screenStateReducer.Init()

	applicationStateReducer := application.NewApplicationStateReducer(store)
	applicationStateReducer.Init()

	repositoryStateReducer := repository.NewRepositoryStateReducer(store)
	repositoryStateReducer.Init()

	networkingStateReducer := networking.NewNetworkingStateReducer(store)
	networkingStateReducer.Init()

	letterStateReducer := letter.NewLetterStateReducer(store)
	letterStateReducer.Init()

	answerInputReducer := answerinput.NewAnswerInputStateReducer(store)
	answerInputReducer.Init()

	promptReducer := prompt.NewPromptStateReducer(store)
	promptReducer.Init()

	eventReducer := event.NewEventStateReducer(store)
	eventReducer.Init()

	soundReducer := sound.NewEventStateReducer(store)
	soundReducer.Init()

	creatorReducer := creator.NewCreatorStateReducer(store)
	creatorReducer.Init()

	statisticsReducer := statistics.NewStatisticsStateReducer(store)
	statisticsReducer.Init()

	metadataReducer := metadata.NewMetadataStateReducer(store)
	metadataReducer.Init()

	sessionReducer := session.NewSessionStateReducer(store)
	sessionReducer.Init()

	store.Reducer(func(action godux.Action) interface{} {
		result := screenStateReducer.GetProcessor()(action)
		if result != nil {
			return result
		}

		result = applicationStateReducer.GetProcessor()(action)
		if result != nil {
			return result
		}

		result = repositoryStateReducer.GetProcessor()(action)
		if result != nil {
			return result
		}

		result = networkingStateReducer.GetProcessor()(action)
		if result != nil {
			return result
		}

		result = letterStateReducer.GetProcessor()(action)
		if result != nil {
			return result
		}

		result = answerInputReducer.GetProcessor()(action)
		if result != nil {
			return result
		}

		result = promptReducer.GetProcessor()(action)
		if result != nil {
			return result
		}

		result = eventReducer.GetProcessor()(action)
		if result != nil {
			return result
		}

		result = soundReducer.GetProcessor()(action)
		if result != nil {
			return result
		}

		result = creatorReducer.GetProcessor()(action)
		if result != nil {
			return result
		}

		result = statisticsReducer.GetProcessor()(action)
		if result != nil {
			return result
		}

		result = metadataReducer.GetProcessor()(action)
		if result != nil {
			return result
		}

		result = sessionReducer.GetProcessor()(action)
		if result != nil {
			return result
		}

		fmt.Println(action.Type)

		return nil
	})

	return store
}
