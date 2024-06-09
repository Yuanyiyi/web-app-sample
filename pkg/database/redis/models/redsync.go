package models

import (
	"context"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/sirupsen/logrus"

	redisConn "github.com/web-app-sample/pkg/database/redis/connection"
)

var (
	redSync *redsync.Redsync
)

func initRedSync() {
	redisDB := redisConn.GetRedisDB(context.Background())
	// create redsync client pool
	pool := goredis.NewPool(redisDB)

	// create redsync instance
	redSync = redsync.New(pool)
}

func GetRedSync() *redsync.Redsync {
	if redSync == nil {
		initRedSync()
	}
	return redSync
}

// 分布式锁
func SetMutex(key string, t1 int) (*redsync.Mutex, bool) {
	if redSync == nil {
		initRedSync()
	}
	if t1 <= 0 {
		t1 = 1000
	}
	options := []redsync.Option{redsync.WithTries(3), redsync.WithExpiry(time.Duration(t1) * time.Millisecond)}
	// create mutex for key
	mutex := redSync.NewMutex(key, options...)

	if err := mutex.Lock(); err != nil {
		return mutex, false
	}
	return mutex, true
}

func DelMutex(mutex *redsync.Mutex) bool {
	if redSync == nil {
		initRedSync()
	}

	if ok, err := mutex.Unlock(); !ok || err != nil {
		logrus.Errorf("redsync del error: %s", err)
		return false
	}
	return true
}
