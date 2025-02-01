package cachemanager

import (
	"encoding/json"
	"time"

	"github.com/server/pkg/db/redis"
	"go.uber.org/zap"
)

type CacheManagerInterface interface {
	Get(key string, out interface{}) error
	Set(key string, value interface{}, ttl time.Duration)
	Delete(pattern string) error
}

type CacheManager struct {
	rdb    redis.RedisInterface
	logger *zap.Logger
}

func New(rdb redis.RedisInterface, logger *zap.Logger) *CacheManager {
	return &CacheManager{
		rdb:    rdb,
		logger: logger,
	}
}

func (s *CacheManager) Get(key string, out interface{}) error {
	data, err := s.rdb.Get(key)
	if err != nil {
		return err
	}
	s.logger.Info("Cache get")
	if err := json.Unmarshal([]byte(data), out); err != nil {
		s.logger.Warn("Failed to unmarshal cache data, deleting key", zap.String("key", key), zap.Error(err))
		if delErr := s.rdb.Del(key); delErr != nil {
			s.logger.Warn("Failed to delete corrupted cache key", zap.String("key", key), zap.Error(delErr))
		}
		return err
	}
	return nil
}

func (s *CacheManager) Set(key string, value interface{}, ttl time.Duration) {
	data, err := json.Marshal(value)
	s.logger.Info("Cache set")
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
		s.logger.Error("Failed to get cache keys", zap.Error(err))
		return err
	}

	for _, key := range keys {
		if err := s.rdb.Del(key); err != nil {
			s.logger.Warn("Failed to delete cache key", zap.String("key", key), zap.Error(err))
		}
	}
	return nil
}
