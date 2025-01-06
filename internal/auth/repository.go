package auth

import (
	"github.com/server/models"
	"github.com/server/pkg/db/postgresql"
	"go.uber.org/zap"
)

type IRepository interface {
	SaveRefreshToken(login string, token string) error
}

type Repository struct {
	db     *postgresql.Db
	logger *zap.Logger
}

func NewRepository(db *postgresql.Db, logger *zap.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (s *Repository) SaveRefreshToken(login string, token string) error {
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
