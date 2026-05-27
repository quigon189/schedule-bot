package api

import (
	"context"
	"web-ui/internal/models"
)

func (c *CoreClient) GetGroups(ctx context.Context, session *Session) ([]models.Group, error) {
	var groups []models.Group
	req := &request{
		method: "GET",
		path: "/groups",
	}	

	if err := c.doWithAuth(ctx, session, req, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}
