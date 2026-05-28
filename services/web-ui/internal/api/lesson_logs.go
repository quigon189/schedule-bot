package api

import (
	"context"
	"fmt"
	"strconv"
	"web-ui/internal/models"
)

type LessonLogFilters struct {
	GroupID          *int    `json:"group_id,omitempty"`
	SubjectID        *int    `json:"subject_id,omitempty"`
	TeacherID        *int    `json:"teacher_id,omitempty"`
	AudienceID       *int    `json:"audience_id,omitempty"`
	AcademicPeriodID *int    `json:"academic_period_id,omitempty"`
	DateFrom         *string `json:"date_from,omitempty"`
	DateTo           *string `json:"date_to,omitempty"`
	Status           *string `json:"status,omitempty"`
	Number           *int    `json:"number,omitempty"`
}

func (c *CoreClient) GetLessonLogs(ctx context.Context, s *Session, filters LessonLogFilters) ([]models.LessonLog, error) {
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
	if filters.AcademicPeriodID != nil {
		query["academic_period_id"] = strconv.Itoa(*filters.AcademicPeriodID)
	}
	if filters.DateFrom != nil {
		query["date_from"] = *filters.DateFrom
	}
	if filters.DateTo != nil {
		query["date_to"] = *filters.DateTo
	}
	if filters.Status != nil {
		query["status"] = *filters.Status
	}
	if filters.Number != nil {
		query["number"] = strconv.Itoa(*filters.Number)
	}
	var logs []models.LessonLog
	req := &request{method: "GET", path: "/lessons", query: query}
	if err := c.doWithAuth(ctx, s, req, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

func (c *CoreClient) CancelLesson(ctx context.Context, s *Session, logID int, comment string) error {
	body := map[string]string{"comment": comment}
	req := &request{method: "POST", path: fmt.Sprintf("/lessons/cancel/%d", logID), body: body}
	return c.doWithAuth(ctx, s, req, nil)
}

type RescheduleLessonRequest struct {
	Date       string `json:"date"`
	Number     int    `json:"number"`
	SubjectID  int    `json:"subject_id"`
	TeacherID  int    `json:"teacher_id"`
	AudienceID int    `json:"audience_id"`
	Comment    string `json:"comment"`
}

func (c *CoreClient) RescheduleLesson(ctx context.Context, s *Session, req RescheduleLessonRequest) error {
	apiReq := &request{method: "POST", path: "/lessons/reschedule", body: req}
	return c.doWithAuth(ctx, s, apiReq, nil)
}

func (c *CoreClient) CompleteLessonsByDate(ctx context.Context, s *Session, date, comment string) error {
	body := map[string]string{"date": date, "comment": comment}
	req := &request{method: "POST", path: "/lessons/complete", body: body}
	return c.doWithAuth(ctx, s, req, nil)
}

type LessonStatistics struct {
	TotalSubjects       int     `json:"total_subjects"`
	OverallProgress     float64 `json:"overall_progress"`
	TotalLessonsAll     int     `json:"total_lessons_all"`
	CompletedLessonsAll int     `json:"completed_lessons_all"`
	RemainingLessonsAll int     `json:"remaining_lessons_all"`
	Subjects            []struct {
		SubjectID          int     `json:"subject_id"`
		SubjectTitle       string  `json:"subject_title"`
		GroupID            int     `json:"group_id"`
		GroupName          string  `json:"group_name"`
		TotalHours         int     `json:"total_hours"`
		TotalLessons       int     `json:"total_lessons"`
		CompletedLessons   int     `json:"completed_lessons"`
		RescheduledLessons int     `json:"rescheduled_lessons"`
		CancelledLessons   int     `json:"cancelled_lessons"`
		PlannedLessons     int     `json:"planned_lessons"`
		RemainingLessons   int     `json:"remaining_lessons"`
		CompletionPercent  float64 `json:"completion_percent"`
		IsOnSchedule       bool    `json:"is_on_schedule"`
	} `json:"subjects"`
}

func (c *CoreClient) GetLessonStatistics(ctx context.Context, s *Session, groupID, teacherID, subjectID, periodID *int) (*LessonStatistics, error) {
	query := make(map[string]string)
	if groupID != nil {
		query["group_id"] = strconv.Itoa(*groupID)
	}
	if teacherID != nil {
		query["teacher_id"] = strconv.Itoa(*teacherID)
	}
	if subjectID != nil {
		query["subject_id"] = strconv.Itoa(*subjectID)
	}
	if periodID != nil {
		query["academic_period_id"] = strconv.Itoa(*periodID)
	}
	var stats LessonStatistics
	req := &request{method: "GET", path: "/lessons/statistics", query: query}
	if err := c.doWithAuth(ctx, s, req, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}
