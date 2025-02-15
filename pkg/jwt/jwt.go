package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/server/configs"
	cookiesmanager "github.com/server/pkg/cookiesManager"
	"go.uber.org/zap"
)

type JWT struct {
	Secret string
	logger *zap.Logger
}

func NewJwt(logger *zap.Logger) *JWT {
	cfg, err := configs.Load(logger)
	if err != nil {
		logger.Error("Failed to load config", zap.Error(err))
		return nil
	}
	return &JWT{
		Secret: cfg.SECRET,
		logger: logger,
	}
}

func (s *JWT) createToken(login string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"login": login,
		"exp":   time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.Secret))
}

func (s *JWT) CreateAccessToken(login string) (string, error) {
	return s.createToken(login, time.Hour*24)
}

func (s *JWT) CreateRefreshToken(login string) (string, error) {
	return s.createToken(login, time.Hour*24*7)
}

func (s *JWT) VerifyToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.Secret), nil
	})

	if err != nil {
		return nil, nil, fmt.Errorf("VerifyToken: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, nil, jwt.ErrInvalidKey
	}

	return token, claims, nil
}

func (s *JWT) ExtractUserFromToken(r *http.Request) (string, error) {
	cookie := cookiesmanager.New(r, s.logger)
	authToken, err := cookie.Get("token")
	if err != nil {
		return "", errors.New("failed extract token from cookie")
	}
	_, claims, err := s.VerifyToken(authToken)
	if err != nil {
		return "", fmt.Errorf("ExtractUserFromToken: %w", err)
	}
	userLogin, ok := claims["login"].(string)
	if !ok {
		return "", fmt.Errorf("ExtractUserFromToken: %w", err)
	}

	return userLogin, nil
}

func (s *JWT) RefreshAccessToken(refreshToken string) (string, error) {
	_, claims, err := s.VerifyToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("RefreshAccessToken: %w", err)
	}

	userLogin, ok := claims["login"].(string)
	if !ok {
		return "", fmt.Errorf("RefreshAccessToken: %w", err)
	}

	return s.CreateAccessToken(userLogin)
}
