package store

import (
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/reducer/application"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/reducer/info"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/reducer/networking"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/reducer/prompt"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/reducer/screen"
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

// GetListenerStartedNetworking retrieves listener started networking state value.
func GetListenerStartedNetworking() string {
	instance := GetInstance()

	return instance.GetState(networking.LISTENER_STARTED_NETWORKING_STATE).(string)
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

// GetInfoUpdated retrieves info updated state value.
func GetInfoUpdated() string {
	instance := GetInstance()

	return instance.GetState(info.UPDATED_INFO_STATE).(string)
}

// GetInfoText retrieves info text state value.
func GetInfoText() string {
	instance := GetInstance()

	return instance.GetState(info.TEXT_INFO_STATE).(string)
}

// GetInfoCancelCallback retrieves info cancel callback state value.
func GetInfoCancelCallback() func() {
	instance := GetInstance()

	return instance.GetState(info.CANCEL_CALLBACK_INFO_STATE).(func())
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

	promptReducer := prompt.NewPromptStateReducer(store)
	promptReducer.Init()

	infoReducer := info.NewInfoStateReducer(store)
	infoReducer.Init()

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

		result = promptReducer.GetProcessor()(action)
		if result != nil {
			return result
		}

		result = infoReducer.GetProcessor()(action)
		if result != nil {
			return result
		}

		return nil
	})

	return store
}
