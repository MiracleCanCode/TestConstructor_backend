package cookie

import (
	"errors"
	"net/http"

	"github.com/server/configs"
	"go.uber.org/zap"
)

type CookieInterface interface {
	Set(name string, value string)
	Get(name string) string
}

type Cookie struct {
	w   http.ResponseWriter
	r   *http.Request
	log *zap.Logger
}

func New(w http.ResponseWriter, r *http.Request, log *zap.Logger) *Cookie {
	return &Cookie{
		w:   w,
		r:   r,
		log: log,
	}
}

func (s *Cookie) Set(name string, value string) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   900,
		HttpOnly: false,
		Secure:   configs.PRODACTION,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(s.w, &cookie)
}

func (s *Cookie) Get(name string) string {
	cookie, err := s.r.Cookie(name)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			s.log.Error("Cookie not found")
		default:
			s.log.Error("Failed to get cookie")
		}

		return ""
	}

	return cookie.Value
}
