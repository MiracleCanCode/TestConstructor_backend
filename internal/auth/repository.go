package auth

import (
	"github.com/server/models"
	"github.com/server/pkg/db/postgresql"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type IRepository interface {
	CreateUser(user models.User) error
	GetUserByLogin(login string) (*models.User, error)
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

func (s *Repository) CreateUser(user models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	result := s.db.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *Repository) GetUserByLogin(login string) (*models.User, error) {

	var user models.User
	result := s.db.Where("login = ?", login).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
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
		s.logger.Error(string(updateResult.Error.Error()))
	}

	return nil
}
