package value

// Describes all the available screen reducer store values.
const (
	ACTIVE_SCREEN_ENTRY_VALUE      = "entry"
	ACTIVE_SCREEN_MENU_VALUE       = "menu"
	ACTIVE_SCREEN_MONITORING_VALUE = "monitoring"
	ACTIVE_SCREEN_SETTINGS_VALUE   = "settings"

	PREVIOUS_SCREEN_MENU_VALUE  = "menu"
	PREVIOUS_SCREEN_EMPTY_VALUE = ""
)

// Describes all the available application reducer store values.
const (
	EXIT_APPLICATION_TRUE_VALUE = "true"

	LOADING_APPLICATION_EMPTY_VALUE = 0
)

// Describes all the available networking reducer store values.
const (
	LISTENER_STARTED_NETWORKING_STATE_TRUE_VALUE  = "true"
	LISTENER_STARTED_NETWORKING_STATE_FALSE_VALUE = "false"
)

// Describes available prompt reducer store values.
const (
	UPDATED_PROMPT_TRUE_VALUE  = "true"
	UPDATED_PROMPT_FALSE_VALUE = "false"
	TEXT_PROMPT_EMPTY_VALUE    = ""
)

// Describes available prompt reducer store values.
var (
	SUBMIT_PROMPT_CALLBACK_EMPTY_VALUE = func() {}
	CANCEL_PROMPT_CALLBACK_EMPTY_VALUE = func() {}
)

// Describes available info reducer store values.
const (
	UPDATED_INFO_TRUE_VALUE  = "true"
	UPDATED_INFO_FALSE_VALUE = "false"
	TEXT_INFO_EMPTY_VALUE    = ""
)

// Describes available info reducer store values.
var (
	CANCEL_INFO_CALLBACK_EMPTY_VALUE = func() {}
)
