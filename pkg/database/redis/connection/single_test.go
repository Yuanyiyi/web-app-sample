package connection

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRedisSingle(t *testing.T) {
	resdisSingle := &RedisSingle{
		Addr: "172.16.32.177:6379",
		Auth: "TNezF2piCu",
	}
	err := resdisSingle.InitSingleRedis(context.Background())
	assert.Nil(t, err)
}
