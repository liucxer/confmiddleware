package memory

import (
	"reflect"
	"sync"
	"time"

	"github.com/liucxer/courier/reflectx"

	"github.com/liucxer/confmiddleware/kvstorage"
)

var _ kvstorage.KVStorage = (*MemoryKVStorage)(nil)

type MemoryKVStorage struct {
	m sync.Map
}

type ValueWithExpire struct {
	Value     interface{}
	Always    bool
	ExpiredAt time.Time
}

func (cache *MemoryKVStorage) Del(key string) error {
	cache.m.Delete(key)
	return nil
}

func (cache *MemoryKVStorage) Store(key string, value interface{}, expiresIn time.Duration) error {
	if expiresIn > 0 {
		cache.m.Store(key, ValueWithExpire{
			Value:     value,
			ExpiredAt: time.Now().Add(expiresIn),
		})
		return nil
	}
	cache.m.Store(key, ValueWithExpire{
		Value:  value,
		Always: true,
	})
	return nil
}

func (cache *MemoryKVStorage) Load(key string, value interface{}) error {
	if val, ok := cache.m.Load(key); ok {
		v := val.(ValueWithExpire)
		if !v.Always {
			if time.Now().After(v.ExpiredAt) {
				cache.Del(key)
				return nil
			}
		}
		rv := reflectx.Indirect(reflect.ValueOf(value))
		rv.Set(reflectx.Indirect(reflect.ValueOf(v.Value)))
		return nil
	}
	return nil
}
