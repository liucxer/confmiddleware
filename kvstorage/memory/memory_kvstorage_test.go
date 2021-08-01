package memory

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMemoryKVStorage(t *testing.T) {
	tt := require.New(t)

	c := &MemoryKVStorage{}

	key := "key"
	value := "value"

	{
		err := c.Store(key, value, -1)
		tt.NoError(err)

		v := ""
		tt.NoError(c.Load(key, &v))
		tt.Equal(value, v)
	}

	{
		err := c.Store(key, value, 1*time.Second)
		tt.NoError(err)
		time.Sleep(2 * time.Second)

		v := ""
		tt.NoError(c.Load(key, &v))
		tt.Empty(v)
	}

	{
		err := c.Store(key, value, -1)
		tt.NoError(err)

		errForDelete := c.Del(key)
		tt.NoError(errForDelete)

		v := ""
		tt.NoError(c.Load(key, &v))
		tt.Empty(v)
	}
}
