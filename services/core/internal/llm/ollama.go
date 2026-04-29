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
	Role      string   `json:"role"`
	Content   string   `json:"content"`
	Images    []string `json:"images,omitempty"`
	ToolCalls []struct {
		Function ollamaFunctionCall `json:"function"`
	} `json:"tool_calls"`
}

type ollamaFunctionCall struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

type ollamaChatRequest struct {
	Model    string                  `json:"model"`
	Messages []ollamaMessage         `json:"messages"`
	Stream   bool                    `json:"stream"`
	Options  map[string]any          `json:"options,omitempty"`
	Tools    []ollamaToolDescroption `json:"tools,omitempty"`
}

type ollamaToolDescroption struct {
	Type     string `json:"type"`
	Function struct {
		Name        string          `json:"name"`
		Description string          `json:"description"`
		Parameters  json.RawMessage `json:"parameters"`
	}
}

type ollamaChatResponse struct {
	Message         ollamaMessage `json:"message"`
	PromptEvalCount int           `json:"prompt_eval_count"`
	EvalCount       int           `json:"eval_count"`
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
		if len(m.ToolCalls) > 0 {
			for _, tc := range m.ToolCalls {
				message.ToolCalls = append(message.ToolCalls, struct {
					Function ollamaFunctionCall `json:"function"`
				}{
					ollamaFunctionCall{
						Name:      tc.Name,
						Arguments: tc.Arguments,
					},
				})
			}
		}
		reqMessages = append(reqMessages, message)
	}

	var toolDescs []ollamaToolDescroption
	if tools, ok := opts["tools"].([]ToolDescription); ok {
		for _, t := range tools {
			toolDescs = append(toolDescs, ollamaToolDescroption{
				Type: t.Type,
				Function: struct {
					Name        string          `json:"name"`
					Description string          `json:"description"`
					Parameters  json.RawMessage `json:"parameters"`
				}{
					Name:        t.Function.Name,
					Description: t.Function.Description,
					Parameters:  t.Function.Parameters,
				},
			})
		}
	}

	reqBody := ollamaChatRequest{
		Model:    c.model,
		Messages: reqMessages,
		Stream:   false,
		Options:  opts,
		Tools:    toolDescs,
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

	m := Message{
		Role:    chatResp.Message.Role,
		Content: chatResp.Message.Content,
	}

	if len(chatResp.Message.ToolCalls) > 0 {
		for _, tc := range chatResp.Message.ToolCalls {
			m.ToolCalls = append(m.ToolCalls, ToolCall{
				Name:      tc.Function.Name,
				Arguments: tc.Function.Arguments,
			})
		}
	}

	return &m, nil
}
