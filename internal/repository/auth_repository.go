package repository

import (
	"github.com/server/internal/models"
	"github.com/server/pkg/db/postgresql"
	"go.uber.org/zap"
)

type AuthInterface interface {
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
	if err := s.db.Where("login = ?", login).First(&user); err != nil {
		return err.Error
	}

	user.RefreshToken = token
	updateResult := s.db.Save(&user)

	if err := updateResult.Error; err != nil {
		s.logger.Error(updateResult.Error.Error(), zap.Error(err))
	}

	return nil
}
