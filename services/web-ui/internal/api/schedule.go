package api

import (
	"context"
	"fmt"
	"strconv"
	"web-ui/internal/models"
)

// GetGroupSchedule возвращает расписание группы
func (c *CoreClient) GetGroupSchedule(ctx context.Context, s *Session, groupID int, periodID *int) (*models.ScheduleData, error) {
	var entries []models.ScheduleEntry
	req := &request{
		method: "GET", 
		path: fmt.Sprintf("/schedule/group/%d", groupID),
	}
	if periodID != nil {
		req.query = map[string]string{
			"period_id": strconv.Itoa(*periodID),
		}
	}
	if err := c.doWithAuth(ctx, s, req, &entries); err != nil {
		return nil, err
	}

	return &models.ScheduleData{Entries: entries}, nil
}

// GetFilteredSchedule получает расписание с фильтрацией (группа, преподаватель, аудитория, период)
func (c *CoreClient) GetFilteredSchedule(ctx context.Context, s *Session, groupID, teacherID, audienceID, periodID *int) ([]models.ScheduleEntry, error) {
	path := "/schedule?"
	if groupID != nil {
		path += fmt.Sprintf("group_id=%d&", *groupID)
	}
	if teacherID != nil {
		path += fmt.Sprintf("teacher_id=%d&", *teacherID)
	}
	if audienceID != nil {
		path += fmt.Sprintf("audience_id=%d&", *audienceID)
	}
	if periodID != nil {
		path += fmt.Sprintf("academic_period_id=%d&", *periodID)
	}
	// убираем последний &
	if len(path) > 0 && path[len(path)-1] == '&' {
		path = path[:len(path)-1]
	}
	var entries []models.ScheduleEntry
	req := &request{method: "GET", path: path}
	if err := c.doWithAuth(ctx, s, req, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}
