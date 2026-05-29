package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"web-ui/internal/models"
)

func (c *CoreClient) GetCurrentUser(ctx context.Context, s *Session) (*models.User, error) {
	var user models.User
	req := &request{
		method: "GET",
		path:   "/user",
	}

	if err := c.doWithAuth(ctx, s, req, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *CoreClient) GetPaginatedUsers(ctx context.Context, s *Session, params models.PaginatedUsersQuery) (*models.PaginatedUsers, error) {
	var paginatedUsers models.PaginatedUsers
	var query map[string]string
	data, _ := json.Marshal(params)
	if err := json.Unmarshal(data, &query); err != nil {
		return nil, err
	}
	req := &request{
		method: "GET",
		path:   "/admin/users/paginated",
		query:  query,
	}

	if err := c.doWithAuth(ctx, s, req, &paginatedUsers); err != nil {
		return nil, err
	}

	return &paginatedUsers, nil
}

func (c *CoreClient) CreateUser(ctx context.Context, s *Session, req models.CreateUserRequest) error {
	var reqBody = make(map[string]any)
	reqBody["username"] = req.Username
	reqBody["full_name"] = req.FullName
	reqBody["email"] = req.Email
	reqBody["password"] = req.Password

	switch req.Role {
	case "student":
		if req.GroupID == nil {
			return errors.New("group is required")
		}
		reqBody["group_id"] = req.GroupID

		apiReq := &request{
			method: "POST",
			path:   "/students",
			body:   reqBody,
		}

		return c.doWithAuth(ctx, s, apiReq, nil)
	case "teacher":
		apiReq := &request{
			method: "POST",
			path:   "/teachers",
			body:   reqBody,
		}

		return c.doWithAuth(ctx, s, apiReq, nil)
	case "user":
		apiReq := &request{
			method: "/POST",
			path:   "/admin/users",
			body:   reqBody,
		}

		return c.doWithAuth(ctx, s, apiReq, nil)
	}
	return errors.New("role is not available")
}

func (c *CoreClient) ChangeUserPasswordAdmin(ctx context.Context, s *Session, userID int, new string) error {
	reqBody := make(map[string]any)
	reqBody["user_id"] = userID
	reqBody["new_password"] = new

	req := &request{
		method: "POST",
		path:   "/admin/users/password",
		body:   reqBody,
	}

	return c.doWithAuth(ctx, s, req, nil)
}

func (c *CoreClient) DeleteUser(ctx context.Context, s *Session, userID int) error {
	return errors.New("method not implemented")
}

func (c *CoreClient) GetUserByID(ctx context.Context, s *Session, userID int) (*models.User, error) {
	var user models.User

	req := &request{
		method: "GET",
		path: fmt.Sprintf("/admin/users/%d", userID),
	}

	if err := c.doWithAuth(ctx, s, req, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *CoreClient) DownloadTeacherTemplate(ctx context.Context, s *Session) ([]byte, error) {
	data, _, err := c.doRawWithAuth(ctx, s, &request{method: "GET", path: "/teachers/template"})
	return data, err
}

func (c *CoreClient) UploadTeachersExcel(ctx context.Context, s *Session, fileData []byte, filename string) ([]models.TeacherCreationResult, error) {
	var result []models.TeacherCreationResult
	err := c.doMultipartWithAuth(ctx, s, "/teachers/upload", fileData, filename, nil, &result)
	return result, err
}
