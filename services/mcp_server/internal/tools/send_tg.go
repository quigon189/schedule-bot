package tools

import (
	"context"
	"log"
	"mcp_server/internal/config"
	"mcp_server/internal/models"
	"mcp_server/internal/service"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

type TGHandler struct {
	tgService       *service.TGService
	scheduleService *service.ChangeService
}

type SendChangesRequest struct {
	ChatID int64     `json:"chat_id" jsonschema_description:"Идентификатор чата пользователя, необходимый для отправки сообщения" jsonschema:"required"`
	Date   time.Time `json:"date" jsonschema_description:"Дата, на которую необходмо найти и отправить изменения в расписании" jsonschema:"required,format=date"`
}

type SendChangesResponse struct {
	Success  bool `json:"success" jsonschema_description:"Сообщение доставленно пользователю"`
	IsChange bool `json:"is_change" jsonschema_description:"Есть изменения на указанную дату"`
	Message string `json:"message" jsonschema_description:"Сообщение MCP клиенту"`
}

func NewTGHandler(cfg *config.Config) *TGHandler {
	return &TGHandler{
		tgService:       service.NewTGService(cfg),
		scheduleService: service.NewChangeService(cfg),
	}
}

func (h *TGHandler) SendChanges(ctx context.Context, req mcp.CallToolRequest, args SendChangesRequest) (SendChangesResponse, error) {
	changesResponse, err := h.scheduleService.GetChanges(args.Date.Format("2005-02-01"))
	if err != nil {
		log.Printf("Failed to get changes: %v", err)
		return SendChangesResponse{
			Success: false,
			IsChange: false,
			Message: "Ошибка при получении данных от сервиса расписаний",
		}, nil
	}

	if changesResponse.Data.Date == "" {
		log.Print("Changes not found")
		return SendChangesResponse{
			Success: false,
			IsChange: false,
			Message: "Нет изменений на выбранную дату",
		}, nil
	}

	if err := h.tgService.SendMessage(&models.SendMessage{
		ChatID:    args.ChatID,
		Message:   changesResponse.Data.Description,
		PhotoURLs: changesResponse.Data.ImgURLs,
	}); err != nil {
		log.Printf("Failed to send tg message: %v", err)
		return SendChangesResponse{
			Success: false,
			IsChange: true,
			Message: "Изменения обнаружены, но произошла ошибка при отправке сообщения в телеграмм бот",
		}, nil
	}

	return SendChangesResponse{
		Success: true,
		IsChange: true,
		Message: "Информация об изменениях отправленна пользователю по указанному chat_id",
	} ,nil
}
