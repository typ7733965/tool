package cache

import (
	"github.com/go-redis/redis"
	"github.com/typ7733965/tool/config"
)

func InitRedis(conf *config.RedisConfig) redis.UniversalClient {
	if conf.IsCluster {
		return redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:         conf.Addrs,
			Password:      conf.Password,
			PoolSize:      conf.PoolSize,
			MinIdleConns:  conf.MinIdleConn,
			RouteRandomly: conf.RouteRandomly,
		})
	} else {
		return redis.NewClient(&redis.Options{
			Addr:         conf.Addrs[0],
			Password:     conf.Password,
			DB:           conf.DB,
			PoolSize:     conf.PoolSize,
			MinIdleConns: conf.MinIdleConn,
		})
	}
}
