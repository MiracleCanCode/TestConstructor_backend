package errors

import (
	mapjson "github.com/server/pkg/mapJson"
	"go.uber.org/zap"
	"net/http"
)

type ErrorHandler interface {
	HandleError(w http.ResponseWriter, r *http.Request, err error)
}

type DefaultErrorHandler struct {
	logger *zap.Logger
	w      http.ResponseWriter
	r      *http.Request
}

func New(logger *zap.Logger, w http.ResponseWriter, r *http.Request) *DefaultErrorHandler {
	return &DefaultErrorHandler{logger: logger, w: w, r: r}
}

func (s *DefaultErrorHandler) HandleError(message string, code int, err error) {
	s.logger.Error(message, zap.Error(err), zap.String("method", s.r.Method), zap.String("endpoint", s.r.URL.Path))

	jsonResponse := mapjson.New(s.logger, s.w, s.r)
	jsonResponse.JsonError(message)

	s.w.WriteHeader(code)
}
