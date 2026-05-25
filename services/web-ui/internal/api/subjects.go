package api

import (
	"context"
	"fmt"
	"web-ui/internal/models"
)

func (c *CoreClient) GetSubjects(ctx context.Context, s *Session, groupID int) ([]models.Subject, error) {
	req := &request{
		method: "GET",
		path: fmt.Sprintf("/subjects/group/%d", groupID),
	}

	var subjects []models.Subject
	if err := c.doWithAuth(ctx, s, req, &subjects); err != nil {
		return nil, fmt.Errorf("failed to get subjects: %w", err)
	}

	return subjects, nil
} 
