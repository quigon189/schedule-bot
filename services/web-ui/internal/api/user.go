package api

import (
	"encoding/json"
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
		query: query,
	}

	if err := c.Do(w, r, req, &paginatedUsers); err != nil {
		return nil, err
	}

	return &paginatedUsers, nil
}
