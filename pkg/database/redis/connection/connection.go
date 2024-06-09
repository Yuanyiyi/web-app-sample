package connection

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strings"

	"github.com/web-app-sample/pkg/utils/startenv"
)

var (
	redisDB        *redis.Client
	redisClusterDB *redis.ClusterClient
)

func redisConnection(ctx context.Context) {
	addr := startenv.GetRedisAddr()
	switch startenv.GetRedisDeployment() {
	case "", "single":
		redisSingle := &RedisSingle{
			Addr:     addr,
			Auth:     startenv.GetRedisPassword(),
			Database: 0,
		}
		redisSingle.InitSingleRedis(ctx)
		redisDB = redisSingle.DB
	case "sentinel":
		addrs := strings.Split(addr, ",")
		redisSentinel := &RedisSentinel{
			Master: "master",
			Addr:   addrs,
			Auth:   startenv.GetRedisPassword(),
		}
		redisSentinel.InitSentinelClient(ctx)
		redisDB = redisSentinel.DB
	case "cluster":
		addrs := strings.Split(addr, ",")
		redisCluster := &RedisCluster{
			Addr: addrs,
			Auth: startenv.GetRedisPassword(),
		}
		redisClusterDB = redisCluster.DB
	}
}

func GetRedisDB(ctx context.Context) *redis.Client {
	if redisDB == nil {
		redisConnection(ctx)
	}
	return redisDB
}

func GetRedisClusterDB(ctx context.Context) *redis.ClusterClient {
	if redisClusterDB == nil {
		redisConnection(ctx)
	}
	return redisClusterDB
}
