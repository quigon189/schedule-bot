package api

import (
	"context"
	"fmt"
	"strconv"
	"web-ui/internal/models"
)

type CreateScheduleTemplateRequest struct {
	DayOfWeek        int `json:"day_of_week"`
	Number           int `json:"number"`
	WeekType         int `json:"week_type"`
	SubjectID        int `json:"subject_id"`
	TeacherID        int `json:"teacher_id"`
	AudienceID       int `json:"audience_id"`
	AcademicPeriodID int `json:"academic_period_id"`
}

func (c *CoreClient) CreateScheduleTemplate(ctx context.Context, s *Session, req CreateScheduleTemplateRequest) (*models.ScheduleTemplate, error) {
	var tmpl models.ScheduleTemplate
	apiReq := &request{method: "POST", path: "/schedule", body: req}
	if err := c.doWithAuth(ctx, s, apiReq, &tmpl); err != nil {
		return nil, err
	}
	return &tmpl, nil
}

type UpdateScheduleTemplateRequest struct {
	DayOfWeek        *int `json:"day_of_week,omitempty"`
	Number           *int `json:"number,omitempty"`
	WeekType         *int `json:"week_type,omitempty"`
	SubjectID        *int `json:"subject_id,omitempty"`
	TeacherID        *int `json:"teacher_id,omitempty"`
	AudienceID       *int `json:"audience_id,omitempty"`
	AcademicPeriodID *int `json:"academic_period_id,omitempty"`
}

func (c *CoreClient) UpdateScheduleTemplate(ctx context.Context, s *Session, id int, req UpdateScheduleTemplateRequest) (*models.ScheduleTemplate, error) {
	var tmpl models.ScheduleTemplate
	apiReq := &request{method: "PATCH", path: fmt.Sprintf("/schedule/%d", id), body: req}
	if err := c.doWithAuth(ctx, s, apiReq, &tmpl); err != nil {
		return nil, err
	}
	return &tmpl, nil
}

func (c *CoreClient) DeleteScheduleTemplate(ctx context.Context, s *Session, id int) error {
	req := &request{method: "DELETE", path: fmt.Sprintf("/schedule/%d", id)}
	return c.doWithAuth(ctx, s, req, nil)
}

func (c *CoreClient) GetScheduleTemplate(ctx context.Context, s *Session, id int) (*models.ScheduleTemplate, error) {
	var tmpl models.ScheduleTemplate
	req := &request{method: "GET", path: fmt.Sprintf("/schedule/%d", id)}
	if err := c.doWithAuth(ctx, s, req, &tmpl); err != nil {
		return nil, err
	}
	return &tmpl, nil
}

type ScheduleFilters struct {
	GroupID          *int `json:"group_id,omitempty"`
	SubjectID        *int `json:"subject_id,omitempty"`
	TeacherID        *int `json:"teacher_id,omitempty"`
	AudienceID       *int `json:"audience_id,omitempty"`
	DayOfWeek        *int `json:"day_of_week,omitempty"`
	WeekType         *int `json:"week_type,omitempty"`
	AcademicPeriodID *int `json:"academic_period_id,omitempty"`
}

