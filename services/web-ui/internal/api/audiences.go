package api

import (
	"context"
	"fmt"
	"web-ui/internal/models"
)

type CreateAudienceRequest struct {
	Name   string `json:"name"`
	Number string `json:"number"`
}

func (c *CoreClient) CreateAudience(ctx context.Context, s *Session, req CreateAudienceRequest) (*models.Audience, error) {
	var audience models.Audience
	apiReq := &request{method: "POST", path: "/audiences", body: req}
	if err := c.doWithAuth(ctx, s, apiReq, &audience); err != nil {
		return nil, err
	}
	return &audience, nil
}

type UpdateAudienceRequest struct {
	Name   *string `json:"name,omitempty"`
	Number *string `json:"number,omitempty"`
}

func (c *CoreClient) UpdateAudience(ctx context.Context, s *Session, id int, req UpdateAudienceRequest) (*models.Audience, error) {
	var audience models.Audience
	apiReq := &request{method: "PATCH", path: fmt.Sprintf("/audiences/%d", id), body: req}
	if err := c.doWithAuth(ctx, s, apiReq, &audience); err != nil {
		return nil, err
	}
	return &audience, nil
}

func (c *CoreClient) DeleteAudience(ctx context.Context, s *Session, id int) error {
	req := &request{method: "DELETE", path: fmt.Sprintf("/audiences/%d", id)}
	return c.doWithAuth(ctx, s, req, nil)
}

func (c *CoreClient) GetAudience(ctx context.Context, s *Session, id int) (*models.Audience, error) {
	var audience models.Audience
	req := &request{method: "GET", path: fmt.Sprintf("/audiences/%d", id)}
	if err := c.doWithAuth(ctx, s, req, &audience); err != nil {
		return nil, err
	}
	return &audience, nil
}

func (c *CoreClient) GetAudiences(ctx context.Context, s *Session) ([]models.Audience, error) {
	var audiences []models.Audience
	req := &request{method: "GET", path: "/audiences"}
	if err := c.doWithAuth(ctx, s, req, &audiences); err != nil {
		return nil, err
	}
	return audiences, nil
}

func (c *CoreClient) DownloadAudienceTemplate(ctx context.Context, s *Session) ([]byte, error) {
	data, _, err := c.doRawWithAuth(ctx, s, &request{method: "GET", path: "/audiences/template"})
	return data, err
}

func (c *CoreClient) UploadAudiencesExcel(ctx context.Context, s *Session, fileData []byte, filename string) ([]models.Audience, error) {
	var result []models.Audience
	err := c.doMultipartWithAuth(ctx, s, "/audiences/upload", fileData, filename, nil, &result)
	return result, err
}
