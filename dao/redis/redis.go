// Package redis 相关连接
package redis

import (
	"context"
	"fmt"

	"blueblog/settings"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	ctx    context.Context
)

// Init 初始化连接
func Init(cfg *settings.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx = context.Background()

	_, err = client.Ping(context.Background()).Result()
	return
}

// 关闭连接
func Close() {
	_ = client.Close()
}
