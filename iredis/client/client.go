package client

import (
	"github.com/go-redis/redis"
	"log"
	"time"
)

type RedisClient interface {
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd

	Get(key string) *redis.StringCmd

	GetClient() *redis.Client
}

type redisClientOpt func(c *redisClient)

var (
	WithLogger = func(logger *log.Logger) redisClientOpt {
		return func(c *redisClient) {
			c.logger = logger
		}
	}
	WithRedisClient = func(host, password string, db int) redisClientOpt {
		return func(c *redisClient) {
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
		}
	}
)

type redisClient struct {
	client *redis.Client
	logger *log.Logger
}

func NewRedisClient(opts ...redisClientOpt) RedisClient {
	c := &redisClient{}

	for _, opt := range opts {
		opt(c)
	}

	if c.logger == nil {
		c.logger = log.Default()
	}

	if !validateClient(c) {
		c.logger.Fatalln("failed client validation")
	}

	return c
}

func validateClient(c *redisClient) bool {
	if c.client == nil {
		c.logger.Println("internal redis client not set, use WithRedisClient")
		return false
	}

	return true
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
