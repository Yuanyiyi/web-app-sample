package connection

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// 哨兵模式连接
// 定义一个RedisClusterObj结构体
type RedisSentinel struct {
	Master string // master- 读写操作
	Addr   []string
	Auth   string
	DB     *redis.Client
}

func (r *RedisSentinel) InitSentinelClient(ctx context.Context) (err error) {
	r.DB = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    r.Master,
		SentinelAddrs: r.Addr,
		DB:            0,
	})
	// valid connect to redis.manager
	res, err := r.DB.Ping(ctx).Result()
	if err != nil {
		logrus.Fatalf("redis connect error: %s", err.Error())
	}
	logrus.Infof("Connect Successful! Ping => %v", res)
	return err
}

func (r *RedisSentinel) Close() {
	r.DB.Close()
}
