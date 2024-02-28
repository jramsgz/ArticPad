package infrastructure

import (
	"io"
	"log"
	"os"

	"github.com/rs/zerolog"
)

type LoggerConfig struct {
	Dir   string
	Level string
}

// startLogger starts the logger
// Returns the logger, the writer and the file
func startLogger(config *LoggerConfig) (zerolog.Logger, io.Writer, *os.File) {
	if _, err := os.Stat(config.Dir); os.IsNotExist(err) {
		os.MkdirAll(config.Dir, 0755)
	}
	file, err := os.OpenFile(config.Dir+"/articpad.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	mw := io.MultiWriter(os.Stdout, file)

	var logLevel zerolog.Level = zerolog.DebugLevel
	desiredLevel, err := zerolog.ParseLevel(config.Level)
	if err == nil {
		logLevel = desiredLevel
	}

	return zerolog.New(mw).With().Timestamp().Logger().Level(logLevel), mw, file
}
