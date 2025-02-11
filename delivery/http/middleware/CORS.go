package middleware

import (
	"net/http"

	"github.com/MiracleCanCode/example_configuration_logger"
	"github.com/rs/cors"
	"github.com/server/configs"
)

func CORSHandler(options cors.Options) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handler := cors.New(options).Handler(next)
		return handler
	}
}

func DefaultCORSMiddleware() func(http.Handler) http.Handler {
	log := logger.Logger(logger.DefaultLoggerConfig())
	cfg, err := configs.Load(log)
	if err != nil {
		log.Error("Failed to load config")
	}
	options := cors.Options{
		AllowedOrigins:   []string{cfg.CLIENT_URL, "http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Requested-With"},
		AllowCredentials: true,
	}

	return CORSHandler(options)
}
