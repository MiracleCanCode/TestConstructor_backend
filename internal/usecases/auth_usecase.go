package usecases

import (
	"errors"
	"net/http"

	"github.com/server/configs"
	"github.com/server/internal/dtos"
	"github.com/server/internal/repository"
	"github.com/server/pkg/cookie"
	errorconstant "github.com/server/pkg/errorConstants"
	"github.com/server/pkg/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthInterface interface {
	Login(data *dtos.LoginRequest, w http.ResponseWriter, r *http.Request) (*dtos.LoginResponse, error)
	Registration(data *dtos.RegistrationRequest) (*dtos.RegistrationResponse, error)
	Logout(w http.ResponseWriter, r *http.Request) error
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

func (s *Auth) Login(data *dtos.LoginRequest, w http.ResponseWriter, r *http.Request) (*dtos.LoginResponse, error) {
	cookies := cookie.New(w, r, s.logger)
	user, err := s.userRepo.GetUserByLogin(data.Login)
	if err != nil || user == nil {
		s.logger.Warn("User not found", zap.String("login", data.Login))
		return nil, errors.New(errorconstant.ErrInvalidCredentials.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		s.logger.Warn("Invalid password", zap.String("login", data.Login))
		return nil, errors.New(errorconstant.ErrInvalidCredentials.Error())
	}

	token, err := s.tokenProvider.CreateAccessToken(user.Login)
	if err != nil {
		s.logger.Error("Failed to create access token", zap.Error(err))
		return nil, errors.New(errorconstant.ErrInternalServer.Error())
	}

	cookies.Set("token", token)

	refreshToken, err := s.tokenProvider.CreateRefreshToken(user.Login)
	if err != nil {
		s.logger.Error("Failed to create refresh token", zap.Error(err))
		return nil, errors.New(errorconstant.ErrInternalServer.Error())
	}

	if err := s.authRepo.SaveRefreshToken(user.Login, refreshToken); err != nil {
		s.logger.Error("Failed to save refresh token", zap.Error(err))
		return nil, errors.New(errorconstant.ErrInternalServer.Error())
	}

	return &dtos.LoginResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *Auth) Registration(data *dtos.RegistrationRequest) (*dtos.RegistrationResponse, error) {
	if err := s.userRepo.CreateUser(data.ToUser()); err != nil {
		s.logger.Error("Failed to register user", zap.Error(err))
		return nil, errors.New(errorconstant.ErrRegisterUser.Error())
	}

	refreshToken, err := s.tokenProvider.CreateRefreshToken(data.Login)
	if err != nil {
		s.logger.Error("Failed to create refresh token", zap.Error(err))
		return nil, errors.New(errorconstant.ErrInternalServer.Error())
	}

	if err := s.authRepo.SaveRefreshToken(data.Login, refreshToken); err != nil {
		s.logger.Error("Failed to save refresh token", zap.Error(err))
		return nil, errors.New(errorconstant.ErrInternalServer.Error())
	}

	return &dtos.RegistrationResponse{
		Name:         data.Name,
		Login:        data.Login,
		Avatar:       data.Avatar,
		Email:        data.Email,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Auth) Logout(w http.ResponseWriter, r *http.Request) error {
	cookies := cookie.New(w, r, s.logger)
	JWT := jwt.NewJwt(s.logger)

	login, err := JWT.ExtractUserFromCookie(r, "token")
	if err != nil {
		s.logger.Error("Failed to extract access_token form cookie", zap.Error(err))
		return errors.New(errorconstant.ErrInternalServer.Error())
	}

	cookies.Delete("token")

	if err := s.userRepo.DeleteRefreshToken(login); err != nil {
		s.logger.Error("Failed to delete refresh token", zap.Error(err))
		return errors.New(errorconstant.ErrInternalServer.Error())
	}

	return nil
}
