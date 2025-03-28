package cachemanager

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisInterface interface {
	Get(key string) (string, error)
	Set(key string, value interface{}, ttl time.Duration) error
	Del(key string) error
	Close() error
	Keys(key string) *redis.StringSliceCmd
}

type CacheManager struct {
	rdb redisInterface
}

func New(rdb redisInterface) *CacheManager {
	return &CacheManager{
		rdb: rdb,
	}
}

func (s *CacheManager) Get(key string, out interface{}) error {
	data, err := s.rdb.Get(key)
	if err != nil {
		return fmt.Errorf("Get: failed get data from cache:%w", err)
	}
	if err := json.Unmarshal([]byte(data), out); err != nil {
		if delErr := s.rdb.Del(key); delErr != nil {
			return fmt.Errorf("Get failed delete data from cache: %w", err)
		}
		return fmt.Errorf("Get: failed get data from cache:%w", err)
	}
	return nil
}

func (s *CacheManager) Set(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("Set: failed marshal data for redis: %w", err)
	}
	if err := s.rdb.Set(key, data, ttl); err != nil {
		return fmt.Errorf("Set: failed set data to redis: %w", err)
	}

	return nil
}

func (s *CacheManager) Delete(pattern string) error {
	keys, err := s.rdb.Keys(pattern).Result()
	if err != nil {
		return fmt.Errorf("Delete: failed get keys by pattern: %w", err)
	}

	for _, key := range keys {
		if err := s.rdb.Del(key); err != nil {
			return fmt.Errorf("Delete: failed delete data from cache: %w", err)
		}
	}
	return nil
}
