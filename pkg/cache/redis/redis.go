package redis

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

type Config struct {
	Addr string `json:"addr"`
	// Maximum number of idle connections in the pool.
	Idle int `json:"max_idle"`

	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	Active int `json:"max_active"`

	// Close connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	IdleTimeout time.Duration `json:"idle_timeout"`

	// Close connections older than this duration. If the value is zero, then
	// the pool does not close connections based on age.
	MaxConnLifetime time.Duration `json:"max_conn_lifetime"`
}

func NewPool(c *Config) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     c.Idle,
		MaxActive:   c.Active,
		IdleTimeout: c.IdleTimeout * time.Second,
		Wait:        true,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) { return redis.Dial("tcp", c.Addr) },
	}
}
