package api

import (
	"fmt"
	"net/http"
	"web-ui/internal/models"
)

func (c *CoreClient) GetCurrentUser(w http.ResponseWriter, r *http.Request) (*models.User, error) {
	var user models.User
	req := request{
		method: "GET",
		path:   "/user",
	}

	if err := c.Do(w, r, req, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *CoreClient) GetStudent(w http.ResponseWriter, r *http.Request, user *models.User) (*models.Student, error) {
	var student models.Student
	req := request{
		method: "GET",
		path:   fmt.Sprintf("/students/%d", user.ID),
	}

	if err := c.Do(w, r, req, &student); err != nil {
		return nil, err
	}

	return &student, nil
}
