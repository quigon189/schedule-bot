package tools

import (
	"context"
	"encoding/json"
	"log"

	"mcp_server/internal/service"

	"github.com/mark3labs/mcp-go/mcp"
)

type ChangeHandler struct {
	service *service.ChangeService
}

func NewChangeHandler(s *service.ChangeService) *ChangeHandler {
	return &ChangeHandler{service: s}
}

func (h *ChangeHandler) HandleChangeRequest(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	date, err := req.RequireString("date")
	if err != nil {
		log.Printf("Failed to get date: %v", err)
		return mcp.NewToolResultError("невалидная дата"), nil
	}

	changesResponse, err := h.service.GetChanges(date)
	if err != nil {
		log.Printf("Failed to get changes: %v", err)
		return mcp.NewToolResultError("ошибка при получении изменеий из сервиса"), nil
	}

	if len(changesResponse) == 0 {
		return mcp.NewToolResultText("нет изменений на выбранную дату"), nil
	}

	jsonData, err := json.Marshal(changesResponse)
	if err != nil {
		log.Printf("Failed to marshal changes response: %v", err)
		return mcp.NewToolResultError("ошибка при обработке данных"), nil
	}

	return mcp.NewToolResultText("Ответ от сервиса баз данных в формате json: "+string(jsonData)), err
}
