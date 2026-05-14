package api

import (
	"encoding/json"
	"errors"
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

func (c *CoreClient) GetPaginatedUsers(w http.ResponseWriter, r *http.Request, params models.PaginatedUsersQuery) (*models.PaginatedUsers, error) {
	var paginatedUsers models.PaginatedUsers
	var query map[string]string
	data, _ := json.Marshal(params)
	if err := json.Unmarshal(data, &query); err != nil {
		return nil, err
	}
	req := request{
		method: "GET",
		path:   "/admin/users/paginated",
		query:  query,
	}

	if err := c.Do(w, r, req, &paginatedUsers); err != nil {
		return nil, err
	}

	return &paginatedUsers, nil
}

func (c *CoreClient) CreateUser(w http.ResponseWriter, r *http.Request, req models.CreateUserRequest) error {
	var reqBody = make(map[string]any)
	reqBody["username"] = req.Username
	reqBody["full_name"] = req.FullName
	reqBody["email"] = req.Email
	reqBody["password"] = req.Password

	switch req.Role {
	case "student":
		if req.Group == nil {
			return errors.New("group is required")
		}
		reqBody["group"] = req.Group

		apiReq := request{
			method: "POST",
			path:   "/students",
			body:   reqBody,
		}

		return c.Do(w, r, apiReq, nil)
	case "teacher":
		apiReq := request{
			method: "POST",
			path:   "/teachers",
			body:   reqBody,
		}

		return c.Do(w, r, apiReq, nil)
	case "user":
		apiReq := request{
			method: "/POST",
			path:   "/admin/users",
			body:   reqBody,
		}

		return c.Do(w, r, apiReq, nil)
	}
	return errors.New("role is not available")
}

func (c *CoreClient) ChangeUserPasswordAdmin(w http.ResponseWriter, r *http.Request, userID int, new string) error {
	reqBody := make(map[string]any)
	reqBody["user_id"] = userID
	reqBody["new_password"] = new

	req := request{
		method: "POST",
		path: "/admin/users/password",
		body: reqBody,
	}

	return c.Do(w, r, req, nil)
}

func (c *CoreClient) DeleteUser(w http.ResponseWriter, r *http.Request, userID int) error {
	return errors.New("method not implemented")
}


