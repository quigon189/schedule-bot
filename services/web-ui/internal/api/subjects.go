package api

import (
	"context"
	"fmt"
	"web-ui/internal/models"
)

type CreateSubjectRequest struct {
	Title     string `json:"title"`
	Semester  int    `json:"semester"`
	HoursLoad int    `json:"hours_load"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	GroupID   int    `json:"group_id"`
}

func (c *CoreClient) CreateSubject(ctx context.Context, s *Session, req CreateSubjectRequest) (*models.Subject, error) {
	var subject models.Subject
	apiReq := &request{method: "POST", path: "/subjects", body: req}
	if err := c.doWithAuth(ctx, s, apiReq, &subject); err != nil {
		return nil, err
	}
	return &subject, nil
}

type UpdateSubjectRequest struct {
	Title     *string `json:"title,omitempty"`
	Semester  *int    `json:"semester,omitempty"`
	HoursLoad *int    `json:"hours_load,omitempty"`
	StartDate *string `json:"start_date,omitempty"`
	EndDate   *string `json:"end_date,omitempty"`
}

func (c *CoreClient) UpdateSubject(ctx context.Context, s *Session, id int, req UpdateSubjectRequest) (*models.Subject, error) {
	var subject models.Subject
	apiReq := &request{method: "PATCH", path: fmt.Sprintf("/subjects/%d", id), body: req}
	if err := c.doWithAuth(ctx, s, apiReq, &subject); err != nil {
		return nil, err
	}
	return &subject, nil
}

func (c *CoreClient) DeleteSubject(ctx context.Context, s *Session, id int) error {
	req := &request{method: "DELETE", path: fmt.Sprintf("/subjects/%d", id)}
	return c.doWithAuth(ctx, s, req, nil)
}

func (c *CoreClient) GetSubject(ctx context.Context, s *Session, id int) (*models.Subject, error) {
	var subject models.Subject
	req := &request{method: "GET", path: fmt.Sprintf("/subjects/%d", id)}
	if err := c.doWithAuth(ctx, s, req, &subject); err != nil {
		return nil, err
	}
	return &subject, nil
}

func (c *CoreClient) GetGroupSubjects(ctx context.Context, s *Session, groupID int) ([]models.Subject, error) {
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
