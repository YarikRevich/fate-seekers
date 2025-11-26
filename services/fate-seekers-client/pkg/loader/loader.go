package loader

import (
	"bytes"
	"encoding/json"
	"image"
	"io/fs"
	"path/filepath"
	"sync"

	"github.com/YarikRevich/fate-seekers/assets"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/logging"
	"github.com/disintegration/imaging"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/lafriks/go-tiled"
	"github.com/pkg/errors"
	"github.com/tidwall/btree"
	"go.uber.org/zap"

	_ "image/jpeg"
	_ "image/png"

	"github.com/Frabjous-Studios/asebiten"
)

var (
	ErrReadingFile                       = errors.New("err happened during file read operation")
	ErrLoadingShader                     = errors.New("err happened during shader loading operation")
	ErrLoadingFont                       = errors.New("err happened during font loading operation")
	ErrLoadingStatic                     = errors.New("err happened during image loading operation")
	ErrLoadingAnimation                  = errors.New("err happened during animation loading operation")
	ErrLoadingMovable                    = errors.New("err happened during movable loading operation")
	ErrIncorrectMovableSkinIdentificator = errors.New("err happened incorrect movable skin identificator has been provided")
	ErrParsingMovable                    = errors.New("err happened during movable parsing operation")
)

var (
	// GetInstance retrieves instance of the asset loader manager, performing initial creation if needed.
	GetInstance = sync.OnceValue[*Loader](newLoader)
)

// Describes all the available maps to be loaded.
const (
	FirstMap = "1"
)

// Describes all the available map layers to be loaded.
const (
	FirstMapThirdLayer  = "third"
	FirstMapSecondLayer = "second"
)

// Describes tilemap configuration source file.
const (
	MapTilemap = "tilemap/tilemap.tmx"
)

// Describes available tilemap properties
const (
	TilemapCollidableProperty = "collidable"
	TilemapSoundProperty      = "sound"
	TilemapSpawnableProperty  = "spawnable"
)

// Describes all the available statics to be loaded.
const (
	ButtonIdleButton  = "ui/button-idle.png"
	ButtonHoverButton = "ui/button-hover.png"

	PanelIdlePanel = "ui/panel-idle.png"

	ListDisabled = "ui/list-disabled.png"
	ListIdle     = "ui/list-idle.png"

	ListMask          = "ui/list-mask.png"
	ListTrackDisabled = "ui/list-track-disabled.png"
	ListTrackIdle     = "ui/list-track-idle.png"

	SliderHandleHover = "ui/slider-handle-hover.png"
	SliderHandleIdle  = "ui/slider-handle-idle.png"
	SliderTrackIdle   = "ui/slider-track-idle.png"

	ComboArrayIdleButton = "ui/arrow-down-idle.png"
	ComboIdleButton      = "ui/combo-button-idle.png"

	TextInputIdle = "ui/text-input-idle.png"

	Heart = "heart/heart.png"

	DefaultLaserGun = "default_laser_gun/default_laser_gun.png"

	Pointer = "pointer/pointer.png"
)

// Describes all the available shaders to be loaded.
const (
	ToxicRainShader = "toxic-rain.kage"
)

// Describes all the available fonts to be loaded.
const (
	KyivRegularFont = "kyiv-regular.ttf"
)

// Describes all the available letters to be loaded.
const (
	LoneManLetter = "lone-man.json"
)

// Describes all the available client templates to be loaded.
const (
	EnglishClientTemplate   = "en/en_client.json"
	UkrainianClientTemplate = "uk/uk_client.json"
)

// Describes all the available shared templates to be loaded.
const (
	EnglishSharedTemplate   = "en/en_shared.json"
	UkrainianSharedTemplate = "uk/uk_shared.json"
)

// Describes all the available movables to be loaded.
const (
	Skins0Movable = "skins/0"
	Skins1Movable = "skins/0"
	Skins2Movable = "skins/0"
	Skins3Movable = "skins/0"
	Skins4Movable = "skins/0"
	Skins5Movable = "skins/0"
	Skins6Movable = "skins/0"
	Skins7Movable = "skins/0"
)

// Describes movable metadata file.
const (
	MovableMetadataFile = "metadata.json"
)

// Describes all the available animations to be loaded.
const (
	SkullAnimation  = "skull/skull.json"
	LogoAnimation   = "logo/logo.json"
	LoaderAnimation = "loader/loader.json"

	Background1Animation = "background/1/background-1.json"
	Background2Animation = "background/2/background-2.json"
	Background3Animation = "background/3/background-3.json"
	Background4Animation = "background/4/background-4.json"
	Background5Animation = "background/5/background-5.json"
	Background6Animation = "background/6/background-6.json"

	BlinkingScreen1Animation = "blinking-screen/1/blinking-screen-1.json"
	BlinkingScreen2Animation = "blinking-screen/2/blinking-screen-2.json"
	BlinkingScreen3Animation = "blinking-screen/3/blinking-screen-3.json"
	BlinkingScreen4Animation = "blinking-screen/4/blinking-screen-4.json"
)

