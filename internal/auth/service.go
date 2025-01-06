package auth

import (
	"errors"

	"github.com/server/configs"
	"github.com/server/internal/user"
	"github.com/server/pkg/db/postgresql"
	"github.com/server/pkg/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db           *postgresql.Db
	log          *zap.Logger
	repo         *Repository
	cfg          *configs.Config
	userRepo *user.Repository
}

func NewService(db *postgresql.Db, log *zap.Logger, cfg *configs.Config, userRepo *user.Repository) *Service {
	return &Service{
		db:   db,
		log:  log,
		repo: NewRepository(db, log),
		cfg:  cfg,
		userRepo: userRepo,
	}
}

func (s *Service) Login(data *LoginRequest) (*LoginResponse, error) {
	if data.Login == "" || data.Password == "" {
		return nil, errors.New("login or password cannot be empty")
	}

	newJwt := jwt.NewJwt("supersecretkeybysuperuser")

	user, err := s.userRepo.GetByLogin(data.Login)
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
	return &LoginResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *Service) Registration(data *RegistrationRequest) (*RegistrationResponse, error) {
	if data.Email == "" || data.Login == "" || data.Password == "" || data.Name == "" {
		return nil, errors.New("exist fields name, login, password, email")
	}

	if err := s.userRepo.Create(data.ToUser()); err != nil {
		s.log.Error("Failed register user")
		return nil, errors.New("failed register")
	}	

	return &RegistrationResponse{}, nil

}
