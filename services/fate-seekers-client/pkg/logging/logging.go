package logging

import (
	"log"
	"os"
	"path"
	"sync"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// GetInstance retrieves instance of the logger, performing initial creation if needed.
	GetInstance = sync.OnceValue[*zap.Logger](configure)
)

// setup performs logger configuration with the help of pre-defined configuration.
func configure() *zap.Logger {
	inputWriter := &lumberjack.Logger{
		Filename:   path.Join(config.GetLoggingDirectory(), config.GetLoggingName()),
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     28,
		LocalTime:  false,
		Compress:   false,
	}
	err := inputWriter.Rotate()
	if err != nil {
		log.Fatalln(err)
	}

	loggingWriter := zapcore.AddSync(inputWriter)

	var loggingConfig zap.Config

	if !config.GetOperationDebug() {
		loggingConfig = zap.NewProductionConfig()

		loggingConfig.DisableCaller = true
	} else {
		loggingConfig = zap.NewDevelopmentConfig()

		loggingConfig.EncoderConfig.LevelKey = "level"
		loggingConfig.EncoderConfig.NameKey = "name"
		loggingConfig.EncoderConfig.MessageKey = "msg"
		loggingConfig.EncoderConfig.CallerKey = "caller"
		loggingConfig.EncoderConfig.StacktraceKey = "stacktrace"

		if config.GetLoggingConsole() {
			loggingWriter = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stderr), loggingWriter)
		}
	}

	err = loggingConfig.Level.UnmarshalText([]byte(config.GetLoggingLevel()))
	if err != nil {
		log.Fatalln(err)
	}

	loggingConfig.EncoderConfig.TimeKey = "timestamp"
	loggingConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := loggingConfig.Build(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewCore(zapcore.NewJSONEncoder(loggingConfig.EncoderConfig), loggingWriter, loggingConfig.Level)
	}))
	if err != nil {
		log.Fatalln(err)
	}

	return logger
}
