package usecases

import (
	"fmt"
	"net/http"

	"github.com/server/configs"
	"github.com/server/entity"
	"github.com/server/internal/dtos"
	"golang.org/x/crypto/bcrypt"
)

type UserRepoInterface interface {
	CreateUser(user entity.User) error
	DeleteRefreshToken(login string) error
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByLogin(login string) (*entity.User, error)
	UpdateUser(user *dtos.UpdateUserRequest) error
}

type AuthRepoInterface interface {
	SaveRefreshToken(login string, token string) error
}

type JWTInterface interface {
	CreateAccessToken(login string) (string, error)
	CreateRefreshToken(login string) (string, error)
}
type Auth struct {
	userRepo      UserRepoInterface
	authRepo      AuthRepoInterface
	tokenProvider JWTInterface
	config        *configs.Config
}

func NewAuth(
	userRepo UserRepoInterface,
	authRepo AuthRepoInterface,
	tokenProvider JWTInterface,
	config *configs.Config,
) *Auth {
	return &Auth{
		userRepo:      userRepo,
		authRepo:      authRepo,
		tokenProvider: tokenProvider,
		config:        config,
	}
}

func (s *Auth) Login(data *dtos.LoginRequest, w http.ResponseWriter, r *http.Request) (*dtos.LoginResponse, error) {
	user, err := s.userRepo.GetUserByLogin(data.Login)
	if err != nil || user == nil {
		return nil, fmt.Errorf("Login: user not found: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return nil, fmt.Errorf("Login: invalid password: %w", err)
	}

	token, err := s.tokenProvider.CreateAccessToken(user.Login)
	if err != nil {
		return nil, fmt.Errorf("Login: failed to create access token: %w", err)
	}

	refreshToken, err := s.tokenProvider.CreateRefreshToken(user.Login)
	if err != nil {
		return nil, fmt.Errorf("Login: failed to create refresh token: %w", err)
	}

	if err := s.authRepo.SaveRefreshToken(user.Login, refreshToken); err != nil {
		return nil, fmt.Errorf("Login: failed to save refresh token: %w", err)
	}

	return &dtos.LoginResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *Auth) Registration(data *dtos.RegistrationRequest) (*dtos.RegistrationResponse, error) {
	if err := s.userRepo.CreateUser(data.ToUser()); err != nil {
		return nil, fmt.Errorf("Registration: failed to register user: %w", err)
	}

	refreshToken, err := s.tokenProvider.CreateRefreshToken(data.Login)
	if err != nil {
		return nil, fmt.Errorf("Registration: failed to create refresh token: %w", err)
	}

	if err := s.authRepo.SaveRefreshToken(data.Login, refreshToken); err != nil {
		return nil, fmt.Errorf("Registration: failed to save refresh token: %w", err)
	}

	return &dtos.RegistrationResponse{
		Name:         data.Name,
		Login:        data.Login,
		Avatar:       data.Avatar,
		Email:        data.Email,
		RefreshToken: refreshToken,
	}, nil
}
