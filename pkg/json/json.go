package json

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/server/pkg/validation"
	"go.uber.org/zap"
)

type DecodeAndEncodeJson struct {
	r      *http.Request
	logger *zap.Logger
	w      http.ResponseWriter
}

func New(r *http.Request, log *zap.Logger, w http.ResponseWriter) *DecodeAndEncodeJson {
	return &DecodeAndEncodeJson{
		r:      r,
		logger: log,
		w:      w,
	}
}

func (s *DecodeAndEncodeJson) Decode(payload any) error {
	if err := json.NewDecoder(s.r.Body).Decode(&payload); err != nil {
		return fmt.Errorf("Decode: failed decode request body: %w", err)
	}

	return nil
}

func (s *DecodeAndEncodeJson) Encode(code int, data any) error {
	s.w.Header().Set("Content-Type", "application/json")
	s.w.WriteHeader(code)
	if err := json.NewEncoder(s.w).Encode(data); err != nil {
		return fmt.Errorf("Encode: failed to encode response body: %w", err)
	}

	return nil
}

func (s *DecodeAndEncodeJson) DecodeAndValidationBody(payload any) error {
	if err := s.Decode(payload); err != nil {
		return fmt.Errorf("DecodeAndValidationBody: failed to decode request body: %w", err)
	}

	if err := validation.Validation(payload); err != nil {
		return fmt.Errorf("DecodeAndValidationBody: failed to validation request body: %w", err)
	}

	return nil
}

func (s *DecodeAndEncodeJson) Marshall(payload any) ([]byte, error) {
	s.w.Header().Set("Content-Type", "application/json")
	result, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Marshall: failed to json marshall: %w", err)
	}

	return result, nil
}
