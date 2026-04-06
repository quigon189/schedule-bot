package llm

import (
	"bytes"
	"context"
	"core/internal/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type OllamaClient struct {
	url   string
	model string
	http  *http.Client
}

type GenerateRequest struct {
	Model   string         `json:"model"`
	Prompt  string         `json:"prompt"`
	Stream  bool           `json:"stream"`
	System  string         `json:"system,omitempty"`
	Options map[string]any `json:"options,omitempty"`
}

type GenerateResponse struct {
	Response string `json:"response"`
}

func NewOllamaClient(cfg *config.OllamaConfig) *OllamaClient {
	return &OllamaClient{
		url:   "http://" + cfg.URL,
		model: cfg.Model,
		http:  &http.Client{Timeout: time.Duration(cfg.Timeout) * time.Second},
	}
}

func (c *OllamaClient) Generate(ctx context.Context, systemPrompt, userMessage string, opts Options) (string, error) {
	reqBody := GenerateRequest{
		Model:   c.model,
		Prompt:  userMessage,
		System:  systemPrompt,
		Stream:  false,
		Options: opts,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	log.Printf("Запрос для Ollama: %s", string(jsonData))

	req, err := http.NewRequestWithContext(ctx, "POST", c.url+"/api/generate", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("post request: %w", err)
	}
	defer resp.Body.Close()

	var genResp GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		return "", err
	}

	return genResp.Response, nil
}
