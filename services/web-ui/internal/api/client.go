package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"web-ui/internal/models"
)

type CoreClient struct {
	baseURL    string
	httpClient *http.Client
	expireDuration time.Duration
}

func NewCoreClient(baseURL string, timeout, expireDuration time.Duration) *CoreClient {
	return &CoreClient{
		baseURL: baseURL,
		expireDuration: expireDuration,
		httpClient: &http.Client{Timeout: timeout},
	}
}

type request struct {
	method  string
	path    string
	headers map[string]string
	query   map[string]string
	body    any
}

func (c *CoreClient) do(ctx context.Context, req *request, result any) error {
	var reqBody []byte
	if req.body != nil {
		var err error
		reqBody, err = json.Marshal(req.body)
		if err != nil {
			return fmt.Errorf("marshal body: %w", err)
		}
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

	r, err := http.NewRequestWithContext(ctx, req.method, c.baseURL+u.String(), bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	r.Header.Set("Content-Type", "application/json")
	for key, value := range req.headers {
		r.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(r)
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

func (c *CoreClient) doWithAuth(ctx context.Context, s *Session, req *request, resp any) error {
	if s.Exire(c.expireDuration) {
		err := c.refreshToken(ctx, s)
		log.Printf("err refresh token: %v", err)
	}
	if req.headers == nil {
		req.headers = make(map[string]string)
	}

	req.headers["Authorization"] = "Bearer "+s.AccessToken

	return c.do(ctx, req, resp)
}

func (c *CoreClient) refreshToken(ctx context.Context, session *Session) error {
	body := map[string]string{
		"refresh_token": session.RefreshToken,
		"session_id":    session.SessionID,
	}
	r := &request{
		method:  "POST",
		path:    "/refresh",
		headers: map[string]string{},
		body:    body,
	}
	var jwt JWT
	if err := c.do(ctx, r, &jwt); err != nil {
		return err
	}

	session.Update(jwt.AccessToken, jwt.RefreshToken, jwt.ExiresAt)

	return nil
}

func (c *CoreClient) Login(r *http.Request, username, password string) (*Session, error) {
	body := map[string]string{"username": username, "password": password}
	var clientIP string
	var jwt JWT
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
	if err := c.do(r.Context(), &request{
		method:  "POST",
		path:    "/login",
		headers: headers,
		body:    body,
	}, &jwt); err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	var session Session
	session.Set(jwt.AccessToken, jwt.RefreshToken, jwt.SessionID, jwt.ExiresAt)

	return &session, nil
}

func (c *CoreClient) Logout(ctx context.Context, s *Session) error {
	req := request{
		method: "GET",
		path:   "/logout",
	}

	return c.doWithAuth(ctx, s, &req, nil)
}
