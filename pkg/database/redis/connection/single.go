package connection

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// 单节点连接
type RedisSingle struct {
	Addr     string
	Auth     string
	Database int
	DB       *redis.Client
}

func (r *RedisSingle) InitSingleRedis(ctx context.Context) (err error) {
	r.DB = redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Auth,
		DB:       r.Database,
		PoolSize: 100,
	})
	// valid connect to redis.manager
	res, err := r.DB.Ping(ctx).Result()
	if err != nil {
		logrus.Fatalf("redis connect error: %s", err.Error())
	}
	logrus.Infof("Connect Redis Successful! Ping => %v", res)
	return err
}

// 释放资源
func (r *RedisSingle) Close() {
	r.DB.Close()
}
