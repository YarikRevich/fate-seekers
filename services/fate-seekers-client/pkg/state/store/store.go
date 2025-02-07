package store

import (
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/application"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/letter"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/reducer/networking"
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

// GetTranslationUpdatedApplication retrieves translation updated application state value.
func GetTranslationUpdatedApplication() string {
	instance := GetInstance()

	return instance.GetState(application.TRANSLATION_UPDATED_APPLICATION_STATE).(string)
}

// GetExitApplication retrieves exit application state value.
func GetExitApplication() string {
	instance := GetInstance()

	return instance.GetState(application.EXIT_APPLICATION_STATE).(string)
}

// GetLoadingApplication retrieves loading application state value.
func GetLoadingApplication() string {
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

		return nil
	})

	return store
}
