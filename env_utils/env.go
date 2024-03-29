package env_utils

import (
	"log/slog"
	"os"
)

func GetEnvStrict(envVar string) string {
	val := os.Getenv(envVar)
	if len(val) == 0 {
		slog.Error("env not found",
			slog.String("env", envVar),
		)
		panic("unable to find env")

	}

	return val
}

func GetEnv(envVar, defaultVal string) string {
	val := os.Getenv(envVar)
	if len(val) == 0 {
		slog.Debug("env not found, setting defaultVal",
			slog.String("env", envVar),
			slog.String("defaultVal", defaultVal),
		)
		return defaultVal
	}

	return val
}
