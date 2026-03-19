package dto

import (
	"core/internal/models"
	"errors"
	"net/mail"
	"regexp"
	"strings"
)

const (
	PASSWORD_LENGHT = 8
)

type CreateUserRequest struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *CreateUserRequest) Validate() error {
	var errs []string

	reUsername := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	if !reUsername.MatchString(r.Username) {
		errs = append(errs, "username error")
	}

	if len(r.FullName) < 3 {
		errs = append(errs, "full_name error")
	}

	if _, err := mail.ParseAddress(r.Email); err != nil {
		errs = append(errs, "email error")
	}

	if len(r.Password) < PASSWORD_LENGHT {
		errs = append(errs, "password error")
	}

	if len(errs) > 0 {
		errStr := strings.Join(errs, "; ")
		return errors.New(errStr)
	}

	return nil
}

type PaginatedUsers struct {
	Users      []models.User `json:"users"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PerPage    int           `json:"per_page"`
	TotalPages int           `json:"total_pages"`
}
