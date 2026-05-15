package api

import (
	"fmt"
	"net/http"
	"web-ui/internal/models"
)

func (c *CoreClient) GetSubjects(w http.ResponseWriter, r *http.Request, groupID int) ([]models.Subject, error) {
	req := request{
		method: "GET",
		path: fmt.Sprintf("/subjects/group/%d", groupID),
	}

	var subjects []models.Subject
	if err := c.Do(w, r, req, &subjects); err != nil {
		return nil, fmt.Errorf("failed to get subjects: %w", err)
	}

	return subjects, nil
} 
