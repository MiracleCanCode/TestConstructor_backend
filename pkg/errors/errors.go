package errors

import (
	"net/http"

	"github.com/server/pkg/constants"
	errorconstant "github.com/server/pkg/errorConstants"
	mapjson "github.com/server/pkg/mapJson"
	"go.uber.org/zap"
)

type ErrorHandler interface {
	HandleError(w http.ResponseWriter, r *http.Request, err error)
}

type DefaultErrorHandler struct {
	logger *zap.Logger
}

func NewErrorHandler(logger *zap.Logger) *DefaultErrorHandler {
	return &DefaultErrorHandler{logger: logger}
}

func (h *DefaultErrorHandler) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	var statusCode int
	var message string

	switch err {
	case errorconstant.ErrInvalidCredentials:
		statusCode = http.StatusUnauthorized
		message = constants.LoginOrPasswordIncorrect
	case errorconstant.ErrLoginAlreadyTaken:
		statusCode = http.StatusConflict
		message = constants.LoginIsExist
	default:
		statusCode = http.StatusInternalServerError
		message = constants.InternalServerError
	}

	h.logger.Warn(message, zap.Error(err))

	jsonResponse := mapjson.New(h.logger, w, r)
	jsonResponse.JsonError(message)

	w.WriteHeader(statusCode)
}

func HandleError(logger *zap.Logger, w http.ResponseWriter, r *http.Request, err error, msg string, code string) {
	jsonError := mapjson.New(logger, w, r)
	logger.Error(msg, zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
	jsonError.JsonError(code)
}
