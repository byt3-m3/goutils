package vars

import "github.com/byt3-m3/goutils/env_utils"

var (
	LogLevel      = env_utils.GetEnv("LOG_LEVEL", "info")
	KafkaUsername = env_utils.GetEnv("KAFKA_USERNAME", "user")
	KafkaPassword = env_utils.GetEnv("KAFKA_PASSWORD", "pass")
)
