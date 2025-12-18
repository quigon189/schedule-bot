package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"mcp_server/internal/config"
	"mcp_server/internal/models"
)

type ChangeService struct {
	BaseURL string
	Client  *http.Client
}

func NewChangeService(cfg *config.Config) *ChangeService {
	timeout, err := strconv.Atoi(cfg.Timeout)
	if err != nil {
		timeout = 30
	}
	return &ChangeService{
		BaseURL: cfg.ScheduleServiceURL + "/api/v1",
		Client:  &http.Client{Timeout: time.Duration(timeout) * time.Second},
	}
}

func (s *ChangeService) GetSchedules(group string, ay string, hy string) (*models.GroupScheduleResponse, error) {
	u, err := url.Parse(s.BaseURL + "/group_schedules")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("group_name", group)
	q.Add("academic_year", ay)
	q.Add("half_year", hy)
	u.RawQuery = q.Encode()

	
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}


	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http error status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}


	gsResp := models.GroupScheduleResponse{}
	if err := json.Unmarshal(body, &gsResp); err != nil {
		return nil, err
	}

	return &gsResp, nil
}

func (s *ChangeService) GetChanges(date string) (*models.ChangesResponse, error) {
	u, err := url.Parse(s.BaseURL + "/changes")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("date", date)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	changesResponse := models.ChangesResponse{}
	if err := json.Unmarshal(body, &changesResponse); err != nil {
		return nil, err
	}

	log.Printf("Decoded body: %v", changesResponse)

	return &changesResponse, nil
}
