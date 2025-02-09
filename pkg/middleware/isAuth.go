package middleware

import (
	"net/http"
	"time"

	logger "github.com/MiracleCanCode/example_configuration_logger"
	"github.com/server/configs"
	"github.com/server/internal/repository"
	cookiesmanager "github.com/server/pkg/cookiesManager"
	"github.com/server/pkg/db/postgresql"
	"github.com/server/pkg/jwt"
	"go.uber.org/zap"
)

func IsAuth(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.Logger(logger.DefaultLoggerConfig())
		JWT := jwt.NewJwt(log)
		cookie := cookiesmanager.New(r, log)

		cfg, err := configs.Load(log)
		if err != nil {
			log.Error("Failed to load config", zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		db, err := postgresql.New(cfg, log)
		if err != nil {
			log.Error("Failed to create DB instance", zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		userRepo := repository.NewUser(db, log)

		login, err := JWT.ExtractUserFromToken(r)
		if err != nil {
			log.Error("Failed to extract user from token", zap.Error(err))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		findUserByLogin, err := userRepo.GetUserByLogin(login)
		if err != nil {
			log.Error("Failed to find user by login", zap.Error(err))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if findUserByLogin.RefreshToken == "" {
			log.Warn("User does not have a refresh token", zap.String("login", login))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		accessToken, err := JWT.RefreshAccessToken(findUserByLogin.RefreshToken)
		if err != nil {
			log.Error("Failed to refresh access token", zap.Error(err))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		cookie.Set("token", accessToken, time.Minute*15, true, w)

		next.ServeHTTP(w, r)
	}
}
