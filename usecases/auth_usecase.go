package usecases

import (
	"errors"

	"github.com/server/configs"
	"github.com/server/dtos"
	"github.com/server/internal/utils/db/postgresql"
	"github.com/server/internal/utils/jwt"
	"github.com/server/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	db       *postgresql.Db
	log      *zap.Logger
	repo     *repository.Auth
	cfg      *configs.Config
	userRepo *repository.User
}

func NewAuth(db *postgresql.Db, log *zap.Logger, cfg *configs.Config, userRepo *repository.User) *Auth {
	return &Auth{
		db:       db,
		log:      log,
		repo:     repository.NewAuth(db, log),
		cfg:      cfg,
		userRepo: userRepo,
	}
}

func (s *Auth) Login(data *dtos.LoginRequest) (*dtos.LoginResponse, error) {
	if data.Login == "" || data.Password == "" {
		return nil, errors.New("login or password cannot be empty")
	}

	newJwt := jwt.NewJwt(s.log)

	user, err := s.userRepo.GetUserByLogin(data.Login)
	if err != nil {
		s.log.Error("Failed to fetch user", zap.Error(err))
		return nil, errors.New("user not found")
	}
	if user == nil {
		s.log.Error("User not found", zap.String("login", data.Login))
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		s.log.Warn("Invalid password", zap.String("login", data.Login))
		return nil, errors.New("incorrect login or password")
	}

	token, err := newJwt.CreateAccessToken(user.Login)
	if err != nil {
		s.log.Error("Failed to create access token", zap.Error(err))
		return nil, errors.New("internal server error")
	}

	refreshToken, err := newJwt.CreateRefreshToken(user.Login)
	if err != nil {
		s.log.Error("Failed to create refresh token", zap.Error(err))
		return nil, errors.New("internal server error")
	}
	if err := s.repo.SaveRefreshToken(user.Login, refreshToken); err != nil {
		s.log.Error(err.Error())
	}
	return &dtos.LoginResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *Auth) Registration(data *dtos.RegistrationRequest) (*dtos.RegistrationResponse, error) {
	if data.Email == "" || data.Login == "" || data.Password == "" || data.Name == "" {
		return nil, errors.New("exist fields name, login, password, email")
	}

	if err := s.userRepo.CreateUser(data.ToUser()); err != nil {
		s.log.Error("Failed register user")
		return nil, errors.New("failed register")
	}

	return &dtos.RegistrationResponse{}, nil

}
