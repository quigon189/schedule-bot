package dto

import (
	"fmt"
	"regexp"
	"slices"

	"auth_service/internal/models"
)

func (r *CreateRegistrationCodeRequest) Validate() error {
	if !slices.Contains([]string{
		models.ManagerRole,
		models.StudentRole,
		models.TeacherRole,
	}, r.RoleName) {
		return fmt.Errorf("invalid role name")
	}

	if r.RoleName == models.StudentRole {
		if r.GroupName == nil {
			return fmt.Errorf("invalid group name")
		}
		ok, err := regexp.Match("^[А-Я]+-[0-9]+", []byte(*r.GroupName))
		if !ok || err != nil {
			return fmt.Errorf("invalid group name")
		}
	}

	if r.MaxUses == 0 {
		r.MaxUses = 1
	}

	if r.MaxUses > 50 {
		return fmt.Errorf("invalid mac uses")
	}

	if r.Expiration == 0 {
		r.Expiration = 21_600
	}

	return nil
}
