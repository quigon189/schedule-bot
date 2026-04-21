package handlers

import (
	"core/internal/dto"
	"core/internal/services"
	"core/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type AudienceHandler struct {
	scheduleService *services.ScheduleService
}

func NewAudienceHandler(scheduleService *services.ScheduleService) *AudienceHandler {
	return &AudienceHandler{scheduleService: scheduleService}
}

func (h *AudienceHandler) CreateAudience(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAudienceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad requset: %v", err))
		return
	}

	if err := dto.ValidateStruct(req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
		return
	}

	audience, err := h.scheduleService.CreateAudience(r.Context(), &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to create audience: %v", err))
		return
	}

	utils.SuccessResponse(w, "audience created", audience)
}

func (h *AudienceHandler) GetAudience(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
		return
	}

	audience, err := h.scheduleService.GetAudience(r.Context(), id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get audience: %v", err))
		return
	}

	utils.SuccessResponse(w, "audience getted", audience)
}

func (h *AudienceHandler) GetAllAudience(w http.ResponseWriter, r *http.Request) {
	audiences, err := h.scheduleService.GetAllAudience(r.Context())
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get audiences: %v", err))
		return
	}

	utils.SuccessResponse(w, "audiences getted", audiences)
}

func (h *AudienceHandler) UpdateAudience(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid id: %v", err))
		return
	}

	var req dto.UpdateAudienceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad requset: %v", err))
		return
	}

	audience, err := h.scheduleService.UpdateAudience(r.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to update audience: %v", err))
		return
	}

	utils.SuccessResponse(w, "audience updated", audience)
}

func (h *AudienceHandler) DeleteAudience(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad requset: %v", err))
		return
	}

	if err := h.scheduleService.DeleteAudience(r.Context(), id); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to delete audience: %v", err))
		return
	}	

	utils.SuccessResponse(w, "audience deleted", nil)
}

// DownloadTemplate скачивает Excel-шаблон для аудиторий
func (h *AudienceHandler) DownloadTemplate(w http.ResponseWriter, r *http.Request) {
    excelSvc := services.NewExcelService()
    data, err := excelSvc.GenerateAudienceTemplate()
    if err != nil {
        utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("generate template: %v", err))
        return
    }
    w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
    w.Header().Set("Content-Disposition", "attachment; filename=audiences_template.xlsx")
    w.Write(data)
}

// UploadExcel загружает Excel и создаёт аудитории
func (h *AudienceHandler) UploadExcel(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseMultipartForm(10 << 20); err != nil {
        utils.ErrorResponse(w, http.StatusBadRequest, "file too large or invalid form")
        return
    }
    file, _, err := r.FormFile("file")
    if err != nil {
        utils.ErrorResponse(w, http.StatusBadRequest, "missing file field")
        return
    }
    defer file.Close()
    
    audiences, err := h.scheduleService.CreateAudiencesFromExcel(r.Context(), file)
    if err != nil {
        utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("import failed: %v", err))
        return
    }
    utils.SuccessResponse(w, "audiences created from Excel", audiences)
}
