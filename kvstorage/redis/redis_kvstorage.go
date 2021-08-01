package redis

import (
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"

	"github.com/liucxer/confmiddleware/kvstorage"
)

type RedisOperator interface {
	// key prefix
	Prefix(key string) string
	// get redis connect
	Get() redis.Conn
}

var _ kvstorage.KVStorage = (*RedisKVStorage)(nil)

func NewRedisKVStorage(op RedisOperator) *RedisKVStorage {
	return &RedisKVStorage{
		op: op,
	}
}

type RedisKVStorage struct {
	op RedisOperator
}

func (s *RedisKVStorage) Do(cmd string, args ...interface{}) (interface{}, error) {
	c := s.op.Get()
	defer c.Close()
	res, err := c.Do(cmd, args...)
	if err != nil {
		logrus.WithField("redis_cmd", cmd).Warningf("%v", err.Error())
	}
	return res, err
}

func (s *RedisKVStorage) Del(key string) error {
	_, err := s.Do("DEL", s.op.Prefix(key))
	return err
}

type data struct {
	Value interface{} `json:"value"`
}

func (s *RedisKVStorage) Store(key string, value interface{}, expiresIn time.Duration) error {
	bytes, err := json.Marshal(data{Value: value})
	if err != nil {
		return err
	}

	if expiresIn > 0 {
		_, err := s.Do("SET", s.op.Prefix(key), string(bytes), "EX", transToSecond(expiresIn))
		if err != nil {
			return err
		}
		return nil
	}

	_, errForExec := s.Do("SET", s.op.Prefix(key), bytes)
	if errForExec != nil {
		return errForExec
	}

	return nil
}

func transToSecond(dur time.Duration) int64 {
	if dur > 0 && dur < time.Second {
		return 0
	}
	return int64(dur / time.Second)
}

func (s *RedisKVStorage) Load(key string, value interface{}) error {
	bytes, err := redis.Bytes(s.Do("GET", s.op.Prefix(key)))
	if err != nil {
		if err == redis.ErrNil {
			return nil
		}
		return err
	}
	return json.Unmarshal(bytes, &data{Value: value})
}
