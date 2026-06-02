package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
	"web-ui/internal/models"
)

type CoreClient struct {
	baseURL        string
	httpClient     *http.Client
	expireDuration time.Duration
}

func NewCoreClient(baseURL string, timeout, expireDuration time.Duration) *CoreClient {
	return &CoreClient{
		baseURL:        baseURL,
		expireDuration: expireDuration,
		httpClient:     &http.Client{Timeout: timeout},
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
	if h, ok := req.headers["Content-Type"]; ok {
		r.Header.Set("Content-Type", h)
	} else {
		r.Header.Set("Content-Type", "application/json")
	}
	for key, value := range req.headers {
		r.Header.Set(key, value)
	}

	log.Printf("do request: %+v", r)

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
	if s.Expire(c.expireDuration) {
		err := c.refreshToken(ctx, s)
		if err != nil {
			return fmt.Errorf("refresh token: %w", err)
		}
	}
	if req.headers == nil {
		req.headers = make(map[string]string)
	}

	req.headers["Authorization"] = "Bearer " + s.AccessToken

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

func (c *CoreClient) doRaw(ctx context.Context, req *request) ([]byte, string, error) {
	// Собираем URL
	u, err := url.Parse(req.path)
	if err != nil {
		return nil, "", fmt.Errorf("parse path: %w", err)
	}
	q := u.Query()
	for key, value := range req.query {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	// Тело запроса (если есть, например для POST)
	var reqBody io.Reader
	if req.body != nil {
		jsonBody, err := json.Marshal(req.body)
		if err != nil {
			return nil, "", fmt.Errorf("marshal body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.method, c.baseURL+u.String(), reqBody)
	if err != nil {
		return nil, "", fmt.Errorf("new request: %w", err)
	}
	// По умолчанию Content-Type: application/json, но для файлов он может быть другим.
	// Устанавливаем, только если не переопределено.
	if _, ok := req.headers["Content-Type"]; !ok && req.body != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	for k, v := range req.headers {
		httpReq.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("read body: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Попробуем распарсить ошибку как JSON от core
		var errRes models.Result
		if json.Unmarshal(body, &errRes) == nil && !errRes.Success {
			return nil, "", fmt.Errorf("API error: %s", errRes.Error)
		}
		return nil, "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	return body, resp.Header.Get("Content-Type"), nil
}

// doRawWithAuth — аналог doWithAuth для raw запросов (с автообновлением токена).
func (c *CoreClient) doRawWithAuth(ctx context.Context, s *Session, req *request) ([]byte, string, error) {
	if s.Expire(c.expireDuration) {
		if err := c.refreshToken(ctx, s); err != nil {
			return nil, "", fmt.Errorf("refresh token: %w", err)
		}
	}
	if req.headers == nil {
		req.headers = make(map[string]string)
	}
	req.headers["Authorization"] = "Bearer " + s.AccessToken
	return c.doRaw(ctx, req)
}

// ---------- Загрузка файлов (multipart) ----------
// doMultipartWithAuth отправляет multipart/form-data с файлом и возвращает JSON-ответ.
func (c *CoreClient) doMultipartWithAuth(ctx context.Context, s *Session, path string, fileData []byte, filename string, extraFields map[string]string, result any) error {
	if s.Expire(c.expireDuration) {
		if err := c.refreshToken(ctx, s); err != nil {
			return fmt.Errorf("refresh token: %w", err)
		}
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	// Добавляем файл
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return fmt.Errorf("create form file: %w", err)
	}
	if _, err := part.Write(fileData); err != nil {
		return fmt.Errorf("write file data: %w", err)
	}
	// Добавляем дополнительные поля
	for k, v := range extraFields {
		if err := writer.WriteField(k, v); err != nil {
			return fmt.Errorf("write field %s: %w", k, err)
		}
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("close writer: %w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+path, body)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	httpReq.Header.Set("Authorization", "Bearer "+s.AccessToken)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}
	var res models.Result
	if err := json.Unmarshal(respBody, &res); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, res.Error)
	}
	if !res.Success {
		return fmt.Errorf("API error: %s", res.Error)
	}
	if result != nil && res.Data != nil {
		return json.Unmarshal(res.Data, result)
	}
	return nil
}
