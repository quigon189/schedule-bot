package api

import (
	"context"
	"fmt"
	"web-ui/internal/models"
)

type PlanScheduleRequest struct {
	AcademicPeriodID int `json:"academic_period_id"`
	Days             int `json:"days"`
	Slots            int `json:"slots"`
	Seed             int `json:"seed"`
	SubjectList      []struct {
		SubjectID    int `json:"subject_id"`
		TeacherID    int `json:"teacher_id"`
		AudienceID   int `json:"audience_id"`
		Priority     int `json:"priority"`
		LessonsCount int `json:"lessons_count"`
	} `json:"subject_list"`
}

type WeeklySchedule struct {
	Seed   int `json:"seed"`
	Grid   map[int]map[int][]struct {
		WeekType int            `json:"week_type"`
		Subject  models.Subject `json:"subject"`
		Teacher  models.Teacher `json:"teacher"`
		Audience models.Audience `json:"audience"`
	} `json:"grid"`
	Groups []models.Group `json:"groups"`
}

func (c *CoreClient) GenerateWeeklySchedule(ctx context.Context, s *Session, req PlanScheduleRequest) (*WeeklySchedule, error) {
	var result WeeklySchedule
	apiReq := &request{method: "POST", path: "/schedule/generate", body: req}
	if err := c.doWithAuth(ctx, s, apiReq, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *CoreClient) DownloadPlannerTemplate(ctx context.Context, s *Session, periodID int) ([]byte, error) {
	data, _, err := c.doRawWithAuth(ctx, s, &request{method: "GET", path: fmt.Sprintf("/schedule/template/%d", periodID)})
	return data, err
}

func (c *CoreClient) UploadPlannerExcel(ctx context.Context, s *Session, fileData []byte, filename string) (*models.GenerationResult, error) {
	var result models.GenerationResult
	err := c.doMultipartWithAuth(ctx, s, "/schedule/generate/upload", fileData, filename, nil, &result)
	return &result, err
}
