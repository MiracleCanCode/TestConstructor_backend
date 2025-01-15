package repository

import (
	"github.com/server/internal/utils/db/postgresql"
	"github.com/server/models"

	"go.uber.org/zap"
)

type IAuth interface {
	SaveRefreshToken(login string, token string) error
}

type Auth struct {
	db     *postgresql.Db
	logger *zap.Logger
}

func NewAuth(db *postgresql.Db, logger *zap.Logger) *Auth {
	return &Auth{
		db:     db,
		logger: logger,
	}
}

func (s *Auth) SaveRefreshToken(login string, token string) error {
	var user models.User
	result := s.db.Where("login = ?", login).First(&user)
	if result.Error != nil {
		return result.Error
	}

	user.RefreshToken = token
	updateResult := s.db.Save(&user)

	if err := updateResult.Error; err != nil {
		s.logger.Error(updateResult.Error.Error())
	}

	return nil
}