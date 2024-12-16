package config

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	ErrReadingFromConfig = errors.New("err happened during config file read operation")
)

var (
	configFile      = flag.String("config", "config.yaml", "a name of configuration file")
	configDirectory = flag.String("configDirectory", getDefaultConfigDirectory(), "a directory where configuration file is located")

	debug bool

	databaseName                 string
	databaseConnectionRetryDelay time.Duration

	loggingLevel                  string
	loggingConsole                bool
	loggingName, loggingDirectory string
)

var (
	// Represents home directory where all application related data will be saved.
	internalDirectory = "/.fate-seekers-client"
)

// SetupDefaultConfig initializes default parameters for the configuration file.
func SetupDefaultConfig() {
	viper.SetDefault("operation.debug", false)
	viper.SetDefault("database.name", "fate_seekers.db")
	viper.SetDefault("database.connection-retry-delay", time.Second*3)
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.console", true)
	viper.SetDefault("logging.name", "fate-seekers.log")
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

	debug = viper.GetBool("operation.debug")
	databaseName = viper.GetString("database.name")
	databaseConnectionRetryDelay = viper.GetDuration("database.connection-retry-delay")
	loggingLevel = viper.GetString("logging.level")
	loggingConsole = viper.GetBool("logging.console")
	loggingName = viper.GetString("loggging.name")

	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	loggingDirectory = filepath.Join(homeDirectory, viper.GetString("logging.directory"))

	if err := os.MkdirAll(loggingDirectory, 0755); err != nil {
		log.Fatalln(err)
	}
}

func GetDebug() bool {
	return debug
}

func GetDatabaseName() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	return filepath.Join(homeDir, internalDirectory, databaseName)
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

func getDefaultConfigDirectory() string {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	return filepath.Join(homeDirectory, internalDirectory)
}
