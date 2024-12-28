package auth

import (
	"github.com/MiracleCanCode/zaperr"
	"github.com/server/models"
	"github.com/server/pkg/db"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository struct {
	db           *db.Db
	logger       *zap.Logger
	handleErrors *zaperr.Zaperr
}

func NewAuthRepository(db *db.Db, logger *zap.Logger, handleErrors *zaperr.Zaperr) *AuthRepository {
	return &AuthRepository{
		db:           db,
		logger:       logger,
		handleErrors: handleErrors,
	}
}

func (s *AuthRepository) CreateUser(user models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	s.handleErrors.LogError(err, string(err.Error()))
	user.Password = string(hashedPassword)

	result := s.db.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *AuthRepository) GetUserByLogin(login string) (*models.User, error) {

	var user models.User
	result := s.db.Where("login = ?", login).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (s *AuthRepository) SaveRefreshToken(login string, token string) error {
	var user models.User
	result := s.db.Where("login = ?", login).First(&user)
	if result.Error != nil {
		return result.Error
	}

	user.RefreshToken = token
	updateResult := s.db.Save(&user)

	s.handleErrors.LogError(updateResult.Error, string(updateResult.Error.Error()))

	return nil
}
