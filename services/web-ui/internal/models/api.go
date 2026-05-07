package models

import "encoding/json"

type Result struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
	Error   string          `json:"error"`
}

type PaginatedUsers struct {
	Users      []User `json:"users"`
	Total      int    `json:"total"`
	TotalPages int    `json:"total_pages"`
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
}

type PaginatedUsersQuery struct {
	Page      *int    `json:"page,string,omitempty"`
	PerPage   *int    `json:"per_page,string,omitempty"`
	SortBy    *string `json:"sort_by,omitempty"`
	SortOrder *string `json:"sort_order,omitempty"`
	FullName  *string `json:"full_name,omitempty"`
	Username  *string `json:"username,omitempty"`
	Email     *string `json:"email,omitempty"`
}
