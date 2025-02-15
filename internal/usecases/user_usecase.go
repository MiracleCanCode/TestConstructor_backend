package usecases

import (
	"fmt"
	"net/http"
	"time"

	"github.com/server/entity"
	"github.com/server/internal/dtos"
	cookiesmanager "github.com/server/pkg/cookiesManager"
	"github.com/server/pkg/jwt"
	"go.uber.org/zap"
)

type CacheManagerV2Interface interface {
	Get(key string, out interface{}) error
	Set(key string, value interface{}, ttl time.Duration)
	Delete(pattern string) error
}

type UserRepoInterfaceReaderAndWriter interface {
	GetUserByLogin(login string) (*entity.User, error)
	UpdateUser(user *dtos.UpdateUserRequest) error
	DeleteRefreshToken(login string) error
}

type User struct {
	userRepo     UserRepoInterfaceReaderAndWriter
	logger       *zap.Logger
	cacheManager CacheManagerV2Interface
}

func NewUser(
	userRepo UserRepoInterfaceReaderAndWriter,
	logger *zap.Logger,
	cacheManager CacheManagerV2Interface,
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
		return nil, fmt.Errorf("FindUserByLogin: failed to find user by login: %w", err)
	}

	s.cacheManager.Set(cacheKey, user, 10*time.Minute)

	return user, nil
}

func (s *User) UpdateUserData(user dtos.UpdateUserRequest) error {
	cacheKey := fmt.Sprintf("user:login:%s", user.UserLogin)

	if err := s.cacheManager.Delete(cacheKey); err != nil {
		return fmt.Errorf("UpdateUserData: failed delete user from cache: %w", err)
	}

	return s.userRepo.UpdateUser(&user)
}

func (s *User) Logout(w http.ResponseWriter, r *http.Request) error {
	cookies := cookiesmanager.New(r, s.logger)
	jwt := jwt.NewJwt(s.logger)
	login, err := jwt.ExtractUserFromToken(r)
	if err != nil {
		return fmt.Errorf("Logout: failed extract user login from token: %w", err)
	}
	cookies.Delete("token", w)
	if err := s.userRepo.DeleteRefreshToken(login); err != nil {
		return fmt.Errorf("Logout: failed delete refresh token: %w", err)
	}

	return nil
}
