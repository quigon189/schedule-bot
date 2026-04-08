package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	BaseURL      string `yaml:"base_url"`
	AccessToken  string `yaml:"access_token"`
	RefreshToken string `yaml:"refresh_token"`
	SessionID    string `yaml:"session_id"`
	ExpiresAt    int64    `yaml:"expires_at"`
	Timeout      int    `yaml:"timeout"`
}

type UpdateToken struct {
	AccessToken  string
	RefreshToken string
	SessionID    string
	ExpiresAt    int64
}

type Manager struct {
	path   string
	config *Config
}

func NewManager() (*Manager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configDir := filepath.Join(home, ".config", "schedule-cli")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}
	configPath := filepath.Join(configDir, "config.yaml")
	cfg := &Config{
		BaseURL: "http://localhost:8080",
		Timeout: 30,
	}
	if data, err := os.ReadFile(configPath); err == nil {
		yaml.Unmarshal(data, cfg)
	}

	return &Manager{path: configPath, config: cfg}, nil
}

func (m *Manager) Get() *Config {
	return m.config
}

func (m *Manager) Save() error {
	data, err := yaml.Marshal(m.config)
	if err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	return os.WriteFile(m.path, data, 0644)
}

func (m *Manager) UpdateToken(req UpdateToken) error {
	m.config.AccessToken = req.AccessToken
	m.config.RefreshToken = req.RefreshToken
	m.config.SessionID = req.SessionID
	m.config.ExpiresAt = req.ExpiresAt
	return m.Save()
}

func (m *Manager) SetBaseURL(url string) error {
	m.config.BaseURL = url
	return m.Save()
}

func (m *Manager) SetTimeout(t int) error {
	m.config.Timeout = t
	return m.Save()
}
