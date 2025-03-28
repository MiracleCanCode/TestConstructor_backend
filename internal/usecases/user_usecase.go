package usecases

import (
	"fmt"
	"net/http"
	"time"

	"github.com/server/entity"
	"github.com/server/internal/dtos"
	"github.com/server/pkg/constants"
	cookiesmanager "github.com/server/pkg/cookiesManager"
	"github.com/server/pkg/jwt"
	"go.uber.org/zap"
)

type CacheManagerV2Interface interface {
	Get(key string, out interface{}) error
	Set(key string, value interface{}, ttl time.Duration) error
	Delete(pattern string) error
}

type UserRepoInterfaceReaderAndWriter interface {
	GetUserByLogin(login string) (*entity.User, error)
	UpdateUser(user *dtos.UpdateUserRequest) error
	DeleteRefreshToken(login string) error
}

type User struct {
	userRepo     UserRepoInterfaceReaderAndWriter
	cacheManager CacheManagerV2Interface
}

func NewUser(
	userRepo UserRepoInterfaceReaderAndWriter,
	cacheManager CacheManagerV2Interface,
) *User {
	return &User{
		userRepo:     userRepo,
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

	if err := s.cacheManager.Set(cacheKey, user, constants.CACHE_HEALTH_TIME); err != nil {
		return nil, fmt.Errorf("FindUserByLogin: failed set user data to redis: %w", err)
	}

	return user, nil
}

func (s *User) UpdateUserData(user dtos.UpdateUserRequest) error {
	cacheKey := fmt.Sprintf("user:login:%s", user.UserLogin)

	if err := s.cacheManager.Delete(cacheKey); err != nil {
		return fmt.Errorf("UpdateUserData: failed delete user from cache: %w", err)
	}

	if err := s.userRepo.UpdateUser(&user); err != nil {
		return fmt.Errorf("UpdateUserData: failed update user data: %w", err)
	}

	return nil
}

func (s *User) Logout(w http.ResponseWriter, r *http.Request, logger *zap.Logger) error {
	cookies := cookiesmanager.New(r, logger)
	jwt := jwt.NewJwt(logger)
	login, err := jwt.ExtractUserFromToken(r)
	if err != nil {
		return fmt.Errorf("Logout: failed extract user login from token: %w", err)
	}
	if err := s.userRepo.DeleteRefreshToken(login); err != nil {
		return fmt.Errorf("Logout: failed delete refresh token: %w", err)
	}

	cookies.Delete("token", w)
	return nil
}
