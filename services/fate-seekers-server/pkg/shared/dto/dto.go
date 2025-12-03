package dto

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// Describes all the available letter attachment types.
const (
	ATTACHMENT_IMAGE_TYPE     = "image"
	ATTACHMENT_ANIMATION_TYPE = "animation"
	ATTACHMENT_AUDIO_TYPE     = "audio"
)

// Describes all the available event duration time.
const (
	EVENT_DURATION_TIME_TOXIC_RAIN = time.Second * 20
)

// Describes all the available event hit rates.
const (
	EVENT_HIT_RATE_TOXIC_RAIN = 2
)

// Describes all the available event frequency rates.
const (
	EVENT_FREQUENCY_RATE_TOXIC_RAIN = time.Second * 5
)

// Describes all the available event names.
const (
	EVENT_NAME_TOXIC_RAIN = "toxic_rain"
	EVENT_NAME_EMPTY      = ""
)

// Describes map, which contains all the available event names.
var EVENTS_NAME_MAP = []string{
	EVENT_NAME_TOXIC_RAIN,
	EVENT_NAME_EMPTY,
}

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
	// Represents text of the notification unit.
	Text string

	// Represents text color of the notification unit.
	Color color.Color

	// Represents duration of which notification unit will be shown.
	Duration time.Duration
}

// AmbientSoundUnit represents ambient sound unit.
type AmbientSoundUnit struct {
	// Represents ambient sound name.
	Name string

	// Represents ambient sound duration.
	Duration time.Duration

	// Represents ambient sound player.
	Player *audio.Player
}

// MusicSoundUnit represents music sound unit.
type MusicSoundUnit struct {
	// Represents music sound duration.
	Duration time.Duration

	// Represents music sound player.
	Player *audio.Player
}

// FXSoundUnit represents fx sound unit.
type FXSoundUnit struct {
	// Represents fx sound duration.
	Duration time.Duration

	// Represents fx sound player.
	Player *audio.Player
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

// SessionsRepositoryInsertOrUpdateRequest represents sessions repository entity update request.
type SessionsRepositoryInsertOrUpdateRequest struct {
	ID      int64
	Name    string
	Seed    int64
	Issuer  int64
	Started bool
}

// SessionsRepositoryGetByFiltersRequest represents sessions repository entity get by filters request.
type SessionsRepositoryGetByFiltersRequest struct {
	Name string
}

// GenerationsRepositoryInsertOrUpdateRequest represents generations repository entity update request.
type GenerationsRepositoryInsertOrUpdateRequest struct {
	ID        int64
	SessionID int64
	Name      string
	Type      string
	Active    bool
}

// AssociationsRepositoryInsertOrUpdateRequest represents associations repository entity update request.
type AssociationsRepositoryInsertOrUpdateRequest struct {
	ID           int64
	SessionID    int64
	GenerationID int64
	Name         string
}

// LobbiesRepositoryInsertOrUpdateRequest represents lobbies repository entity update request.
type LobbiesRepositoryInsertOrUpdateRequest struct {
	UserID         int64
	SessionID      int64
	Skin           uint64
	Health         uint64
	Active         bool
	Host           bool
	Eliminated     bool
	PositionX      float64
	PositionY      float64
	PositionStatic bool
}

// CacheSessionEntity represent cache session entity used by global networking cache.
type CacheSessionEntity struct {
	ID      int64
	Seed    int64
	Name    string
	Started bool
}

// CacheMetadataEntity represent cache metadata entity used by global networking cache.
type CacheMetadataEntity struct {
	LobbyID        int64
	SessionID      int64
	PositionX      float64
	PositionY      float64
	PositionStatic bool
	Skin           uint64
	Health         uint64
	Active         bool
	Eliminated     bool
	Host           bool
}

// CacheLobbySetEntity represent cache lobby set entity used by global networking cache.
type CacheLobbySetEntity struct {
	ID     int64
	Issuer string
	Skin   uint64
	Host   bool
}

// SessionEvent represents session event description.
type SessionEvent struct {
	EndRate       time.Time
	FrequencyRate time.Time
	PauseRate     time.Time
	// If name is empty, but other fields are not, it means that event
	// has ended and awaits for the pause to end to start another event shuffle.
	Name string
}

// CacheGenerationsEntity represent generations entity used by global networking cache.
type CacheGenerationsEntity struct {
	ID        int64
	SessionID int64
	Name      string
	Type      string
	Active    bool
}

// Position represents a set of coordinates.
type Position struct {
	X, Y int
}
