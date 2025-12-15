package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Sender struct {
	url    string
	client *http.Client
}

func NewSender(url string) *Sender {
	return &Sender{
		url: url,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *Sender) Send(data any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()


	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Schedule service return status: %s", resp.Status)
	}
	
	return nil
}
