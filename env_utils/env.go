package env_utils

import (
	"log"
	"os"
)

func GetEnvStrict(envVar string) string {
	val := os.Getenv(envVar)
	if len(val) == 0 {
		log.Fatalf("'%s' not found", envVar)
	}

	return val
}

func GetEnv(envVar, defaultVal string) string {
	val := os.Getenv(envVar)
	if len(val) == 0 {
		return defaultVal
	}

	return val
}