// Describes all the available sounds to be loaded.
const (
	AmbientMusicSound   = "music/ambient/ambient.mp3"
	EnergetykMusicSound = "music/energetyk/energetyk.mp3"

	ButtonFXSound    = "fx/button/button.ogg"
	ToxicRainFXSound = "fx/toxicrain/toxicrain.ogg"
	RockFootFXSound  = "fx/rock_foot/rock_foot.ogg"
)

// Decsribes all the embedded files specific paths.
const (
	MapsPath       = "maps"
	ShadersPath    = "shaders"
	FontsPath      = "fonts"
	StaticsPath    = "statics"
	LettersPath    = "letters"
	TemplatesPath  = "templates"
	AnimationsPath = "animations"
	MovablePath    = "movable"
	SoundsPath     = "sounds"
)

// Loader represents low level asset loading manager, which operates in a lazy mode manner.
type Loader struct {
	// Represents cache map of embedded maps.
	maps sync.Map

	// Represents cache map of embedded statics.
	statics sync.Map

	// Represents cache map of embedded shaders.
	shaders sync.Map

	// Represents cache map of embedded fonts.
	fonts sync.Map

	// Represents cache map of embedded letters.
	letters sync.Map

	// Represents cache map of embedded templates.
	templates sync.Map

	// Represents cache map of embedded movables.
	movable sync.Map

	// Represents cache map of embedded animations.
	animations sync.Map

	// Represents cache map of embedded sounds.
	sounds sync.Map
}

// GetMap retrieves map content with the given name.
func (l *Loader) GetMap(name string) *tiled.Map {
	result, ok := l.maps.Load(name)
	if ok {
		return result.(*tiled.Map)
	}

	file, err := tiled.LoadFile(
		filepath.Join(common.ClientBasePath, MapsPath, name, MapTilemap),
		tiled.WithFileSystem(assets.AssetsClient))
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrReadingFile.Error()).Error())
	}

	l.maps.Store(name, file)

	logging.GetInstance().Debug("Map has been loaded", zap.String("name", name))

	return file
}

// GetMapLayerTiles retrieves map layer tiles for the provided layer.
func GetMapLayerTiles(layer *tiled.Layer, height, width, tileHeight, tileWidth int) (
	*btree.Map[float64, []*dto.ProcessedTile], []dto.Position, []*dto.CollidableTile, []*dto.SoundableTile) {
	var (
		result      = btree.NewMap[float64, []*dto.ProcessedTile](32)
		spawnables  []dto.Position
		collidables []*dto.CollidableTile
		soundables  []*dto.SoundableTile
	)

	var tiles sync.Map

	i := 0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if layer.Tiles[i].IsNil() {
				i++

				continue
			}

			tileImage, ok := tiles.Load(layer.Tiles[i].Tileset.FirstGID + layer.Tiles[i].ID)
			if !ok {
				for k := uint32(0); k < uint32(layer.Tiles[i].Tileset.TileCount); k++ {
					tiles.Store(k+layer.Tiles[i].Tileset.FirstGID, getMapTileImage(
						layer.Tiles[i].Tileset.GetFileFullPath(layer.Tiles[i].Tileset.Image.Source),
						layer.Tiles[i].Tileset.GetTileRect(uint32(k))))

					if layer.Tiles[i].ID == k {
						tileImage, _ = tiles.Load(layer.Tiles[i].Tileset.FirstGID + layer.Tiles[i].ID)
					}
				}
			}

			position := getMapTilePosition(x, y, tileWidth, tileHeight)

			processedTile := &dto.ProcessedTile{
				Position: position,
				Image:    tileImage.(*ebiten.Image),
			}

			for _, tile := range layer.Tiles[i].Tileset.Tiles {
				if layer.Tiles[i].Tileset.FirstGID+layer.Tiles[i].ID == tile.ID+layer.Tiles[i].Tileset.FirstGID {
					collidableProperty := tile.Properties.GetBool(TilemapCollidableProperty)
					if collidableProperty {
						// red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
						// processedTile.Image.Fill(red)

						collidables = append(collidables, &dto.CollidableTile{
							Position:   position,
							TileWidth:  tileWidth,
							TileHeight: tileHeight,
						})
					}

					soundableProperty := tile.Properties.GetString(TilemapSoundProperty)
					if soundableProperty != "" {
						soundables = append(soundables, &dto.SoundableTile{
							Position:   position,
							Name:       soundableProperty,
							TileWidth:  tileWidth,
							TileHeight: tileHeight,
						})
					}

					spawnableProperty := tile.Properties.GetBool(TilemapSpawnableProperty)
					if spawnableProperty {
						spawnables = append(spawnables, position)
					}
				}
			}

			var values []*dto.ProcessedTile

			values, ok = result.Get(position.Y)
			if ok {
				values = append(values, processedTile)
			} else {
				values = []*dto.ProcessedTile{processedTile}
			}

			result.Set(position.Y, values)

			i++
		}
	}

	return result, spawnables, collidables, soundables
}

