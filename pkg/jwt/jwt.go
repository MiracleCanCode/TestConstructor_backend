package jwt

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/server/configs"
)

type JWT struct {
	Secret string
}

func NewJwt(secret string) *JWT {
	return &JWT{
		Secret: secret,
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
		return nil, nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, nil, jwt.ErrInvalidKey
	}

	return token, claims, nil
}

func ExtractUserFromAuthHeader(r *http.Request, cfg *configs.Config) (string, error) {
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	_, claims, err := NewJwt(cfg.SECRET).VerifyToken(tokenString)
	if err != nil {
		return "", err
	}

	userLogin, ok := claims["login"].(string)
	if !ok {
		return "", errors.New("invalid login claim type")
	}

	return userLogin, nil
}
