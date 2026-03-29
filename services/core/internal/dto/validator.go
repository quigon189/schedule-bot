package dto

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s any) error {
	return validate.Struct(s)
}
