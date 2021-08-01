package kvstorage

import (
	"time"
)

type KVStorage interface {
	Store(key string, value interface{}, expiresIn time.Duration) error
	Load(key string, value interface{}) error
	Del(key string) error
}
