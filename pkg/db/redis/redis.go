package redis

import (
	"context"
	"time"

	"github.com/MiracleCanCode/example_configuration_logger/pkg/logger"
	"github.com/redis/go-redis/v9"
	"github.com/server/configs"
)

type Redis struct {
	client *redis.Client
}

func New() *Redis {
	log := logger.Logger(logger.DefaultLoggerConfig())
	cfg, err := configs.Load(log)
	if err != nil {
		return nil
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.REDIS_HOST,
		Password: "",
		DB:       0,
	})

	return &Redis{
		client: rdb,
	}
}

func (r *Redis) Get(key string) (string, error) {
	ctx := context.Background()
	return r.client.Get(ctx, key).Result()
}

func (r *Redis) Set(key string, value interface{}, ttl time.Duration) error {
	ctx := context.Background()
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *Redis) Del(key string) error {
	ctx := context.Background()
	return r.client.Del(ctx, key).Err()
}

func (r *Redis) Close() error {
	return r.client.Close()
}
