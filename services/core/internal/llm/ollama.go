package llm

import (
	"bytes"
	"context"
	"core/internal/config"
	"encoding/base64"
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
	Think   bool           `json:"think"`
	Options map[string]any `json:"options,omitempty"`
}

type ollamaMessage struct {
	Role    string   `json:"role"`
	Content string   `json:"content"`
	Images  []string `json:"images,omitempty"`
}

type ollamaChatRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
	Options  map[string]any  `json:"options,omitempty"`
}

type ollamaChatResponse struct {
	Message         Message `json:"message"`
	PromptEvalCount int     `json:"prompt_eval_count"`
	EvalCount       int     `json:"eval_count"`
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
		Think:   false,
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

func (c *OllamaClient) Chat(ctx context.Context, messages []Message, opts Options) (*Message, error) {
	reqMessages := []ollamaMessage{}
	for _, m := range messages {
		message := ollamaMessage{
			Role:    m.Role,
			Content: m.Content,
		}
		if len(m.Images) > 0 {
			var images []string
			for _, i := range m.Images {
				base64Image := base64.StdEncoding.EncodeToString(i)
				images = append(images, base64Image)
			}
		}
		reqMessages = append(reqMessages, message)
	}
	reqBody := ollamaChatRequest{
		Model:    c.model,
		Messages: reqMessages,
		Stream:   false,
		Options:  opts,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.url+"/api/chat", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("post request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("response status %s", resp.Status)
	}

	var chatResp ollamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, err
	}

	log.Printf("Prompt tokens %d", chatResp.PromptEvalCount)
	log.Printf("Answer tokens %d", chatResp.EvalCount)

	return &chatResp.Message, nil
}
