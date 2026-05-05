package api

import (
	"context"
	"net/http"
	"web-ui/internal/models"
)

func (c *CoreClient) GetCurrentUser(ctx context.Context) (*models.User, error) {
	req, _ :=http.NewRequestWithContext(ctx, "GET", c.baseURL+"/user", nil)
	req.Header.Set("Authoriztion", "Bearer ")

	return nil, nil
}
