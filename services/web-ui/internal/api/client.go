package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"web-ui/internal/models"
	"web-ui/internal/session"
)

type CoreClient struct {
	baseURL        string
	httpClient     *http.Client
	sessionManager *session.SessionManager
}

func NewCoreClient(baseURL string, timeout time.Duration, sm *session.SessionManager) *CoreClient {
	return &CoreClient{
		baseURL:        baseURL,
		httpClient:     &http.Client{Timeout: timeout},
		sessionManager: sm,
	}
}

func (c *CoreClient) doRequest(ctx context.Context, method, path string, headers map[string]string, body any, result any) error {
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bytes.NewReader(reqBody))
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

	var res models.Result
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("decode response body: %w", err)
	}

	if !res.Success {
		return fmt.Errorf("API error: %w", res.Error)
	}

	if result != nil && res.Data != nil {
		return json.Unmarshal(res.Data, result)
	}

	return nil
}

func (c *CoreClient) refreshToken(ctx context.Context, jwt *models.JWT) error {
	body := map[string]string{"refresh_token": jwt.Refresh, "session_id": jwt.SessionID}
	headers := map[string]string{}
	return c.doRequest(ctx, "POST", "/refresh", headers, body, jwt)
}

func (c *CoreClient) getAccessToken(w http.ResponseWriter, r *http.Request) (string, error) {
	jwt, ok := c.sessionManager.Get(r, "jwt").(models.JWT)
	if !ok {
		return "", fmt.Errorf("failed to get jwt from session")
	}
	exiresAt :=  time.Unix(jwt.ExpiresAt, 0)
	if time.Now().After(exiresAt) {
		if err := c.refreshToken(r.Context(), &jwt); err != nil {
			return "", err
		}
		c.sessionManager.Set(w, r, "jwt", jwt)
	}
	return jwt.AccessToken, nil
}

func (c *CoreClient) Login(r *http.Request, username, password string) (*models.JWT, error) {
	body := map[string]string{"username": username, "password": password}
	var jwt models.JWT
	var clientIP string
	if xForwardedFor := r.Header.Get("X-Forwarded-For"); xForwardedFor != "" {
		clientIP = xForwardedFor
	} else {
		ips := strings.Split(r.RemoteAddr, ",")
		clientIP = strings.TrimSpace(ips[0])
	}
	headers := map[string]string{
		"User-Agent": r.Header.Get("User-Agent"),
		"clientIP": clientIP,
	}
	if err := c.doRequest(r.Context(), "POST", "/login", headers, body, &jwt); err != nil {
		return nil, err
	}

	return &jwt, nil
}


