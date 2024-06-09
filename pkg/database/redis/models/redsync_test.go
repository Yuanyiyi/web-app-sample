package models

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetMutex(t *testing.T) {
	os.Setenv("REDIS_ADDR", "172.16.32.177:6379")
	os.Setenv("REDIS_PASSWORD", "TNezF2piCu")

	key := "00001"
	mutex1, res1 := SetMutex(key, 0)
	assert.Equal(t, true, res1)

	_, res2 := SetMutex(key, 0)
	assert.Equal(t, false, res2)

	res3 := DelMutex(mutex1)
	assert.Equal(t, true, res3)

	_, res4 := SetMutex(key, 0)
	assert.Equal(t, true, res4)

	time.Sleep(time.Second)
	_, res4 = SetMutex(key, 0)
	assert.Equal(t, true, res4)
}
