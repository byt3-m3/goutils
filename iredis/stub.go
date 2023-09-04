package iredis

import (
	"github.com/go-redis/redis"
	"time"
)

type StubRedisClient struct {
	SetStubReturn       func(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	GetStubReturn       func(key string) *redis.StringCmd
	GetClientStubReturn func() *redis.Client
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
