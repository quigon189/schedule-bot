package dto

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	validate := validator.New()

	validate.RegisterValidation("academic_year", isAcademicYear)
	validate.RegisterValidation("group_name", isGroupName)

	return validate
}

func isAcademicYear(fl validator.FieldLevel) bool {
	ok, _ := regexp.MatchString("^(\\d{4})/(\\d{4})", fl.Field().String())
	return ok
}

func isGroupName(fl validator.FieldLevel) bool {
	ok, _ := regexp.MatchString("^[А-Я]+-\\d+", fl.Field().String())
	return ok
}
