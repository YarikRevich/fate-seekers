package value

// Describes all the available screen reducer store values.
const (
	ACTIVE_SCREEN_INTRO_VALUE        = "intro"
	ACTIVE_SCREEN_ENTRY_VALUE        = "entry"
	ACTIVE_SCREEN_MENU_VALUE         = "menu"
	ACTIVE_SCREEN_SETTINGS_VALUE     = "settings"
	ACTIVE_SCREEN_SESSION_VALUE      = "session"
	ACTIVE_SCREEN_TRAVEL_VALUE       = "travel"
	ACTIVE_SCREEN_ANSWER_INPUT_VALUE = "answer_input"
)

// Describes all the available application reducer store values.
const (
	TRANSLATION_UPDATED_TRUE_VALUE  = "true"
	TRANSLATION_UPDATED_FALSE_VALUE = "false"

	EXIT_APPLICATION_TRUE_VALUE = "true"

	LOADING_APPLICATION_TRUE_VALUE  = "true"
	LOADING_APPLICATION_FALSE_VALUE = "false"
)

// Describes all the available networking reducer store values.
const (
	ENTRY_HANDSHAKE_STARTED_NETWORKING_TRUE_VALUE  = "true"
	ENTRY_HANDSHAKE_STARTED_NETWORKING_FALSE_VALUE = "false"
)

// Describes all the available letter reducer store values.
const (
	LETTER_UPDATED_TRUE_VALUE  = "true"
	LETTER_UPDATED_FALSE_VALUE = "false"
	LETTER_NAME_EMPTY_VALUE    = ""
	LETTER_IMAGE_EMPTY_VALUE   = ""
)
