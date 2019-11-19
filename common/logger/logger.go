package logger

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

var Log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC822}).With().Timestamp().Logger()

func SetLogLevel(logLevel string) {
	levels := map[string]zerolog.Level{
		"DEBUG": zerolog.DebugLevel,
		"INFO":  zerolog.InfoLevel,
		"WARN":  zerolog.WarnLevel,
	}

	if logLevel == "" {
		logLevel = "INFO"
	}

	Log = Log.Level(levels[logLevel])
}
