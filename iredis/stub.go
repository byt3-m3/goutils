package iredis

import (
	"github.com/go-redis/redis"
	"log/slog"
	"time"
)

type StubRedisClient struct {
	SetStubReturn        func(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	GetStubReturn        func(key string) *redis.StringCmd
	GetClientStubReturn  func() *redis.Client
	MustValidateReturn   func()
	WithLoggerReturn     func(logger *slog.Logger) Client
	WithConnectionReturn func(host, password string, db int) Client
}

func (s *StubRedisClient) MustValidate() {
	s.MustValidateReturn()
}

func (s *StubRedisClient) WithLogger(logger *slog.Logger) Client {
	s.WithLoggerReturn(logger)
	return s
}

func (s *StubRedisClient) WithConnection(host, password string, db int) Client {
	s.WithConnectionReturn(host, password, db)
	return s
}

func (s *StubRedisClient) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return s.SetStubReturn(key, value, expiration)
}

type GetStubReturn struct {
	Cmd *redis.StringCmd
}

func (s *StubRedisClient) Get(key string) *redis.StringCmd {
	return s.GetStubReturn(key)
}

type GetClientStubReturn struct {
	Client *redis.Client
}

func (s *StubRedisClient) GetClient() *redis.Client {
	return s.GetClientStubReturn()
}