// getMapTileImage reads cropped tile from the provided tilemap and the provided tile dimension.
func getMapTileImage(path string, rect image.Rectangle) *ebiten.Image {
	file, err := fs.ReadFile(assets.AssetsClient, path)
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrReadingFile.Error()).Error())
	}

	image, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrReadingFile.Error()).Error())
	}

	return ebiten.NewImageFromImage(imaging.Crop(image, rect))
}

// getMapTilePosition retrieves map tile position by the provided dimension.
func getMapTilePosition(x, y, tileWidth, tileHeight int) dto.Position {
	return dto.Position{
		X: float64((x - y) * (tileWidth / 2)),
		Y: float64((x + y) * (tileHeight / 2)),
	}
}

// GetStatic retrieves static content with the given name.
func (l *Loader) GetStatic(name string) *ebiten.Image {
	result, ok := l.statics.Load(name)
	if ok {
		return result.(*ebiten.Image)
	}

	file, err := common.ReadFile(filepath.Join(StaticsPath, name))
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrReadingFile.Error()).Error())
	}

	source, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrLoadingStatic.Error()).Error())
	}

	image := ebiten.NewImageFromImage(source)

	l.statics.Store(name, image)

	logging.GetInstance().Debug("Static has been loaded", zap.String("name", name))

	return image
}

// GetShader retrieves shader content with the given name.
func (l *Loader) GetShader(name string) *ebiten.Shader {
	result, ok := l.shaders.Load(name)
	if ok {
		return result.(*ebiten.Shader)
	}

	file, err := common.ReadFile(filepath.Join(ShadersPath, name))
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrReadingFile.Error()).Error())
	}

	shader, err := ebiten.NewShader(file)
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrLoadingShader.Error()).Error())
	}

	l.shaders.Store(name, shader)

	logging.GetInstance().Debug("Shader has been loaded", zap.String("name", name))

	return shader
}

// GetFont retrieves font content with the given name.
func (l *Loader) GetFont(name string) *text.GoTextFaceSource {
	result, ok := l.fonts.Load(name)
	if ok {
		return result.(*text.GoTextFaceSource)
	}

	file, err := common.ReadFile(filepath.Join(FontsPath, name))
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrReadingFile.Error()).Error())
	}

	font, err := text.NewGoTextFaceSource(bytes.NewReader(file))
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrLoadingFont.Error()).Error())
	}

	l.fonts.Store(name, font)

	logging.GetInstance().Debug("Font has been loaded", zap.String("name", name))

	return font
}

// GetLetter retrieves letter content with the given name.
func (l *Loader) GetLetter(name string) dto.LetterLoaderUnit {
	result, ok := l.letters.Load(name)
	if ok {
		return result.(dto.LetterLoaderUnit)
	}

	file, err := common.ReadFile(filepath.Join(LettersPath, name))
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrReadingFile.Error()).Error())
	}

	var data dto.LetterLoaderUnit

	err = json.Unmarshal(file, &data)
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrReadingFile.Error()).Error())
	}

	l.letters.Store(name, data)

	logging.GetInstance().Debug("Letter has been loaded", zap.String("name", name))

	return data
}

// GetTemplate retrieves template content with the given name.
func (l *Loader) GetTemplate(name string) []byte {
	result, ok := l.templates.Load(name)
	if ok {
		return result.([]byte)
	}

	file, err := common.ReadFile(filepath.Join(TemplatesPath, name))
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrReadingFile.Error()).Error())
	}

	l.templates.Store(name, file)

	logging.GetInstance().Debug("Template has been loaded", zap.String("name", name))

	return file
}

// GetSoundMusic retrieves music sound content with the given name.
func (l *Loader) GetSoundMusic(name string) *mp3.Stream {
	result, ok := l.sounds.Load(name)
	if ok {
		stream, err := mp3.DecodeF32(bytes.NewReader(result.([]byte)))
		if err != nil {
			logging.GetInstance().Fatal(errors.Wrap(err, ErrReadingFile.Error()).Error())
		}

		return stream
	}

	file, err := common.ReadFile(filepath.Join(SoundsPath, name))
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrReadingFile.Error()).Error())
	}

	l.sounds.Store(name, file)

	stream, err := mp3.DecodeF32(bytes.NewReader(file))
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrReadingFile.Error()).Error())
	}

	logging.GetInstance().Debug("Music sound has been loaded", zap.String("name", name))

	return stream
}

