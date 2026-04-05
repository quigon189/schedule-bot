package llm

import (
	"context"
	"core/internal/services"
	"encoding/json"
	"fmt"
)

type Options map[string]any

type Client interface {
	Generate(prompt string, opts Options) (string, error)
}

type LLMService struct {
	client          Client
	scheduleService services.ScheduleService
}

func NewLLMService(client Client, svc services.ScheduleService) *LLMService {
	return &LLMService{
		client: client,
		scheduleService: svc,
	}
}

func (s *LLMService) Call(ctx context.Context, action string, params map[string]any) (string, error) {
	switch action {
	case "get_group_schedule":
		groupID, ok := params["group_id"].(int)
		if !ok {
			return "", fmt.Errorf("group_id required")
		}

		var pID *int
		periodID, ok := params["period_id"].(int)
		if ok {
			pID = &periodID
		}

		schedule, err := s.scheduleService.GetGroupSchedule(ctx, groupID, pID)
		if err != nil {
			return "", fmt.Errorf("get group schedule: %w", err)
		}

		jsonData, err := json.Marshal(schedule)
		return string(jsonData), err
	default:
		return "", fmt.Errorf("unknown action: %s", action)
	}
}
