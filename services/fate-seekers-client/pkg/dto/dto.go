package dto

import "time"

// Describes all the available letter attachment types.
const (
	ATTACHMENT_IMAGE_TYPE     = "image"
	ATTACHMENT_ANIMATION_TYPE = "animation"
	ATTACHMENT_AUDIO_TYPE     = "audio"
)

// GeneratedQuestionUnit represents a generated question unit.
type GeneratedQuestionUnit struct {
	// Represents generated question body.
	Question string

	// Represents answer for the generated question.
	Answer int
}

// LetterLoaderAttachmentUnit represents a letter attachment unit.
type LetterLoaderCollectionUnit struct {
	// Represents collection name.
	Name string `json:"name"`

	// Represents collection max value.
	Max int64 `json:"max"`

	// Represents collection index in the context of max collection value.
	Index int64 `json:"index"`
}

// LetterLoaderAttachmentUnit represents a letter attachment unit.
type LetterLoaderAttachmentUnit struct {
	// Represents attachment type. Currently supported are 'image', 'animation', 'audio'.
	Type string `json:"type"`

	// Represents attachment location.
	Location string `json:"location"`
}

// LetterLoaderUnit represents a letter unit used by loader to interprite raw file.
type LetterLoaderUnit struct {
	// Represents letter text.
	Text string `json:"text"`

	// Represents letter collection.
	Collection LetterLoaderCollectionUnit `json:"collection"`

	// Represents letter attachment.
	Attachment LetterLoaderAttachmentUnit `json:"attachment"`
}

// LetterUnit represents a letter unit.
type LetterUnit struct {
}

// SubtitlesUnit represents a subtitle unit.
type SubtitlesUnit struct {
	// Represents text for of the subtitles unit.
	Text string

	// Represents duration of which subtitles unit will be shown.
	Duration time.Duration
}

// NotificationUnit represents a notification unit.
type NotificationUnit struct {
	// Represents text for of the notification unit.
	Text string

	// Represents duration of which notification unit will be shown.
	Duration time.Duration
}

// ReducerResult represents result of reducer execution operation.
type ReducerResult map[string]interface{}

// ReducerResultUnit represents result unit of reducer execution operation.
type ReducerResultUnit struct {
	// Represents reducer result key.
	Key string

	// Represents reducer result value.
	Value interface{}
}

// ComposeReducerResult composes reducer result from the given reducer result units.
func ComposeReducerResult(units ...ReducerResultUnit) ReducerResult {
	result := make(map[string]interface{})

	for _, unit := range units {
		result[unit.Key] = unit.Value
	}

	return result
}
