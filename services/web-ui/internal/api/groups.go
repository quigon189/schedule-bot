package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
	"web-ui/internal/models"
)

func (c *CoreClient) GetGroups(ctx context.Context, session *Session) ([]models.Group, error) {
	var groups []models.Group
	req := &request{
		method: "GET",
		path:   "/groups",
	}

	if err := c.doWithAuth(ctx, session, req, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

type CreateGroupRequest struct {
	Name          string `json:"name"`
	Specialty     string `json:"specialty"`
	AdmissionYear int    `json:"admission_year"`
}

type FileResponse struct {
	Success bool `json:"success"`
	Data    struct {
		FilePath  string    `json:"file_path"`
		ExpiresAt time.Time `json:"expires"`
	} `json:"data"`
}

func (c *CoreClient) CreateGroup(ctx context.Context, session *Session, req CreateGroupRequest) (*models.Group, error) {
	var group models.Group
	apiReq := &request{
		method: "POST",
		path:   "/groups",
		body:   req,
	}
	if err := c.doWithAuth(ctx, session, apiReq, &group); err != nil {
		return nil, err
	}
	return &group, nil
}

type UpdateGroupRequest struct {
	Name          *string `json:"name,omitempty"`
	Specialty     *string `json:"specialty,omitempty"`
	AdmissionYear *int    `json:"admission_year,omitempty"`
}

func (c *CoreClient) UpdateGroup(ctx context.Context, s *Session, id int, req UpdateGroupRequest) (*models.Group, error) {
	var group models.Group
	apiReq := &request{method: "PATCH", path: fmt.Sprintf("/groups/%d", id), body: req}
	if err := c.doWithAuth(ctx, s, apiReq, &group); err != nil {
		return nil, err
	}
	return &group, nil
}

func (c *CoreClient) DeleteGroup(ctx context.Context, s *Session, id int) error {
	req := &request{method: "DELETE", path: fmt.Sprintf("/groups/%d", id)}
	return c.doWithAuth(ctx, s, req, nil)
}

func (c *CoreClient) GetGroup(ctx context.Context, s *Session, id int) (*models.Group, error) {
	var group models.Group
	req := &request{method: "GET", path: fmt.Sprintf("/groups/%d", id)}
	if err := c.doWithAuth(ctx, s, req, &group); err != nil {
		return nil, err
	}
	return &group, nil
}

func (c *CoreClient) DownloadGroupTemplate(ctx context.Context, s *Session) ([]byte, error) {
	raw, _, err := c.doRawWithAuth(ctx, s, &request{method: "GET", path: "/groups/template"})
	var data FileResponse
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}
	file, _, err := c.doRawWithAuth(ctx, s, &request{method: "GET", path: data.Data.FilePath})
	return file, err
}

type File struct {
	Filename string
	FileData []byte
}

func (c *CoreClient) DownloadAIGroupTemplate(ctx context.Context, s *Session, message string, files []File) ([]byte, error) {
	var result struct {
		FilePath  string    `json:"file_path"`
		ExpiresAt time.Time `json:"expires_at"`
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	// Добавляем файлы
	for _, f := range files {
		part, err := writer.CreateFormFile("files", f.Filename)
		if err != nil {
			return nil, fmt.Errorf("create form file: %w", err)
		}
		if _, err := part.Write(f.FileData); err != nil {
			return nil, fmt.Errorf("write file data: %w", err)
		}
	}

	if err := writer.WriteField("message", message); err != nil {
		return nil, fmt.Errorf("write field message: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close writer: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/groups/template/ai", body)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	httpReq.Header.Set("Authorization", "Bearer "+s.AccessToken)

	client := http.Client{Timeout: 20 * time.Minute}

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	var res models.Result
	if err := json.Unmarshal(respBody, &res); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, res.Error)
	}
	if !res.Success {
		return nil, fmt.Errorf("API error: %s", res.Error)
	}
	if err := json.Unmarshal(res.Data, &result); err != nil{
		return nil, fmt.Errorf("unmarshal result: %w", err)
	}

	file, _, err := c.doRawWithAuth(ctx, s, &request{method: "GET", path: result.FilePath})
	return file, err
}

func (c *CoreClient) UploadGroupExcel(ctx context.Context, s *Session, fileData []byte, filename string) (*models.GroupWithCurriculumResponse, error) {
	var result models.GroupWithCurriculumResponse
	err := c.doMultipartWithAuth(ctx, s, "/groups/upload", fileData, filename, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
