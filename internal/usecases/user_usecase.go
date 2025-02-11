package usecases

import (
	"fmt"
	"net/http"
	"time"

	"github.com/server/entity"
	"github.com/server/internal/dtos"
	"github.com/server/internal/repository"
	cachemanager "github.com/server/pkg/cacheManager"
	cookiesmanager "github.com/server/pkg/cookiesManager"
	"github.com/server/pkg/jwt"
	"go.uber.org/zap"
)

type UserInterface interface {
	FindUserByLogin(login string) (*entity.User, error)
	UpdateUserData(user dtos.UpdateUserRequest) error
	Logout(w http.ResponseWriter, r *http.Request) error
}

type User struct {
	userRepo     repository.UserInterface
	logger       *zap.Logger
	cacheManager cachemanager.CacheManagerInterface
}

func NewUser(
	userRepo repository.UserInterface,
	logger *zap.Logger,
	cacheManager cachemanager.CacheManagerInterface,
) *User {
	return &User{
		userRepo:     userRepo,
		logger:       logger,
		cacheManager: cacheManager,
	}
}

func (s *User) FindUserByLogin(login string) (*entity.User, error) {
	cacheKey := fmt.Sprintf("user:login:%s", login)
	var result entity.User

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

func (s *User) UpdateUserData(user dtos.UpdateUserRequest) error {
	cacheKey := fmt.Sprintf("user:login:%s", user.UserLogin)

	if err := s.cacheManager.Delete(cacheKey); err != nil {
		return err
	}

	return s.userRepo.UpdateUser(&user)
}

func (s *User) Logout(w http.ResponseWriter, r *http.Request) error {
	cookies := cookiesmanager.New(r, s.logger)
	jwt := jwt.NewJwt(s.logger)
	login, err := jwt.ExtractUserFromToken(r)
	if err != nil {
		s.logger.Error("Extract user login from token", zap.Error(err))
		return err
	}
	cookies.Delete("token", w)
	if err := s.userRepo.DeleteRefreshToken(login); err != nil {
		s.logger.Error("Delete refresh token", zap.Error(err))
		return err
	}

	return nil
}
