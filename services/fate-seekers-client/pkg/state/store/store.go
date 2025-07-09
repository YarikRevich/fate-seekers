package store

import (
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/answerinput"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/application"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/creator"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/event"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/letter"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/networking"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/prompt"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/repository"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/screen"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/sound"
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

// GetSessionJoiningStartedNetworking retrieves session joining started networking state value.
func GetSessionJoiningStartedNetworking() string {
	instance := GetInstance()

	return instance.GetState(networking.SESSION_JOINING_STARTED_NETWORKING_STATE).(string)
}

// GetSessionRemovalStartedNetworking retrieves session removal started networking state value.
func GetSessionRemovalStartedNetworking() string {
	instance := GetInstance()

	return instance.GetState(networking.SESSION_REMOVAL_STARTED_NETWORKING_STATE).(string)
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

// GetSoundFXUpdated retrieves sound fx updated state value.
func GetSoundFXUpdated() string {
	instance := GetInstance()

	return instance.GetState(sound.FX_UPDATED_SOUND_STATE).(string)
}

// GetSoundMusicUpdated retrieves sound music updated state value.
func GetSoundMusicUpdated() string {
	instance := GetInstance()

	return instance.GetState(sound.MUSIC_UPDATED_SOUND_STATE).(string)
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

		return nil
	})

	return store
}
