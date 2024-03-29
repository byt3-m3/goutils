package iredis

import (
	"github.com/go-redis/redis"
	"log/slog"

	"time"
)

type redisClient struct {
	client *redis.Client
	logger *slog.Logger
}

func New() Client {
	c := &redisClient{}

	return c
}

func (c *redisClient) MustValidate() {
	if c.logger == nil {
		c.logger = slog.Default()
	}

	if c.client == nil {
		panic("client not set, use WithConnection")

	}

}

func (c *redisClient) WithLogger(logger *slog.Logger) Client {
	c.logger = logger
	return c
}

func (c *redisClient) WithConnection(host, password string, db int) Client {
	rc := redis.NewClient(&redis.Options{
		Network:            "",
		Addr:               host,
		Dialer:             nil,
		OnConnect:          nil,
		Password:           password,
		DB:                 db,
		MaxRetries:         0,
		MinRetryBackoff:    0,
		MaxRetryBackoff:    0,
		DialTimeout:        0,
		ReadTimeout:        0,
		WriteTimeout:       0,
		PoolSize:           0,
		MinIdleConns:       0,
		MaxConnAge:         0,
		PoolTimeout:        0,
		IdleTimeout:        0,
		IdleCheckFrequency: 0,
		TLSConfig:          nil,
	})

	c.client = rc
	return c
}

func (c *redisClient) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return c.client.Set(key, value, expiration)

}

func (c *redisClient) GetClient() *redis.Client {
	return c.client
}

func (c *redisClient) Get(key string) *redis.StringCmd {
	return c.client.Get(key)
}
