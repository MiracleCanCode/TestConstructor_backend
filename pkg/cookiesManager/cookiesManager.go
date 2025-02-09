package cookiesmanager

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type CookiesManagerInterface interface {
	Get(name string) (string, error)
	Set(name string, value string, maxAge time.Duration, httpOnly bool, w http.ResponseWriter)
	Delete(name string, w http.ResponseWriter)
}

type CookiesManager struct {
	r      *http.Request
	logger *zap.Logger
}

func New(r *http.Request, logger *zap.Logger) *CookiesManager {
	return &CookiesManager{
		r:      r,
		logger: logger,
	}
}

func (s *CookiesManager) Get(name string) (string, error) {
	value, err := s.r.Cookie(name)
	if err != nil {
		s.logger.Error("Extract cookie", zap.Error(err))
		return "", err
	}

	return value.Value, nil
}

func (s *CookiesManager) Set(name string, value string, maxAge time.Duration, httpOnly bool, w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(maxAge),
		HttpOnly: httpOnly,
		Secure:   false,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
}

func (s *CookiesManager) Delete(name string, w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    name,
		Value:   "",
		Expires: time.Unix(0, 0),
		Path:    "/",
	})
}
