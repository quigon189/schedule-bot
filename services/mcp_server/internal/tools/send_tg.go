package tools

import (
	"context"
	"log"
	"mcp_server/internal/config"
	"mcp_server/internal/models"
	"mcp_server/internal/service"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
)

type TGHandler struct {
	tgService       *service.TGService
	scheduleService *service.ChangeService
}

func NewTGHandler(cfg *config.Config) *TGHandler {
	return &TGHandler{
		tgService:       service.NewTGService(cfg),
		scheduleService: service.NewChangeService(cfg),
	}
}

func (h *TGHandler) SendChanges(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	date, err := req.RequireString("date")
	if err != nil {
		log.Printf("Failed to get send changes date: %v", err)
		return mcp.NewToolResultError("ошибка при проверке даты"), nil
	}

	chatIDs, err := req.RequireString("chat_id")
	if err != nil {
		log.Printf("Failed to get send changes chat_id: %v", err)
		return mcp.NewToolResultError("ошибка при проверке chat_id"), nil
	}

	chatID, err := strconv.ParseInt(chatIDs, 10, 64)
	if err != nil {
		log.Printf("Failed to get send changes chat_id: %v", err)
		return mcp.NewToolResultError("ошибка при проверке chat_id"), nil
	}

	changesResponse, err := h.scheduleService.GetChanges(date)
	if err != nil {
		log.Printf("Failed to get changes: %v", err)
		return mcp.NewToolResultError("ошибка при получении изменеий из сервиса"), nil
	}

	if err := h.tgService.SendMessage(&models.SendMessage{
		ChatID:    chatID,
		Message:   changesResponse.Data.Description,
		PhotoURLs: changesResponse.Data.ImgURLs,
	}); err != nil {
		log.Printf("Failed to send tg message: %v", err)
		return mcp.NewToolResultError("ошибка при отправке сообщения в тг"), nil
	}

	return mcp.NewToolResultText("Пользователь получил информацию об изменениях"), nil
}
