package dto

import "time"

type FileResponse struct {
	Path    string    `json:"file_path"`
	Expires time.Time `json:"expires"`
}
