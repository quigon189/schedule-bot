package llm

import (
	"bytes"
	"encoding/json"
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
	Options map[string]any `json:"options,omitempty"`
}

type GenerateResponse struct {
	Response string `json:"response"`
}

func NewOllamaClient(url, model string) *OllamaClient {
	return &OllamaClient{
		url:   url,
		model: model,
		http:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *OllamaClient) Generate(prompt string, opts Options) (string, error) {
	reqBody := GenerateRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: false,
		Options: opts,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := c.http.Post(c.url+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var genResp GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		return "", err
	}

	return genResp.Response, nil
}
