package logging

import (
	"log"
	"log/slog"
	"os"
	"strings"
)

const (
	LevelDebug   = "DEBUG"
	LevelInfo    = "INFO"
	LevelWarning = "WARNING"
	LevelError   = "ERROR"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Llongfile)
	switch strings.ToUpper(os.Getenv("LOG_LEVEL")) {
	case LevelDebug:
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case LevelInfo:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	case LevelWarning:
		slog.SetLogLoggerLevel(slog.LevelWarn)
	case LevelError:
		slog.SetLogLoggerLevel(slog.LevelError)
	default:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}
}
