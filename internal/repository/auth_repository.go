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

func (r *Auth) SaveRefreshToken(login string, refreshToken string) error {
	user := &models.User{}
	if err := r.db.Model(user).Where("login = ?", login).Update("refresh_token", refreshToken).Error; err != nil {
		return err
	}
	return nil
}
