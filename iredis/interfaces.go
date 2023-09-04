package iredis

import (
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"time"
)

type (
	ClientOptionSetter interface {
		WithLogger(logger *log.Logger) Client

		WithConnection(host, password string, db int) Client
	}

	ClientValidator interface {
		MustValidate()
	}

	Client interface {
		ClientValidator
		ClientOptionSetter
		Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd

		Get(key string) *redis.StringCmd

		GetClient() *redis.Client
	}
)
