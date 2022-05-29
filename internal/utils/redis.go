package utils

import (
	"fmt"
	"os"

	"github.com/gomodule/redigo/redis"
)

func NewRedisPool() (*redis.Pool, error) {
	addr, ok := os.LookupEnv("REDIS_ADDR")
	if !ok {
		return nil, fmt.Errorf("no redis address")
	}
	return &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		},
	}, nil
}
