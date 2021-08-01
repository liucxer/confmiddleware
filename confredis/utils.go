package confredis

import "github.com/gomodule/redigo/redis"

type Conn = redis.Conn

func Command(name string, args ...interface{}) *CMD {
	return &CMD{
		name: name,
		args: args,
	}
}

type CMD struct {
	name string
	args []interface{}
}

type RedisOperator interface {
	// key prefix
	Prefix(key string) string
	// get redis connect
	Get() Conn

	// exec
	Exec(cmd *CMD, others ...*CMD) (interface{}, error)
}
