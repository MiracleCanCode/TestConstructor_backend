package middleware

import (
	"net/http"

	"github.com/MiracleCanCode/example_configuration_logger/pkg/logger"
	"github.com/server/pkg/cookie"
	"go.uber.org/zap"
)

func IsAuth(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.Logger(logger.DefaultLoggerConfig())
		cookies := cookie.New(w, r, log)

		token := cookies.Get("token")

		if token == "" {
			log.Warn("Unauthorized access attempt", zap.String("path", r.URL.Path))
			http.Redirect(w, r, "/auth", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	}
}
