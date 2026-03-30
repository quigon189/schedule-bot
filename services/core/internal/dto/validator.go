package dto

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	return regexp.MustCompile("^[a-zA-Z0-9_]{3,20}$").MatchString(username)
}

func SetupValidator() {
	validate.RegisterValidation("username", validateUsername)
}

func ValidateStruct(s any) error {
	return validate.Struct(s)
}

