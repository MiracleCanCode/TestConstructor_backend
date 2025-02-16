package cachemanager

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisInterface interface {
	Get(key string) (string, error)
	Set(key string, value interface{}, ttl time.Duration) error
	Del(key string) error
	Close() error
	Keys(key string) *redis.StringSliceCmd
}

type CacheManager struct {
	rdb    RedisInterface
	logger *zap.Logger
}

func New(rdb RedisInterface, logger *zap.Logger) *CacheManager {
	return &CacheManager{
		rdb:    rdb,
		logger: logger,
	}
}

func (s *CacheManager) Get(key string, out interface{}) error {
	data, err := s.rdb.Get(key)
	if err != nil {
		return fmt.Errorf("Get: failed get data from cache:%w", err)
	}
	if err := json.Unmarshal([]byte(data), out); err != nil {
		if delErr := s.rdb.Del(key); delErr != nil {
			s.logger.Warn("Failed to delete corrupted cache key", zap.String("key", key), zap.Error(delErr))
		}
		return fmt.Errorf("Get: failed get data from cache:%w", err)
	}
	return nil
}

func (s *CacheManager) Set(key string, value interface{}, ttl time.Duration) {
	data, err := json.Marshal(value)
	if err != nil {
		s.logger.Warn("Failed to marshal data for cache", zap.String("key", key), zap.Error(err))
		return
	}
	if err := s.rdb.Set(key, data, ttl); err != nil {
		s.logger.Warn("Failed to cache data", zap.String("key", key), zap.Error(err))
	}
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
