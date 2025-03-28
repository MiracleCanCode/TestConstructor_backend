package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/server/pkg/logger"
	"go.uber.org/zap"
)

func TraceLogger(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		log := logger.GetInstance()

		log.Info("request", zap.Any("request_data", r.Body), zap.String("request_id", requestID))
		next.ServeHTTP(w, r)
	}
}
