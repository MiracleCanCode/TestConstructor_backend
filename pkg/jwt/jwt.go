package jwt

import (
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
	token, err := s.createToken(login, time.Hour*24)
	if err != nil {
		return "", fmt.Errorf("CreateAccessToken: failed to create access token: %w", err)
	}
	return token, nil
}

func (s *JWT) CreateRefreshToken(login string) (string, error) {
	token, err := s.createToken(login, time.Hour*24*7)
	if err != nil {
		return "", fmt.Errorf("CreateRefreshToken: failed to create refresh token: %w", err)
	}
	return token, nil
}

func (s *JWT) VerifyToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.Secret), nil
	})

	if err != nil {
		return nil, nil, fmt.Errorf("VerifyToken: failed verify auth token: %w", err)
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
		return "", fmt.Errorf("ExtractUserFromToken: failed extract token from cookie: %w", err)
	}
	_, claims, err := s.VerifyToken(authToken)
	if err != nil {
		return "", fmt.Errorf("ExtractUserFromToken: failed verify auth token: %w", err)
	}

	userLogin, ok := claims["login"].(string)
	if !ok {
		return "", fmt.Errorf("ExtractUserFromToken: failed get login from claims: %w", err)
	}

	return userLogin, nil
}

func (s *JWT) RefreshAccessToken(refreshToken string) (string, error) {
	_, claims, err := s.VerifyToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("RefreshAccessToken: failed verify refresh token: %w", err)
	}

	userLogin, ok := claims["login"].(string)
	if !ok {
		return "", fmt.Errorf("RefreshAccessToken: failed get login from claims: %w", err)
	}

	return s.CreateAccessToken(userLogin)
}
