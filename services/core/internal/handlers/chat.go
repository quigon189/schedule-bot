package handlers

import (
	"encoding/json"
	"net/http"

	"core/internal/llm"
	"core/internal/services"
	"core/pkg/utils"
)

type ChatHandler struct {
	agent *llm.LLMAgent
}

func NewChatHandler(client llm.Client, svc services.ScheduleService) *ChatHandler {
	systemPromt := `Ты — ассистент по расписанию занятий. Твоя задача — понять, какую информацию хочет получить пользователь, выбрать подходящий инструмент из списка и дать пользователю короткий ответ.
Отвечай только на вопросы связанные с расписанием заняти.`
	agent := llm.NewAgent(client, svc, systemPromt)
	return &ChatHandler{agent: agent}
}

type ChatRequest struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Reply string `json:"reply"`
}

func (h *ChatHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}
	
	if req.Message == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "message is required")
		return
	}
	
	reply, err := h.agent.ProcessMessage(r.Context(), req.Message)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "processing failed: "+err.Error())
		return
	}
	
	utils.SuccessResponse(w, "ok", ChatResponse{Reply: reply})
}
