package logging

import (
	"github.com/byt3-m3/goutils/vars"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

var logLevelMap = map[string]log.Level{
	"DEBUG": log.DebugLevel,
	"INFO":  log.InfoLevel,
}

func NewLogger() log.Logger {
	return log.Logger{
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
