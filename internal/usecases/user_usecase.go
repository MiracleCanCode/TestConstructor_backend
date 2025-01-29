package usecases

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/server/internal/models"
	"github.com/server/internal/repository"
	"github.com/server/pkg/db/redis"
	"go.uber.org/zap"
)

type UserInterface interface {
	FindUserData(login string) (*models.User, error)
	FindUserByLogin(login string) (*models.User, error)
}

type User struct {
	userRepo repository.UserInterface
	logger   *zap.Logger
}

func NewUser(userRepo repository.UserInterface, logger *zap.Logger) *User {
	return &User{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *User) FindUserByLogin(login string) (*models.User, error) {
	rdb := redis.New()
	cacheKey := fmt.Sprintf("user:login:%s", login)

	cachedData, err := rdb.Get(cacheKey)
	if err == nil {
		var result struct {
			User *models.User `json:"user"`
		}

		if err := json.Unmarshal([]byte(cachedData), &result); err == nil {
			return result.User, nil
		}

		s.logger.Error("Failed to unmarshal cache", zap.Error(err))
	}

	user, err := s.userRepo.GetUserByLogin(login)
	if err != nil {
		return nil, err
	}

	cacheValue, err := json.Marshal(map[string]interface{}{
		"user": user,
	})
	if err == nil {
		_ = rdb.Set(cacheKey, cacheValue, 10*time.Minute)
	} else {
		s.logger.Warn("Failed to marshal cache data", zap.Error(err))
	}

	return user, nil
}
