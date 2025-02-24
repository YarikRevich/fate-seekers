package config

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	ErrReadingFromConfig                 = errors.New("err happened during config file read operation")
	ErrReadingSettingsLanguageFromConfig = errors.New("err happened during config file settings language read operation")
)

var (
	configFile      = flag.String("config", "config.yaml", "a name of configuration file")
	configDirectory = flag.String("configDirectory", getDefaultConfigDirectory(), "a directory where configuration file is located")

	networkingHost string

	settingsLanguage string

	debug bool

	databaseName                 string
	databaseConnectionRetryDelay time.Duration

	loggingLevel                  string
	loggingConsole                bool
	loggingName, loggingDirectory string
)

// Represents window configuration.
const (
	windowName = "Fate Seekers"
)

// Represents internal world size.
const (
	worldWidth  = 640 * 2
	worldHeight = 360 * 2
)

// Represents internal min static size.
const (
	minStaticWidth  = 640
	minStaticHeight = 360
)

// Represents all the available settings language values.
const (
	SETTINGS_LANGUAGE_ENGLISH   = "en"
	SETTINGS_LANGUAGE_UKRAINIAN = "uk"
)

const (
	// Represents home directory where all application related data is located.
	internalGlobalDirectory = "/.fate-seekers-client"

	// Represents directory where all application configuration files are located.
	internalConfigDirectory = "/config"

	// Represents database directory where all the database files is located.
	internalDatabaseDirectory = "/internal/database"
)

// SetupDefaultConfig initializes default parameters for the configuration file.
func SetupDefaultConfig() {
	viper.SetDefault("settings.window.width", 1920)
	viper.SetDefault("settings.window.height", 1080)
	viper.SetDefault("settings.networking.host", "localhost:8080")
	viper.SetDefault("settings.language", SETTINGS_LANGUAGE_ENGLISH)
	viper.SetDefault("operation.debug", false)
	viper.SetDefault("database.name", "fate_seekers.db")
	viper.SetDefault("database.connection-retry-delay", time.Second*3)
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.console", true)
	viper.SetDefault("logging.name", "fate_seekers.log")
	viper.SetDefault("logging.directory", "log")
}

// Init initializes the configuration using provided configuration files and parameters.
func Init() {
	flag.Parse()

	viper.AddConfigPath(*configDirectory)

	viper.SetConfigFile(filepath.Join(*configDirectory, *configFile))
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(ErrReadingFromConfig.Error(), zap.String("configFile", *configFile), zap.Error(err))
	}

	windowWidth := viper.GetInt("settings.window.width")
	windowHeight := viper.GetInt("settings.window.height")
	networkingHost = viper.GetString("settings.networking.host")
	settingsLanguage = viper.GetString("settings.language")

	if settingsLanguage != SETTINGS_LANGUAGE_ENGLISH &&
		settingsLanguage != SETTINGS_LANGUAGE_UKRAINIAN {
		log.Fatalln(
			ErrReadingSettingsLanguageFromConfig.Error(),
			zap.String("configFile", *configFile),
			zap.String("settingsLanguage", settingsLanguage))
	}

	debug = viper.GetBool("operation.debug")
	databaseName = viper.GetString("database.name")
	databaseConnectionRetryDelay = viper.GetDuration("database.connection-retry-delay")
	loggingLevel = viper.GetString("logging.level")
	loggingConsole = viper.GetBool("logging.console")
	loggingName = viper.GetString("logging.name")

	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	loggingDirectory = filepath.Join(homeDirectory, internalGlobalDirectory, viper.GetString("logging.directory"))

	if err := os.MkdirAll(loggingDirectory, 0755); err != nil {
		log.Fatalln(err)
	}

	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle(windowName)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetVsyncEnabled(true)
}

func SetWindowSize(width, height int) {
	viper.Set("settings.window.width", width)
	viper.Set("settings.window.height", height)

	ebiten.SetWindowSize(width, height)
}

func SetNetworkingHost(value string) {
	viper.Set("settings.networking.host", value)
}

func GetNetworkingHost() string {
	return networkingHost
}

func SetSettingsLanguage(value string) {
	viper.Set("settings.language", value)
}

func GetSettingsLanguage() string {
	return settingsLanguage
}

func GetDebug() bool {
	return debug
}

func GetDatabaseName() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	return filepath.Join(homeDir, internalGlobalDirectory, internalDatabaseDirectory, databaseName)
}

func GetDatabaseConnectionRetryDelay() time.Duration {
	return databaseConnectionRetryDelay
}

func GetLoggingLevel() string {
	return loggingLevel
}

func GetLoggingConsole() bool {
	return loggingConsole
}

func GetLoggingName() string {
	return loggingName
}

func GetLoggingDirectory() string {
	return loggingDirectory
}

func GetWorldWidth() int {
	return worldWidth
}

func GetWorldHeight() int {
	return worldHeight
}

func GetMinStaticWidth() int {
	return minStaticWidth
}

func GetMinStaticHeight() int {
	return minStaticHeight
}

func getDefaultConfigDirectory() string {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	return filepath.Join(homeDirectory, internalGlobalDirectory, internalConfigDirectory)
}
