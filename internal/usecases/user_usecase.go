package usecases

import (
	"fmt"
	"time"

	"github.com/server/internal/dtos"
	"github.com/server/internal/models"
	"github.com/server/internal/repository"
	cachemanager "github.com/server/pkg/cacheManager"
	"go.uber.org/zap"
)

type UserInterface interface {
	FindUserByLogin(login string) (*models.User, error)
	UpdateUserData(user dtos.UpdateUserRequest) error
}

type User struct {
	userReader   repository.UserReader
	userWriter   repository.UserWriter
	logger       *zap.Logger
	cacheManager cachemanager.CacheManagerInterface
}

func NewUser(
	userReader repository.UserReader,
	userWriter repository.UserWriter,
	logger *zap.Logger,
	cacheManager cachemanager.CacheManagerInterface,
) *User {
	return &User{
		userReader:   userReader,
		userWriter:   userWriter,
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

	user, err := s.userReader.GetUserByLogin(login)
	if err != nil {
		return nil, err
	}

	s.cacheManager.Set(cacheKey, user, 10*time.Minute)

	return user, nil
}

func (s *User) UpdateUserData(user dtos.UpdateUserRequest) error {
	cacheKey := fmt.Sprintf("user:login:%s", user.UserLogin)

	if err := s.cacheManager.Delete(cacheKey); err != nil {
		return err
	}

	return s.userWriter.UpdateUser(&user)
}
