package dto

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func NewScheduleValidator() *validator.Validate {
	validate := validator.New()

	validate.RegisterValidation("academic_year", isAcademicYear)
	validate.RegisterValidation("group_name", isGroupName)

	return validate
}

func NewChangeValidator() *validator.Validate {
	validate := validator.New()

	validate.RegisterValidation("date", isDate)

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

func isDate(fl validator.FieldLevel) bool {
	pattern := "^\\d{4}-([0][1-9]|[1][0-2])-([0][1-9]|[12][0-9]|3[01])$"
	ok, _ := regexp.MatchString(pattern, fl.Field().String())
	return ok
}
