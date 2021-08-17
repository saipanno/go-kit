package client

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/saipanno/go-kit/logger"
)

// CreateRedisPool ...
func CreateRedisPool(conf *DBConfig) *redis.Pool {

	logger.Infof("create redis connect %s", conf.URI)

	return &redis.Pool{
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(conf.URI)
		},
	}
}
