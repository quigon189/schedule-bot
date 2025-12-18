package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"mcp_server/internal/service"

	"github.com/mark3labs/mcp-go/mcp"
)

type SchedulesHandler struct {
	service *service.ChangeService
}

func NewSchedulesHandler(s *service.ChangeService) *SchedulesHandler {
	return &SchedulesHandler{
		service: s,
	}
}

func (h *SchedulesHandler) CurrentDateHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	resultText := fmt.Sprintf("Текущая дата в формате ISO 8601: %s", time.Now().Format("2006-01-02"))
	return mcp.NewToolResultText(resultText), nil
}

func (h *SchedulesHandler) CurrentStudyPeriodHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var studyPeriod string
	var half_yaer int
	currentYear, currentMonth, _ := time.Now().Date()
	if currentMonth >= 7 {
		half_yaer = 1
		studyPeriod = fmt.Sprintf("%d/%d", currentYear, currentYear+1)
	} else {
		half_yaer = 2
		studyPeriod = fmt.Sprintf("%d/%d", currentYear-1, currentYear)
	}

	resultText := fmt.Sprintf("Текущий учебный период: %s\nТекущее полугодие: %d", studyPeriod, half_yaer)
	return mcp.NewToolResultText(resultText), nil
}

func (h *SchedulesHandler) GetGroupSchedule(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	groupName, err := req.RequireString("group_name")
	if err != nil {
		log.Printf("invalid group_name: %v", err)
		return mcp.NewToolResultError("невалидное название группы"), nil
	}

	academicYear, err := req.RequireString("academic_year")
	if err != nil {
		log.Printf("invalid academic_year: %v", err)
		return mcp.NewToolResultError("невалидный учебный год"), nil
	}

	halfYear, err := req.RequireString("half_year")
	if err != nil {
		log.Printf("invalid half_year: %v", err)
		return mcp.NewToolResultError("невалидное полугодие"), nil
	}

	log.Printf("%s %s %s", groupName, academicYear, halfYear)

	gsResp, err := h.service.GetSchedules(groupName, academicYear, halfYear)
	if err != nil {
		log.Printf("failed getting group schedules from service: %v", err)
		return mcp.NewToolResultError("ошибка при получении расписания группы из сервиса"), nil
	}

	log.Printf("resp: %+v", gsResp)

	jsonData, err := json.Marshal(gsResp.Data)
	if err != nil {
		log.Printf("Failed to marshal changes response: %v", err)
		return mcp.NewToolResultError("ошибка при обработке данных"), nil
	}

	return mcp.NewToolResultText("Ответ от сервиса расписаний в формате json: "+string(jsonData)), err
}
