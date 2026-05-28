package models

type AlertType string

const (
	AlertSuccess AlertType = "success"
	AlertError   AlertType = "danger"
	AlertInfo    AlertType = "info"
	AlertWarning AlertType = "warning"
)

type Alert struct {
	ID      string    `json:"id"`
	Type    AlertType `json:"type"`
	Message string    `json:"message"`
}
