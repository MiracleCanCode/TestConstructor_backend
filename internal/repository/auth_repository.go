package repository

import (
	"fmt"

	"github.com/server/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

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
		return fmt.Errorf("SaveRefreshToken: failed to update user data: %w", err)
	}
	return nil
}
