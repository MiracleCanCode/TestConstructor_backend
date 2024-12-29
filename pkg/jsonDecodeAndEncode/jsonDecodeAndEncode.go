package jsonDecodeAndEncode

import (
	"encoding/json"
	"net/http"

	"github.com/server/pkg/validation"
	"go.uber.org/zap"
)

type DecodeAndEncodeJson struct {
	r      *http.Request
	logger *zap.Logger
	w      http.ResponseWriter
}

func NewDecodeAndEncodeJson(r *http.Request, log *zap.Logger, w http.ResponseWriter) *DecodeAndEncodeJson {
	return &DecodeAndEncodeJson{
		r:      r,
		logger: log,
		w:      w,
	}
}

func (s *DecodeAndEncodeJson) Decode(payload any) error {
	if err := json.NewDecoder(s.r.Body).Decode(&payload); err != nil {
		s.logger.Error("Failed to decode login data", zap.Error(err))
		http.Error(s.w, "Invalid request body", http.StatusBadRequest)
		return err
	}

	return nil
}

func (s *DecodeAndEncodeJson) Encode(code int, data any) error {
	s.w.Header().Set("Content-Type", "application/json")
	s.w.WriteHeader(code)
	if err := json.NewEncoder(s.w).Encode(data); err != nil {
		http.Error(s.w, "Failed to encode response", http.StatusInternalServerError)
		return err
	}

	return nil
}

func (s *DecodeAndEncodeJson) DecodeAndValidationBody(payload any) error {
	if err := s.Decode(payload); err != nil {
		return err
	}

	if err := validation.Validation(payload); err != nil {
		s.logger.Error("Validation failed", zap.Error(err))
		http.Error(s.w, "Validation error: "+err.Error(), http.StatusBadRequest)
		return err
	}

	return nil
}

func (s *DecodeAndEncodeJson) Marshall(payload any) ([]byte, error) {
	s.w.Header().Set("Content-Type", "application/json")
	return json.Marshal(payload)
}
