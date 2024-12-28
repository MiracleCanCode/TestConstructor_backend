package validation

import (
	"github.com/go-playground/validator/v10"
)

func Validation(payload interface{}) error {
	validate := validator.New()
	return validate.Struct(payload)
}
