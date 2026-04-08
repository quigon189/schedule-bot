package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"schedule-cli/internal/config"
	"schedule-cli/internal/models"
	"time"
)

type APIClient struct {
	httpClient *http.Client
	cfgMgr     *config.Manager
}

func NewAPIClient(mgr *config.Manager) *APIClient {
	cfg := mgr.Get()
	return &APIClient{
		cfgMgr:     mgr,
		httpClient: &http.Client{Timeout: time.Duration(cfg.Timeout) * time.Second},
	}
}

func (c *APIClient) doRequest(method, path string, headers map[string]string, body any, result any) error {
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal body: %w", err)
		}
	}

	cfg := c.cfgMgr.Get()
	url := cfg.BaseURL + path

	req, err := http.NewRequest(method, url, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	var apiResp models.Response
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("decode response body: %w", err)
	}

	if !apiResp.Success {
		return fmt.Errorf("API error: %s", apiResp.Error)
	}

	if result != nil && apiResp.Data != nil {
		dataBytes, _ := json.Marshal(apiResp.Data)
		return json.Unmarshal(dataBytes, result)
	}

	return nil
}

func (c *APIClient) Refresh() error {
	cfg := c.cfgMgr.Get()
	req := models.RefreshRequest{
		RefreshToken: cfg.RefreshToken,
		SessionID:    cfg.SessionID,
	}
	var resp models.LoginResponse
	if err := c.doRequest("POST", "/refresh", map[string]string{}, req, &resp); err != nil {
		return err
	}
	return c.cfgMgr.UpdateToken(config.UpdateToken{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		SessionID:    cfg.SessionID,
		ExpiresAt:    resp.ExiresAt,
	})
}

func (c *APIClient) Login(req models.LoginRequest) error {
	var resp models.LoginResponse
	headers := map[string]string{
		"User-Agent": "schedule-cli",
	}
	if err := c.doRequest("POST", "/login", headers, req, &resp); err != nil {
		return fmt.Errorf("request: %w", err)
	}

	return c.cfgMgr.UpdateToken(config.UpdateToken{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		SessionID:    resp.SessionID,
		ExpiresAt:    resp.ExiresAt,
	})
}

func (c *APIClient) Get(path string, query map[string]string, result any) error {
	cfg := c.cfgMgr.Get()
	expiresAt := time.Unix(cfg.ExpiresAt, 0)
	if time.Now().After(expiresAt) {
		if cfg.RefreshToken == "" {
			return  fmt.Errorf("you must login")
		}
		if err := c.Refresh(); err != nil {
			return fmt.Errorf("refresh token: %w", err)
		}
		cfg = c.cfgMgr.Get()
	}

	headrs := map[string]string{
		"Authorization": "Bearer "+cfg.AccessToken,
	}

	u, err := url.Parse(path)
	if err != nil {
		return fmt.Errorf("parse url path")
	}

	q := u.Query()
	for key, value := range query {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	return c.doRequest("GET", u.String(), headrs, nil, result)
}
