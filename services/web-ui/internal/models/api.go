package models

import "encoding/json"

type Result struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
	Error   string          `json:"error"`
}
