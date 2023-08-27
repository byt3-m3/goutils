package client

import (
	"fmt"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	rc := NewRedisClient(
		WithRedisClient("localhost:6379", "", 0),
	)

	result := rc.Set("test", "t", 20+time.Second)
	fmt.Println(result)
	data := rc.Get("test")
	fmt.Println(data)
}
