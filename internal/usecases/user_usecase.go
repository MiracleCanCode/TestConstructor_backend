package usecases

import (
	"fmt"
	"time"

	"github.com/server/internal/models"
	"github.com/server/internal/repository"
	cachemanager "github.com/server/pkg/cacheManager"
	"go.uber.org/zap"
)

type UserInterface interface {
	FindUserByLogin(login string) (*models.User, error)
}

type User struct {
	userRepo     repository.UserInterface
	logger       *zap.Logger
	cacheManager cachemanager.CacheManagerInterface
}

func NewUser(userRepo repository.UserInterface, logger *zap.Logger, cacheManager cachemanager.CacheManagerInterface) *User {
	return &User{
		userRepo:     userRepo,
		logger:       logger,
		cacheManager: cacheManager,
	}
}

func (s *User) FindUserByLogin(login string) (*models.User, error) {
	cacheKey := fmt.Sprintf("user:login:%s", login)
	var result models.User

	if err := s.cacheManager.Get(cacheKey, &result); err == nil {
		return &result, nil
	}

	user, err := s.userRepo.GetUserByLogin(login)
	if err != nil {
		return nil, err
	}

	s.cacheManager.Set(cacheKey, user, 10*time.Minute)

	return user, nil
}
