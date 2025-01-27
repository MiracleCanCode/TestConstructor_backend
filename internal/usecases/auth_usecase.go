package usecases

import (
	"errors"

	"github.com/server/configs"
	"github.com/server/internal/dtos"
	"github.com/server/internal/repository"
	"github.com/server/pkg/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthInterface interface {
	Login(data *dtos.LoginRequest) (*dtos.LoginResponse, error)
	Registration(data *dtos.RegistrationRequest) (*dtos.RegistrationResponse, error)
}

type Auth struct {
	userRepo      repository.UserInterface
	authRepo      repository.AuthInterface
	logger        *zap.Logger
	tokenProvider jwt.JWT
	config        *configs.Config
}

func NewAuth(
	userRepo repository.UserInterface,
	authRepo repository.AuthInterface,
	logger *zap.Logger,
	tokenProvider jwt.JWT,
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

func (uc *Auth) Login(data *dtos.LoginRequest) (*dtos.LoginResponse, error) {
	user, err := uc.userRepo.GetUserByLogin(data.Login)
	if err != nil || user == nil {
		uc.logger.Warn("User not found", zap.String("login", data.Login))
		return nil, errors.New("incorrect login or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		uc.logger.Warn("Invalid password", zap.String("login", data.Login))
		return nil, errors.New("incorrect login or password")
	}

	token, err := uc.tokenProvider.CreateAccessToken(user.Login)
	if err != nil {
		uc.logger.Error("Failed to create access token", zap.Error(err))
		return nil, errors.New("internal server error")
	}

	refreshToken, err := uc.tokenProvider.CreateRefreshToken(user.Login)
	if err != nil {
		uc.logger.Error("Failed to create refresh token", zap.Error(err))
		return nil, errors.New("internal server error")
	}

	if err := uc.authRepo.SaveRefreshToken(user.Login, refreshToken); err != nil {
		uc.logger.Error("Failed to save refresh token", zap.Error(err))
	}

	return &dtos.LoginResponse{
		User:  user,
		Token: token,
	}, nil
}

func (uc *Auth) Registration(data *dtos.RegistrationRequest) (*dtos.RegistrationResponse, error) {
	if err := uc.userRepo.CreateUser(data.ToUser()); err != nil {
		uc.logger.Error("Failed to register user", zap.Error(err))
		return nil, errors.New("failed to register user")
	}
	return &dtos.RegistrationResponse{}, nil
}