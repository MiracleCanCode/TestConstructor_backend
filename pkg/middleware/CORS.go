package middleware

import (
	"github.com/rs/cors"
	"net/http"
)

func CORSHandler(options cors.Options) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handler := cors.New(options).Handler(next)
		return handler
	}
}

func DefaultCORSMiddleware() func(http.Handler) http.Handler {
	options := cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Requested-With"},
		AllowCredentials: false,
	}

	return CORSHandler(options)
}