func (c *CoreClient) GetScheduleTemplates(ctx context.Context, s *Session, filters ScheduleFilters) ([]models.ScheduleTemplate, error) {
	query := make(map[string]string)
	if filters.GroupID != nil {
		query["group_id"] = strconv.Itoa(*filters.GroupID)
	}
	if filters.SubjectID != nil {
		query["subject_id"] = strconv.Itoa(*filters.SubjectID)
	}
	if filters.TeacherID != nil {
		query["teacher_id"] = strconv.Itoa(*filters.TeacherID)
	}
	if filters.AudienceID != nil {
		query["audience_id"] = strconv.Itoa(*filters.AudienceID)
	}
	if filters.DayOfWeek != nil {
		query["day_of_week"] = strconv.Itoa(*filters.DayOfWeek)
	}
	if filters.WeekType != nil {
		query["week_type"] = strconv.Itoa(*filters.WeekType)
	}
	if filters.AcademicPeriodID != nil {
		query["academic_period_id"] = strconv.Itoa(*filters.AcademicPeriodID)
	}
	var templates []models.ScheduleTemplate
	req := &request{method: "GET", path: "/schedule", query: query}
	if err := c.doWithAuth(ctx, s, req, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

// GetGroupSchedule возвращает расписание группы (шаблоны)
func (c *CoreClient) GetGroupSchedule(ctx context.Context, s *Session, groupID int, periodID *int) ([]models.ScheduleTemplate, error) {
	path := fmt.Sprintf("/schedule/group/%d", groupID)
	query := make(map[string]string)
	if periodID != nil {
		query["period_id"] = strconv.Itoa(*periodID)
	}
	var templates []models.ScheduleTemplate
	req := &request{method: "GET", path: path, query: query}
	if err := c.doWithAuth(ctx, s, req, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

func (c *CoreClient) GetTeacherSchedule(ctx context.Context, s *Session, teacherID int, periodID *int) ([]models.ScheduleTemplate, error) {
	path := fmt.Sprintf("/schedule/teacher/%d", teacherID)
	query := make(map[string]string)
	if periodID != nil {
		query["period_id"] = strconv.Itoa(*periodID)
	}
	var templates []models.ScheduleTemplate
	req := &request{method: "GET", path: path, query: query}
	if err := c.doWithAuth(ctx, s, req, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

func (c *CoreClient) GetAudienceSchedule(ctx context.Context, s *Session, audienceID int, periodID *int) ([]models.ScheduleTemplate, error) {
	path := fmt.Sprintf("/schedule/audience/%d", audienceID)
	query := make(map[string]string)
	if periodID != nil {
		query["period_id"] = strconv.Itoa(*periodID)
	}
	var templates []models.ScheduleTemplate
	req := &request{method: "GET", path: path, query: query}
	if err := c.doWithAuth(ctx, s, req, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

type SemesterScheduleRequest struct {
	GroupID          int `json:"group_id"`
	AcademicPeriodID int `json:"academic_period_id"`
	ScheduleEntries  []struct {
		DayOfWeek  int `json:"day_of_week"`
		Number     int `json:"number"`
		WeekType   int `json:"week_type"`
		SubjectID  int `json:"subject_id"`
		TeacherID  int `json:"teacher_id"`
		AudienceID int `json:"audience_id"`
	} `json:"schedule_entries"`
}

func (c *CoreClient) CreateSemesterSchedule(ctx context.Context, s *Session, req SemesterScheduleRequest) error {
	apiReq := &request{method: "POST", path: "/schedule/semester", body: req}
	return c.doWithAuth(ctx, s, apiReq, nil)
}

func (c *CoreClient) ExportGroupSchedule(ctx context.Context, s *Session, groupID int, periodID *int, format string) ([]byte, string, error) {
	path := fmt.Sprintf("/schedule/group/%d/export", groupID)
	query := map[string]string{"format": format}
	if periodID != nil {
		query["period_id"] = strconv.Itoa(*periodID)
	}
	req := &request{method: "GET", path: path, query: query}
	return c.doRawWithAuth(ctx, s, req)
}

func (c *CoreClient) ExportTeacherSchedule(ctx context.Context, s *Session, teacherID int, periodID *int, format string) ([]byte, string, error) {
	path := fmt.Sprintf("/schedule/teacher/%d/export", teacherID)
	query := map[string]string{"format": format}
	if periodID != nil {
		query["period_id"] = strconv.Itoa(*periodID)
	}
	req := &request{method: "GET", path: path, query: query}
	return c.doRawWithAuth(ctx, s, req)
}

func (c *CoreClient) ExportAudienceSchedule(ctx context.Context, s *Session, audienceID int, periodID *int, format string) ([]byte, string, error) {
	path := fmt.Sprintf("/schedule/audience/%d/export", audienceID)
	query := map[string]string{"format": format}
	if periodID != nil {
		query["period_id"] = strconv.Itoa(*periodID)
	}
	req := &request{method: "GET", path: path, query: query}
	return c.doRawWithAuth(ctx, s, req)
}
