package models

import "time"

type ChangesRequest struct {
	Date time.Time `json:"date" jsonschema:"description=Дата изменений в формате ISO 8601, если значение пустое, то возвращает изменения на текущую дату"`
}

type GroupScheduleRequest struct {
	GroupName    string `json:"group_name" jsonschema:"description=Название учебной группы в формате АБ-123, где АБ короткое наименование специальности а 123 номер группы. Например СА-501 Системное и сетевое администрирование группа 501"`
	AcademicYear string `json:"academic_year" jsonschema:"description=Академический год в формате YYYY/YYYY. Например 2025/2026"`
	HalfYear     int    `json:"half_year" jsonschema:"description=Полугодие. Одно из двух значение 1 или 2"`
}

type ChangesResponse struct {
	Data struct {
		ID          int       `json:"id"`
		Date        string    `json:"date"`
		ImgURLs     []string  `json:"image_urls"`
		Description string    `json:"description"`
		CreatedAt   time.Time `json:"created_at"`
	} `json:"data"`
}

type GroupScheduleResponse struct {
	Data []struct {
		AcademicYear   string    `json:"academic_year"`
		HalfYear       int       `json:"half_year"`
		GroupName      string    `json:"group_name"`
		Semester       int       `json:"semester"`
		ScheduleImgURL string    `json:"schedule_img_url"`
		CreatedAt      time.Time `json:"created_at"`
	} `json:"data"`
}

type SendMessage struct {
	ChatID    int64    `json:"chat_id"`
	Message   string   `json:"message"`
	PhotoURLs []string `json:"photo_urls"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    []byte `json:"data"`
	Error   string `json:"error,omitempty"`
}
