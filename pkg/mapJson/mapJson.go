package mapjson

import (
	"net/http"

	"github.com/server/pkg/json"
	"go.uber.org/zap"
)

type MessageResponse struct {
	logger *zap.Logger
	w http.ResponseWriter
	r *http.Request
}

func New(logger *zap.Logger,
	w http.ResponseWriter,
	r *http.Request) *MessageResponse {
		return &MessageResponse{
			w: w,
			logger: logger,
			r: r,
		}
	}

func (s *MessageResponse) JsonError(message string) {
	json := json.New(s.r, s.logger, s.w)
	s.w.WriteHeader(http.StatusBadRequest)

	jsonMap := map[string]string{
		"error": message,
	}

	dataMarshaled, err := json.Marshall(jsonMap)
	if err != nil {
		s.logger.Error("Failed to marshal error response", zap.Error(err))
		http.Error(s.w, "Internal server error", http.StatusInternalServerError)
		return
	}

	s.w.Header().Set("Content-Type", "application/json")
	s.w.Write(dataMarshaled)
}

func (s *MessageResponse) JsonSuccess(message string) {
	json := json.New(s.r, s.logger, s.w)
	s.w.WriteHeader(http.StatusAccepted)

	jsonMap := map[string]string{
		"success": message,
	}

	dataMarshaled, err := json.Marshall(jsonMap)
	if err != nil {
		s.logger.Error("Failed to marshal success response", zap.Error(err))
		http.Error(s.w, "Internal server error", http.StatusInternalServerError)
		return
	}

	s.w.Header().Set("Content-Type", "application/json")
	s.w.Write(dataMarshaled)
}