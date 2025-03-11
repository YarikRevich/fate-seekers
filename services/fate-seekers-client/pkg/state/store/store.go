package store

import (
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/answerinput"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/application"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/event"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/letter"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/networking"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/prompt"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/screen"
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

// GetApplicationExit retrieves exit application state value.
func GetApplicationExit() string {
	instance := GetInstance()

	return instance.GetState(application.EXIT_APPLICATION_STATE).(string)
}

// GetApplicationLoading retrieves loading application state value.
func GetApplicationLoading() string {
	instance := GetInstance()

	return instance.GetState(application.LOADING_APPLICATION_STATE).(string)
}

// GetEntryHandshakeStartedNetworking retrieves entry handshake started networking state value.
func GetEntryHandshakeStartedNetworking() string {
	instance := GetInstance()

	return instance.GetState(networking.ENTRY_HANDSHAKE_STARTED_NETWORKING_STATE).(string)
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

// newStore creates new instance of application store.
func newStore() *godux.Store {
	store := godux.NewStore()

	screenStateReducer := screen.NewScreenStateReducer(store)
	screenStateReducer.Init()

	applicationStateReducer := application.NewApplicationStateReducer(store)
	applicationStateReducer.Init()

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

	store.Reducer(func(action godux.Action) interface{} {
		result := screenStateReducer.GetProcessor()(action)
		if result != nil {
			return result
		}

		result = applicationStateReducer.GetProcessor()(action)
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

		return nil
	})

	return store
}
