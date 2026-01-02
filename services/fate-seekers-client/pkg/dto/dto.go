package dto

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/solarlune/resolv"
)

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
	// Reprsents letter title in the collections view.
	Title string `json:"title"`

	// Represents letter text.
	Text string `json:"text"`

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

	// Rperensets fx sound interruptable mode.
	Interruptable bool
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

// RetrievedSessionMetadata represents retrieved session holder for reducer components.
type RetrievedSessionMetadata struct {
	SessionID int64
	Name      string
	Seed      uint64
}

// RetrievedChests represents retrieved session holder for reducer components.
type RetrievedChests struct {
	SessionID  int64
	ID         int64
	Position   Position
	ChestItems []RetrievedChestItems
}

// RetrievedChestItems represents retrieved chest items.
type RetrievedChestItems struct {
	ID     int64
	Name   string
	Active bool
}

// RetrievedLobbySetMetadata represents retrieved lobby set holder for reducer components.
type RetrievedLobbySetMetadata struct {
	Issuer string
	Skin   uint64
	Host   bool
}

// GetFilteredSessionsRequest represents filtered sessions retrieval request.
type GetFilteredSessionsRequest struct {
	Name string
}

// SelectedSessionMetadata represents selected session metadata.
type SelectedSessionMetadata struct {
	ID   int64
	Name string
	Seed uint64
}

// SelectedLobbySetUnitMetadata represents selected lobby set unit metadata.
type SelectedLobbySetUnitMetadata struct {
	ID     int64
	Issuer string
	Skin   uint64
	Host   bool
}

// Position represents lobby user position.
type Position struct {
	X float64
	Y float64
}

const (
	SELECTED_MOVABLE_OBJECT = iota
	SELECTED_LOCAL_STATIC_OBJECT
	SELECTED_TILE_OBJECT
)

// SelectedObjectDetails represents selected object details.
type SelectedObjectDetails struct {
	Position Position

	// Represents a kind of a selectable tile used for selected worker.
	Kind int
}

// RetrievedInventoryUnit represents retrieved inventory unit.
type RetrievedInventoryUnit struct {
	ID   int64
	Name string
}

// RetrievedUsersMetadataSessionUnit represents retrieved users metadata session content unit.
type RetrievedUsersMetadataSessionUnit struct {
	Health             uint64
	Skin               uint64
	Active             bool
	Eliminated         bool
	AnimationDirection string
	AnimationStatic    bool
	Position           Position
	Change             time.Time
	Inventory          []RetrievedInventoryUnit
}

// RetrievedUsersMetadataSessionSet represents retrieved users metadata seession content set of units
type RetrievedUsersMetadataSessionSet map[string]RetrievedUsersMetadataSessionUnit

// RawMovableMetadata represents provided raw movable metadata
type RawMovableMetadata struct {
	Rotations  map[string]string   `json:"rotations"`
	Animations map[string][]string `json:"animations"`
}

// ProcessedMovableMetadata represents processed movable metadata unit.
type ProcessedMovableMetadataUnit struct {
	Rotation *ebiten.Image
	Frames   []*ebiten.Image
}

// ProcessedMovableMetadataSet represents movable metadata set.
type ProcessedMovableMetadataSet map[string]ProcessedMovableMetadataUnit

// Describes all the available moveable rotation directions
const (
	LeftMovableRotation      = "left"
	RightMovableRotation     = "right"
	UpMovableRotation        = "up"
	UpLeftMovableRotation    = "up-left"
	UpRightMovableRotation   = "up-right"
	DownMovableRotation      = "down"
	DownLeftMovableRotation  = "down-left"
	DownRightMovableRotation = "down-right"
)

// ProcessedTile represents processed tile.
type ProcessedTile struct {
	Position              Position
	TileWidth, TileHeight int
	Image                 *ebiten.Image
}

// CollidableTile represents collidable tile.
type CollidableTile struct {
	Position              Position
	TileWidth, TileHeight int
}

// CollidableStatic represents collidable static.
type CollidableStatic struct {
	Position              Position
	TileWidth, TileHeight int
}

// SelectableTile represents selectable tile.
type SelectableTile struct {
	Position Position

	TileWidth, TileHeight int
}

// SelectableStatic represents selectable static.
type SelectableStatic struct {
	Position Position

	TileWidth, TileHeight int
}

// SoundableTile represents soundable tile.
type SoundableTile struct {
	Position              Position
	Name                  string
	TileWidth, TileHeight int
}

// InventoryElement represents inventory element component.
type InventoryElement struct {
	// Represents inventory element image.
	Image *ebiten.Image

	// Represents apply interaction callback
	ApplyCallback func(success func())

	// Represents remove interaction callback
	RemoveCallback func(success func())
}

// ChestElement represents chest element component.
type ChestElement struct {
	// Represents inventory element image.
	Image *ebiten.Image

	// Represents chests click interaction callback.
	Callback func(success func())
}

// ExternalSounderObject represents external sounder object.
type ExternalSounderObject struct {
	// Represents if the position has been updated since the last update.
	Updated bool

	// Represents convex polygon value.
	Polygon *resolv.ConvexPolygon
}

// RendererPositionItem represents renderer position item.
type RendererPositionItem struct {
	Name string
	Type int
}

// Describes all the available animator movable object position item type.
const (
	RendererPositionItemMainCenteredMovable = iota
	RendererPositionItemSecondaryExternalMovable
	RendererPositionItemSecondaryTile
	RendererPositionItemSecondaryStatic
)

// Represents direction type used for gamepad stick direction processing.
type Direction int

// Represents available stick directions.
const (
	DirNone Direction = iota
	DirUp
	DirDown
	DirLeft
	DirRight
	DirUpLeft
	DirUpRight
	DirDownLeft
	DirDownRight
)

const (
	CHEST_ITEM_HEALTH_PACK_TYPE = "standard_health_pack"
	CHEST_ITEM_LETTER_TYPE      = "letter"
)
