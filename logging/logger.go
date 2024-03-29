package logging

import (
	"github.com/byt3-m3/goutils/vars"
	log "github.com/sirupsen/logrus"
	"log/slog"
	"os"
	"strings"
)

var logLevelMap = map[string]log.Level{
	"DEBUG": log.DebugLevel,
	"INFO":  log.InfoLevel,
}

func NewLogger() *log.Logger {
	return &log.Logger{
		Out:   os.Stdout,
		Level: logLevelMap[strings.ToUpper(vars.LogLevel)],

		Formatter: &log.JSONFormatter{
			TimestampFormat:   "",
			DisableTimestamp:  false,
			DisableHTMLEscape: false,
			DataKey:           "",
			FieldMap:          nil,
			CallerPrettyfier:  nil,
			PrettyPrint:       false,
		},
	}
}

// NewJSONLogger creates a Custom slog.Logger at the requested log level, if setDefault is True, then the custom logger
// will be set as the default logger.
func NewJSONLogger(level slog.Level, setDefault bool) *slog.Logger {

	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,
		Level:       level,
		ReplaceAttr: nil,
	}))

	if setDefault {
		slog.SetDefault(l)
	}

	return l

}
