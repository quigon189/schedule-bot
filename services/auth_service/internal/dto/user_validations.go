package dto

import "fmt"

func (r *CreateUserRequest) Validate() error {
	if r.FullName == "" {
		return fmt.Errorf("invalid full name")
	}

	if r.TelegramID == 0 {
		return fmt.Errorf("invalid telegram id")
	}

	return nil
}
