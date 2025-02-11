package redis

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/server/configs"
	"github.com/server/pkg/logger"
)

type RedisInterface interface {
	Get(key string) (string, error)
	Set(key string, value interface{}, ttl time.Duration) error
	Del(key string) error
	Close() error
	Keys(key string) *redis.StringSliceCmd
}

type Redis struct {
	client *redis.Client
}

var (
	instance *redis.Client
	once     sync.Once
)

func New() *Redis {
	once.Do(func() {
		log := logger.GetInstance()
		cfg, err := configs.Load(log)
		if err != nil {
			return
		}
		instance = redis.NewClient(&redis.Options{
			Addr:     cfg.REDIS_HOST,
			Password: "",
			DB:       0,
		})
	})
	return &Redis{
		client: instance,
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

func (r *Redis) Keys(key string) *redis.StringSliceCmd {
	ctx := context.Background()
	return r.client.Keys(ctx, key)
}
