package api

import (
	"context"
	"fmt"
	"web-ui/internal/models"
)

type CreateAcademicPeriodRequest struct {
	Year      string `json:"year"`
	Semester  int    `json:"semester"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (c *CoreClient) CreateAcademicPeriod(ctx context.Context, s *Session, req CreateAcademicPeriodRequest) (*models.AcademicPeriod, error) {
	var period models.AcademicPeriod
	apiReq := &request{method: "POST", path: "/academic-periods", body: req}
	if err := c.doWithAuth(ctx, s, apiReq, &period); err != nil {
		return nil, err
	}
	return &period, nil
}

type UpdateAcademicPeriodRequest struct {
	Year      *string `json:"year,omitempty"`
	Semester  *int    `json:"semester,omitempty"`
	StartDate *string `json:"start_date,omitempty"`
	EndDate   *string `json:"end_date,omitempty"`
}

func (c *CoreClient) UpdateAcademicPeriod(ctx context.Context, s *Session, id int, req UpdateAcademicPeriodRequest) (*models.AcademicPeriod, error) {
	var period models.AcademicPeriod
	apiReq := &request{method: "PATCH", path: fmt.Sprintf("/academic-periods/%d", id), body: req}
	if err := c.doWithAuth(ctx, s, apiReq, &period); err != nil {
		return nil, err
	}
	return &period, nil
}

func (c *CoreClient) DeleteAcademicPeriod(ctx context.Context, s *Session, id int) error {
	req := &request{method: "DELETE", path: fmt.Sprintf("/academic-periods/%d", id)}
	return c.doWithAuth(ctx, s, req, nil)
}

func (c *CoreClient) GetAcademicPeriod(ctx context.Context, s *Session, id int) (*models.AcademicPeriod, error) {
	var period models.AcademicPeriod
	req := &request{method: "GET", path: fmt.Sprintf("/academic-periods/%d", id)}
	if err := c.doWithAuth(ctx, s, req, &period); err != nil {
		return nil, err
	}
	return &period, nil
}

func (c *CoreClient) GetAcademicPeriods(ctx context.Context, s *Session) ([]models.AcademicPeriod, error) {
	var periods []models.AcademicPeriod
	req := &request{method: "GET", path: "/academic-periods"}
	if err := c.doWithAuth(ctx, s, req, &periods); err != nil {
		return nil, err
	}
	return periods, nil
}
