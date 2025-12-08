package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

type SchedulesHandler struct {
}

func NewSchedulesHandler() *SchedulesHandler {
	return &SchedulesHandler{}
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
