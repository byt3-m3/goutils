package client

import (
	"github.com/go-redis/redis"
	"time"
)

type StubRedisClient struct {
	SetStubReturn       func(key string, value interface{}, expiration time.Duration) SetStubReturn
	GetStubReturn       func(key string) GetStubReturn
	GetClientStubReturn func() GetClientStubReturn
}

type SetStubReturn struct {
	Cmd *redis.StatusCmd
}

func (s *StubRedisClient) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return s.SetStubReturn(key, value, expiration).Cmd
}

type GetStubReturn struct {
	Cmd *redis.StringCmd
}

func (s *StubRedisClient) Get(key string) *redis.StringCmd {
	return s.GetStubReturn(key).Cmd
}

type GetClientStubReturn struct {
	Client *redis.Client
}

func (s *StubRedisClient) GetClient() *redis.Client {
	return s.GetClientStubReturn().Client
}
