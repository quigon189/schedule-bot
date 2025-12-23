package service

import (
	"bytes"
	"encoding/json"
	"mcp_server/internal/config"
	"mcp_server/internal/models"
	"net/http"
	"strconv"
	"time"
)

type TGService struct {
	BaseURL string
	Client  *http.Client
}

func NewTGService(cfg *config.Config) *TGService {
	timeout, err := strconv.Atoi(cfg.Timeout)
	if err != nil {
		timeout = 30
	}
	return &TGService{
		BaseURL: cfg.TGServiceURL,
		Client:  &http.Client{Timeout: time.Duration(timeout) * time.Second},
	}
}

func (s *TGService) SendMessage(sm *models.SendMessage) error {
	url := s.BaseURL + "/send_message"

	jsonData, err := json.Marshal(sm)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	_, err = s.Client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
