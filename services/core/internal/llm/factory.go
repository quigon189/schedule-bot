package llm

import (
	"core/internal/config"
	"fmt"
)

func NewLLMClient(cfg *config.LLMConfig) (Client, error) {
	switch cfg.Provider {
	case "ollama":
		if cfg.Ollama.URL == "" {
			return nil, fmt.Errorf("ollama URL is required")
		}
		return NewOllamaClient(&cfg.Ollama), nil

	case "gigachat":
		if cfg.GigaChat.AuthorizationKey == "" {
			if cfg.GigaChat.ClientID == "" || cfg.GigaChat.ClientSecret == "" {
				return nil, fmt.Errorf("gigachat client credentials are required")
			}
		}
		return NewGigaChatClient(&cfg.GigaChat), nil

	default:
		return nil, fmt.Errorf("unknown LLM provider: %s", cfg.Provider)
	}
}
