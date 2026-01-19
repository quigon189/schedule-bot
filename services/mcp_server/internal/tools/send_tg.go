package tools

import (
	"context"
	"fmt"
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

type SendChangesRequest struct {
	ChatID string `json:"chat_id" jsonschema_description:"Идентификатор чата пользователя, необходимый для отправки сообщения. Положительное целое число в формате int64" jsonschema:"required,format=int64"`
	Date   string `json:"date" jsonschema_description:"Дата, на которую необходмо найти и отправить изменения в расписании ГГГГ-ММ-ДД" jsonschema:"required,format=date"`
}

type SendChangesResponse struct {
	Success  bool   `json:"success" jsonschema_description:"Сообщение доставленно пользователю"`
	IsChange bool   `json:"is_change" jsonschema_description:"Есть изменения на указанную дату"`
	Message  string `json:"message" jsonschema_description:"Сообщение MCP клиенту"`
}

func NewTGHandler(cfg *config.Config) *TGHandler {
	return &TGHandler{
		tgService:       service.NewTGService(cfg),
		scheduleService: service.NewChangeService(cfg),
	}
}

func (h *TGHandler) SendChanges(ctx context.Context, req mcp.CallToolRequest, args SendChangesRequest) (SendChangesResponse, error) {
	log.Printf("Send changes request: %+v", args)
	changesResponse, err := h.scheduleService.GetChanges(args.Date)
	if err != nil {
		log.Printf("Failed to get changes: %v", err)
		return SendChangesResponse{
			Success:  false,
			IsChange: false,
			Message:  "Ошибка при получении данных от сервиса расписаний",
		}, nil
	}

	if changesResponse.Data.Date == "" {
		log.Print("Changes not found")
		return SendChangesResponse{
			Success:  false,
			IsChange: false,
			Message:  "Нет изменений на выбранную дату",
		}, nil
	}
	chatID, _ := strconv.ParseInt(args.ChatID, 10, 64)

	if err := h.tgService.SendMessage(&models.SendMessage{
		ChatID:    chatID,
		Message:   changesResponse.Data.Description,
		PhotoURLs: changesResponse.Data.ImgURLs,
	}); err != nil {
		log.Printf("Failed to send tg message: %v", err)
		return SendChangesResponse{
			Success:  false,
			IsChange: true,
			Message:  fmt.Sprintf("Изменения обнаружены, но произошла ошибка при отправке сообщения в телеграмм бот: %v", err),
		}, nil
	}

	return SendChangesResponse{
		Success:  true,
		IsChange: true,
		Message:  "Информация об изменениях отправленна пользователю по указанному chat_id",
	}, nil
}
