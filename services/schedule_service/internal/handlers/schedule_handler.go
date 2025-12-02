package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"schedule-service/internal/dto"
	"schedule-service/internal/models"
	"schedule-service/internal/repository"
	"schedule-service/internal/service"
	"schedule-service/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

type ScheduleHandler struct {
	service  *service.ScheduleService
	validate *validator.Validate
}

func NewScheduleHandler(db *repository.DB) *ScheduleHandler {
	repo := repository.NewScheduleRepository(db)
	return &ScheduleHandler{
		service:  service.NewScheduleService(&repo),
		validate: dto.NewScheduleValidator(),
	}
}

func (h *ScheduleHandler) AddGroupSchedule(w http.ResponseWriter, r *http.Request) {
	req := &dto.AddGroupScheduleRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("AddGroupSchedule handler error: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		log.Printf("Invalid request: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %v", err))
		return
	}

	gs, err := h.service.AddGroupSchedule(req)
	if err != nil {
		log.Printf("AddGroupSchedule handler error: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, "Group schedule added", gsToResponse(gs))
}

func (h *ScheduleHandler) GetGroupSchedule(w http.ResponseWriter, r *http.Request) {
	gsQueryParams := dto.GroupScheduleQueryParams{}
	if err := schema.NewDecoder().Decode(&gsQueryParams, r.URL.Query()); err != nil {
		log.Printf("GetGroupSchedule handler error: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid query params")
		return
	}

	if err := h.validate.Struct(gsQueryParams); err != nil {
		log.Printf("Invalid request: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %v", err))
		return
	}

	gss, err := h.service.GetGroupSchedule(&gsQueryParams)
	if err != nil {
		log.Printf("GetGroupSchedule handler error: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := []dto.GroupScheduleResponse{}
	for _, gs := range gss {
		resp = append(resp, gsToResponse(&gs))
	}

	utils.SuccessResponse(w, "Group schedule sent", resp)
}

func (h *ScheduleHandler) RemoveGroupSchedule(w http.ResponseWriter, r *http.Request) {
	gsQueryParams := dto.GroupScheduleQueryParams{}
	if err := schema.NewDecoder().Decode(&gsQueryParams, r.URL.Query()); err != nil {
		log.Printf("GetGroupSchedule handler error: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid query params")
		return
	}

	if gsQueryParams.AcademicYear == "" || gsQueryParams.HalfYear == 0 || gsQueryParams.GroupName == "" {
		log.Print("Invalid request: all field required")
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request: all field required")
		return
	}

	if err := h.validate.Struct(gsQueryParams); err != nil {
		log.Printf("Invalid request: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %v", err))
		return
	}

	err := h.service.RemoveScheduleService(&gsQueryParams)
	if err != nil {
		log.Printf("RemoveGroupSchedule handler error: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, "Group schedule removed", nil)
}

func gsToResponse(gs *models.GroupSchedule) dto.GroupScheduleResponse {
	return dto.GroupScheduleResponse{
		GroupName:      gs.Group.Name,
		AcademicYear:   gs.StudyPeriod.AcademicYear,
		HalfYear:       gs.StudyPeriod.HalfYear,
		Semester:       gs.Semester,
		ScheduleImgURL: gs.ScheduleImgURL,
		CreatedAt:      gs.CreatedAt,
	}
}
