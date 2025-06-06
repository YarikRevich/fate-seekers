package value

// Describes all the available screen reducer store values.
const (
	ACTIVE_SCREEN_LOGO_VALUE  = "logo"
	ACTIVE_SCREEN_INTRO_VALUE = "intro"

	// Entry screen is expected to be used for initialization operations.
	ACTIVE_SCREEN_ENTRY_VALUE = "entry"

	ACTIVE_SCREEN_MENU_VALUE         = "menu"
	ACTIVE_SCREEN_SETTINGS_VALUE     = "settings"
	ACTIVE_SCREEN_SELECTOR_VALUE     = "selector"
	ACTIVE_SCREEN_CREATOR_VALUE      = "creator"
	ACTIVE_SCREEN_LOBBY_VALUE        = "lobby"
	ACTIVE_SCREEN_SESSION_VALUE      = "session"
	ACTIVE_SCREEN_TRAVEL_VALUE       = "travel"
	ACTIVE_SCREEN_ANSWER_INPUT_VALUE = "answer_input"
	ACTIVE_SCREEN_RESUME_VALUE       = "resume"

	PREVIOUS_SCREEN_MENU_VALUE   = "menu"
	PREVIOUS_SCREEN_RESUME_VALUE = "resume"
	PREVIOUS_SCREEN_EMPTY_VALUE  = ""
)

// Describes all the available application reducer store values.
const (
	TRANSLATION_UPDATED_TRUE_VALUE  = "true"
	TRANSLATION_UPDATED_FALSE_VALUE = "false"

	EXIT_APPLICATION_TRUE_VALUE = "true"

	LOADING_APPLICATION_EMPTY_VALUE = 0
)

// Describes all the available repository reducer store values.
const (
	UUID_REPOSITORY_EMPTY_VALUE = ""

	UUID_CHECKED_REPOSITORY_TRUE_VALUE  = "true"
	UUID_CHECKED_REPOSITORY_FALSE_VALUE = "false"
)

// Describes all the available networking reducer store values.
const (
	ENTRY_HANDSHAKE_STARTED_NETWORKING_TRUE_VALUE  = "true"
	ENTRY_HANDSHAKE_STARTED_NETWORKING_FALSE_VALUE = "false"

	PING_CONNECTION_STARTED_NETWORKING_TRUE_VALUE  = "true"
	PING_CONNECTION_STARTED_NETWORKING_FALSE_VALUE = "false"

	SESSION_RETRIEVAL_STARTED_NETWORKING_TRUE_VALUE  = "true"
	SESSION_RETRIEVAL_STARTED_NETWORKING_FALSE_VALUE = "false"

	SESSION_CREATION_STARTED_NETWORKING_TRUE_VALUE  = "true"
	SESSION_CREATION_STARTED_NETWORKING_FALSE_VALUE = "false"

	SESSION_JOINING_STARTED_NETWORKING_TRUE_VALUE  = "true"
	SESSION_JOINING_STARTED_NETWORKING_FALSE_VALUE = "false"

	SESSION_REMOVAL_STARTED_NETWORKING_TRUE_VALUE  = "true"
	SESSION_REMOVAL_STARTED_NETWORKING_FALSE_VALUE = "false"
)

// Describes all the available letter reducer store values.
const (
	LETTER_UPDATED_TRUE_VALUE  = "true"
	LETTER_UPDATED_FALSE_VALUE = "false"
	LETTER_NAME_EMPTY_VALUE    = ""
	LETTER_IMAGE_EMPTY_VALUE   = ""
)

// Describes all the available answer input reducer store values.
const (
	ANSWER_INPUT_QUESTION_UPDATED_TRUE_VALUE  = "true"
	ANSWER_INPUT_QUESTION_UPDATED_FALSE_VALUE = "false"
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

// Describes available event reducer store values.
var (
	EVENT_NAME_EMPTY_VALUE      = ""
	EVENT_NAME_TOXIC_RAIN_VALUE = "toxic_rain"
	EVENT_STARTED_FALSE_VALUE   = "false"
	EVENT_STARTED_TRUE_VALUE    = "true"
	EVENT_ENDING_FALSE_VALUE    = "false"
	EVENT_ENDING_TRUE_VALUE     = "true"
)

// Describes available sound reducer store values.
var (
	SOUND_FX_UPDATED_FALSE_VALUE    = "false"
	SOUND_FX_UPDATED_TRUE_VALUE     = "true"
	SOUND_MUSIC_UPDATED_FALSE_VALUE = "false"
	SOUND_MUSIC_UPDATED_TRUE_VALUE  = "true"
)
