package middleware

import (
	"net/http"

	"github.com/MiracleCanCode/example_configuration_logger/pkg/logger"
	"github.com/server/configs"
	"github.com/server/internal/repository"
	"github.com/server/pkg/cookie"
	"github.com/server/pkg/db/postgresql"
	"github.com/server/pkg/jwt"
	"go.uber.org/zap"
)

func IsAuth(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.Logger(logger.DefaultLoggerConfig())
		JWT := jwt.NewJwt(log)

		// Загрузка конфигурации
		cfg, err := configs.Load(log)
		if err != nil {
			log.Error("Failed to load config", zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Инициализация базы данных
		db, err := postgresql.New(cfg, log)
		if err != nil {
			log.Error("Failed to create DB instance", zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		userRepo := repository.NewUser(db, log)
		cookies := cookie.New(w, r, log)

		// Извлечение пользователя из токена
		login, err := JWT.ExtractUserFromCookie(r, "token")
		if err != nil {
			log.Error("Failed to extract user from token", zap.Error(err))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Поиск пользователя по логину
		findUserByLogin, err := userRepo.GetUserByLogin(login)
		if err != nil {
			log.Error("Failed to find user by login", zap.Error(err))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Проверка на наличие refresh token
		if findUserByLogin.RefreshToken == "" {
			log.Warn("User does not have a refresh token", zap.String("login", login))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Обновление access token
		accessToken, err := JWT.RefreshAccessToken(findUserByLogin.RefreshToken)
		if err != nil {
			log.Error("Failed to refresh access token", zap.Error(err))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Установка нового токена в cookie
		cookies.Set("token", accessToken)

		// Продолжение обработки запроса
		next.ServeHTTP(w, r)
	}
}
