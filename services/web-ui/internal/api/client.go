package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"web-ui/internal/models"
)

type CoreClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewCoreClient(baseURL string, timeout time.Duration) *CoreClient {
	return &CoreClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: timeout},
	}
}

func (c *CoreClient) refreshToken(ctx context.Context, refresh, sessionID string) {}

func (c *CoreClient) Login(ctx context.Context, username, password string) (*models.JWT, error) {
	body := map[string]string{"username": username, "password": password}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequestWithContext(ctx, "POST", c.baseURL + "/login", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w",err)
	}
	defer resp.Body.Close()

	var result models.Result
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("deocode body: %w",err)
	}

	if !result.Success {
		return nil, fmt.Errorf(result.Error)
	}

	var jwt models.JWT
	if err := json.Unmarshal(result.Data, &jwt); err != nil {
		return nil, fmt.Errorf("unmarshal data: %w",err)
	}

	return &jwt, nil
}
