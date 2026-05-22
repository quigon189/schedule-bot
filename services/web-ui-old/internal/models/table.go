package models

type Column struct {
	Key   string
	Title string
	Width string //css ширина
}

type TableConfig struct {
	ID         string
	Title      string
	Columns    []Column
	HasActions bool
	EditURL    string //шаблон с %d
	DeleteURL  string
}

type TableRowData map[string]any

type PaginatedTableData struct {
	Rows       []TableRowData
	Total      int
	TotalPages int
	Page       int
	PerPage    int
}
