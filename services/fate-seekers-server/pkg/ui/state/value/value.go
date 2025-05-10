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
	TRANSLATION_UPDATED_TRUE_VALUE  = "true"
	TRANSLATION_UPDATED_FALSE_VALUE = "false"

	EXIT_APPLICATION_TRUE_VALUE = "true"

	LOADING_APPLICATION_EMPTY_VALUE = 0
)

// Describes all the available networking reducer store values.
const (
	ENTRY_HANDSHAKE_STARTED_NETWORKING_TRUE_VALUE  = "true"
	ENTRY_HANDSHAKE_STARTED_NETWORKING_FALSE_VALUE = "false"
)

// Describes available prompt reducer store values.
const (
	PROMPT_UPDATED_TRUE_VALUE  = "true"
	PROMPT_UPDATED_FALSE_VALUE = "false"
	PROMPT_TEXT_EMPTY_VALUE    = ""
)

// Describes available prompt reducer store values.
var (
	PROMPT_SUBMIT_CALLBACK_EMPTY_VALUE = func() {}
	PROMPT_CANCEL_CALLBACK_EMPTY_VALUE = func() {}
)

// Describes available info reducer store values.
const (
	INFO_UPDATED_TRUE_VALUE  = "true"
	INFO_UPDATED_FALSE_VALUE = "false"
	INFO_TEXT_EMPTY_VALUE    = ""
)

// Describes available info reducer store values.
var (
	INFO_CANCEL_CALLBACK_EMPTY_VALUE = func() {}
)
