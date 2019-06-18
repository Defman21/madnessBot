package common

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

var Log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC822}).With().Timestamp().Logger()

func SetLogLevel() {
	levels := map[string]zerolog.Level{
		"DEBUG": zerolog.DebugLevel,
		"INFO":  zerolog.InfoLevel,
		"WARN":  zerolog.WarnLevel,
	}

	logLevel := os.Getenv("LOG_LEVEL")

	if logLevel == "" {
		logLevel = "INFO"
	}

	Log = Log.Level(levels[logLevel])
}
