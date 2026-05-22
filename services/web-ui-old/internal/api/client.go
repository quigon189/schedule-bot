package api

import ( "bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

type request struct {
	method  string
	path    string
	headers map[string]string
	query   map[string]string
	body    any
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
	json.NewDecoder(resp.Body).Decode(&res)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("api return status %s: %s", resp.Status, res.Error)
	}

	if !res.Success {
		return fmt.Errorf("API error: %s", res.Error)
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
	exiresAt := time.Unix(jwt.ExpiresAt, 0)
	if time.Now().Add(30*time.Second).After(exiresAt) {
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
		"User-Agent":      r.Header.Get("User-Agent"),
		"X-Forwarded-For": clientIP,
	}
	if err := c.doRequest(r.Context(), "POST", "/login", headers, body, &jwt); err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return &jwt, nil
}

func (c *CoreClient) Logout(w http.ResponseWriter, r *http.Request) error {
	req := request {
		method: "GET",
		path: "/logout",
	}

	c.sessionManager.Logout(w, r)

	return c.Do(w, r, req, nil)
}

func (c *CoreClient) Do(w http.ResponseWriter, r *http.Request, req request, result any) error {
	accessToken, err := c.getAccessToken(w, r)
	if err != nil {
		return fmt.Errorf("get access token: %w", err)
	}

	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", accessToken)

	for key, value := range req.headers {
		headers[key] = value
	}

	u, err := url.Parse(req.path)
	if err != nil {
		return fmt.Errorf("parse path: %w", err)
	}

	q := u.Query()
	for key, value := range req.query {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	return c.doRequest(r.Context(), req.method, u.String(), headers, req.body, result)
}
