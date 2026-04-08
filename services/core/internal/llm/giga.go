package llm

import (
	"bytes"
	"context"
	"core/internal/config"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type GigaChatClient struct {
	clientID         string
	clientSecret     string
	authorizationKey string
	baseURL          string
	model            string
	http             *http.Client
	token            string
	tokenExpiry      time.Time
	tokenMu          sync.RWMutex
}

func NewGigaChatClient(cfg *config.GigaChatConfig) *GigaChatClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &GigaChatClient{
		clientID:         cfg.ClientID,
		clientSecret:     cfg.ClientSecret,
		authorizationKey: cfg.AuthorizationKey,
		baseURL:          cfg.BaseURL,
		model:            cfg.Model,
		http: &http.Client{
			Timeout:   time.Duration(cfg.Timeout) * time.Second,
			Transport: tr,
		},
	}
}

// getToken получает или обновляет access token
func (c *GigaChatClient) getToken(ctx context.Context) (string, error) {
	c.tokenMu.RLock()
	if c.token != "" && time.Now().Before(c.tokenExpiry) {
		token := c.token
		c.tokenMu.RUnlock()
		return token, nil
	}
	c.tokenMu.RUnlock()

	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()

	// Double-check после получения блокировки
	if c.token != "" && time.Now().Before(c.tokenExpiry) {
		return c.token, nil
	}

	// Запрос токена
	url := c.baseURL + "/oauth"
	body := bytes.NewBufferString("scope=GIGACHAT_API_PERS")

	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return "", fmt.Errorf("create token request: %w", err)
	}
	if c.authorizationKey != "" {
		req.Header.Set("Authorization", "Basic "+c.authorizationKey)
	} else {
		req.SetBasicAuth(c.clientID, c.clientSecret)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("RqUID", generateUUID())

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("response status %d", resp.StatusCode)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresAt   int64 `json:"expires_at"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("decode token response: %w", err)
	}

	c.token = tokenResp.AccessToken
	c.tokenExpiry = time.Unix(tokenResp.ExpiresAt, 0)

	return c.token, nil
}

// generateUUID генерирует простой UUID для RqUID
func generateUUID() string {
	id := uuid.New()
	return id.String()
}

type gigachatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type gigachatRequest struct {
	Model    string            `json:"model"`
	Messages []gigachatMessage `json:"messages"`
	Stream   bool              `json:"stream"`
	Temp     float64           `json:"temperature,omitempty"`
}

type gigachatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (c *GigaChatClient) Generate(ctx context.Context, systemPrompt, userPrompt string, opts Options) (string, error) {
	token, err := c.getToken(ctx)
	if err != nil {
		return "", fmt.Errorf("get token: %w", err)
	}

	messages := []gigachatMessage{}
	if systemPrompt != "" {
		messages = append(messages, gigachatMessage{Role: "system", Content: systemPrompt})
	}
	messages = append(messages, gigachatMessage{Role: "user", Content: userPrompt})

	var temp float64
	temp, ok := opts["temperature"].(float64)
	if !ok {
		temp = 0.1
	}

	reqBody := gigachatRequest{
		Model:    c.model,
		Messages: messages,
		Stream:   false,
		Temp:     temp,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	// другой url для запросов
	url := "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("response status: %d body: %s", resp.StatusCode, string(body))
	}

	var chatResp gigachatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("gigachat error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from gigachat")
	}

	return chatResp.Choices[0].Message.Content, nil
}

func (c *GigaChatClient) GetModelName() string {
	return c.model
}
