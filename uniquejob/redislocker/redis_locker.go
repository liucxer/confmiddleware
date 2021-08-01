package redislocker

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

type RedisConn interface {
	Get() redis.Conn
	Prefix(key string) string
}

func NewRedisLocker(redisConn RedisConn) *RedisLocker {
	return &RedisLocker{
		redisConn: redisConn,
	}
}

type RedisLocker struct {
	redisConn RedisConn
}

func (locker *RedisLocker) Lock(key string, expiresIn time.Duration) (bool, error) {
	conn := locker.redisConn.Get()
	defer conn.Close()

	k := locker.redisConn.Prefix(key)

	result, err := conn.Do("SET", k, "ok", "EX", int(expiresIn/time.Second), "NX")
	if err != nil {
		return false, err
	}

	return result != nil, nil
}
