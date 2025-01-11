package value

// Describes all the available screen reducer store values.
const (
	ACTIVE_SCREEN_INTRO_VALUE    = "intro"
	ACTIVE_SCREEN_ENTRY_VALUE    = "entry"
	ACTIVE_SCREEN_MENU_VALUE     = "menu"
	ACTIVE_SCREEN_SETTINGS_VALUE = "settings"
)

// Describes all the available application reducer store values.
const (
	EXIT_APPLICATION_TRUE_VALUE = "true"

	LOADING_APPLICATION_TRUE_VALUE  = "true"
	LOADING_APPLICATION_FALSE_VALUE = "false"
)

// Describes all the available networking reducer store values.
const (
	ENTRY_HANDSHAKE_STARTED_NETWORKING_TRUE_VALUE  = "true"
	ENTRY_HANDSHAKE_STARTED_NETWORKING_FALSE_VALUE = "false"
)
