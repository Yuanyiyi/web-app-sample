package connection

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// redis集群模式连接
// 定义一个redisCluster结构体
type RedisCluster struct {
	Addr []string
	Auth string
	DB   *redis.ClusterClient
}

func (r *RedisCluster) InitClusterClient(ctx context.Context) (err error) {
	r.DB = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    r.Addr,
		Password: r.Auth,
	})
	// valid connect to redis.manager
	res, err := r.DB.Ping(ctx).Result()
	if err != nil {
		logrus.Fatalf("redis connect error: %s", err.Error())
	}
	logrus.Infof("Connect Successful! Ping => %v\n", res)
	return err
}

func (r *RedisCluster) Close() {
	r.DB.Close()
}
