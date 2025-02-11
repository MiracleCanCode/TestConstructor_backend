package usecases

import (
	"errors"
	"net/http"

	"github.com/server/configs"
	"github.com/server/internal/dtos"
	"github.com/server/internal/repository"
	"github.com/server/pkg/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthInterface interface {
	Login(data *dtos.LoginRequest, w http.ResponseWriter, r *http.Request) (*dtos.LoginResponse, error)
	Registration(data *dtos.RegistrationRequest) (*dtos.RegistrationResponse, error)
}

type Auth struct {
	userRepo      repository.UserInterface
	authRepo      repository.AuthInterface
	logger        *zap.Logger
	tokenProvider jwt.JWTInterface
	config        *configs.Config
}

func NewAuth(
	userRepo repository.UserInterface,
	authRepo repository.AuthInterface,
	logger *zap.Logger,
	tokenProvider jwt.JWTInterface,
	config *configs.Config,
) *Auth {
	return &Auth{
		userRepo:      userRepo,
		authRepo:      authRepo,
		logger:        logger,
		tokenProvider: tokenProvider,
		config:        config,
	}
}

func (s *Auth) Login(data *dtos.LoginRequest, w http.ResponseWriter, r *http.Request) (*dtos.LoginResponse, error) {
	user, err := s.userRepo.GetUserByLogin(data.Login)
	if err != nil || user == nil {
		s.logger.Error("User not found", zap.String("login", data.Login))
		return nil, errors.New("user not found by login")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		s.logger.Error("Invalid password", zap.String("login", data.Login))
		return nil, errors.New("invalid password")
	}

	token, err := s.tokenProvider.CreateAccessToken(user.Login)
	if err != nil {
		s.logger.Error("Failed to create access token", zap.Error(err))
		return nil, errors.New("failed to create access token")
	}

	refreshToken, err := s.tokenProvider.CreateRefreshToken(user.Login)
	if err != nil {
		s.logger.Error("Failed to create refresh token", zap.Error(err))
		return nil, errors.New("failed to generate access token")
	}

	if err := s.authRepo.SaveRefreshToken(user.Login, refreshToken); err != nil {
		s.logger.Error("Failed to save refresh token", zap.Error(err))
		return nil, errors.New("failed to save access token")
	}

	return &dtos.LoginResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *Auth) Registration(data *dtos.RegistrationRequest) (*dtos.RegistrationResponse, error) {
	if err := s.userRepo.CreateUser(data.ToUser()); err != nil {
		s.logger.Error("Failed to register user", zap.Error(err))
		return nil, errors.New("failed to registration user")
	}

	refreshToken, err := s.tokenProvider.CreateRefreshToken(data.Login)
	if err != nil {
		s.logger.Error("Failed to create refresh token", zap.Error(err))
		return nil, errors.New("failed to create refresh token")
	}

	if err := s.authRepo.SaveRefreshToken(data.Login, refreshToken); err != nil {
		s.logger.Error("Failed to save refresh token", zap.Error(err))
		return nil, errors.New("failed to save refresh token")
	}

	return &dtos.RegistrationResponse{
		Name:         data.Name,
		Login:        data.Login,
		Avatar:       data.Avatar,
		Email:        data.Email,
		RefreshToken: refreshToken,
	}, nil
}
