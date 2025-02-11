package repository

import (
	"github.com/server/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthInterface interface {
	SaveRefreshToken(login string, token string) error
}

type Auth struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewAuth(db *gorm.DB, logger *zap.Logger) *Auth {
	return &Auth{
		db:     db,
		logger: logger,
	}
}

func (s *Auth) SaveRefreshToken(login string, refreshToken string) error {
	user := &entity.User{}
	if err := s.db.Model(user).Where("login = ?", login).Update("refresh_token", refreshToken).Error; err != nil {
		return err
	}
	return nil
}