// GetSoundFX retrieves FX sound content with the given name.
func (l *Loader) GetSoundFX(name string) *vorbis.Stream {
	result, ok := l.sounds.Load(name)
	if ok {
		stream, err := vorbis.DecodeF32(bytes.NewReader(result.([]byte)))
		if err != nil {
			logging.GetInstance().Fatal(errors.Wrap(err, ErrReadingFile.Error()).Error())
		}

		return stream
	}

	file, err := common.ReadFile(filepath.Join(SoundsPath, name))
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrReadingFile.Error()).Error())
	}

	l.sounds.Store(name, file)

	stream, err := vorbis.DecodeF32(bytes.NewReader(file))
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrReadingFile.Error()).Error())
	}

	logging.GetInstance().Debug("FX sound has been loaded", zap.String("name", name))

	return stream
}

// GetMovableSkinsPath retrieves movable skins path according to the provided
// skin identificator.
func GetMovableSkinsPath(skin uint64) string {
	switch skin {
	case 0:
		return Skins0Movable
	case 1:
		return Skins1Movable
	case 2:
		return Skins2Movable
	case 3:
		return Skins3Movable
	case 4:
		return Skins4Movable
	case 5:
		return Skins5Movable
	case 6:
		return Skins6Movable
	case 7:
		return Skins7Movable
	default:
		logging.GetInstance().Fatal(ErrIncorrectMovableSkinIdentificator.Error())

		return ""
	}
}

// GetMovable retrieves movable content with the given name.
func (l *Loader) GetMovable(name string) dto.ProcessedMovableMetadataSet {
	result, ok := l.movable.Load(name)
	if ok {
		return result.(dto.ProcessedMovableMetadataSet)
	}

	file, err := common.ReadFile(filepath.Join(MovablePath, name, MovableMetadataFile))
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrReadingFile.Error()).Error())
	}

	var raw dto.RawMovableMetadata
	if err := json.Unmarshal(file, &raw); err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrParsingMovable.Error()).Error())
	}

	set := make(dto.ProcessedMovableMetadataSet)

	var (
		rotationFile   []byte
		rotationSource image.Image
		rotationImage  *ebiten.Image

		frameFile   []byte
		frameSource image.Image
		frameImage  *ebiten.Image
	)

	for direction, frames := range raw.Animations {
		rotationFile, err = common.ReadFile(filepath.Join(MovablePath, name, raw.Rotations[direction]))
		if err != nil {
			logging.GetInstance().Fatal(errors.Wrap(err, ErrReadingFile.Error()).Error())
		}

		rotationSource, _, err = image.Decode(bytes.NewReader(rotationFile))
		if err != nil {
			logging.GetInstance().Fatal(errors.Wrap(err, ErrLoadingMovable.Error()).Error())
		}

		rotationImage = ebiten.NewImageFromImage(rotationSource)

		var frameImages []*ebiten.Image

		for _, frame := range frames {
			frameFile, err = common.ReadFile(filepath.Join(MovablePath, name, frame))
			if err != nil {
				logging.GetInstance().Fatal(errors.Wrap(err, ErrReadingFile.Error()).Error())
			}

			frameSource, _, err = image.Decode(bytes.NewReader(frameFile))
			if err != nil {
				logging.GetInstance().Fatal(errors.Wrap(err, ErrLoadingMovable.Error()).Error())
			}

			frameImage = ebiten.NewImageFromImage(frameSource)

			frameImages = append(frameImages, frameImage)
		}

		set[direction] = dto.ProcessedMovableMetadataUnit{
			Rotation: rotationImage,
			Frames:   frameImages,
		}
	}

	l.movable.Store(name, set)

	logging.GetInstance().Debug("Movable has been loaded", zap.String("name", name))

	return set
}

// GetAnimation retrieves animation content with the given name. Allows to load new instance everytime.
// In order to load global object, it's required to set second argumeent to 'true'.
func (l *Loader) GetAnimation(name string, shared bool) *asebiten.Animation {
	if shared {
		result, ok := l.animations.Load(name)
		if ok {
			return result.(*asebiten.Animation)
		}
	}

	animation, err := common.LoadAnimation(filepath.Join(AnimationsPath, name))
	if err != nil {
		logging.GetInstance().Fatal(errors.Wrap(err, ErrLoadingAnimation.Error()).Error())
	}

	if shared {
		l.animations.Store(name, animation)

		logging.GetInstance().Debug("Animation has been loaded", zap.String("name", name))
	}

	return animation
}

// newLoader initializes Loader.
func newLoader() *Loader {
	return new(Loader)
}
