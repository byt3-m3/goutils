package vars

import "github.com/byt3-m3/goutils/env_utils"

var (
	LogLevel = env_utils.GetEnv("LOG_LEVEL", "info")
)
