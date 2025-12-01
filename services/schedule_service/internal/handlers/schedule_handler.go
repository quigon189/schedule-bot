package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"schedule-service/internal/dto"
	"schedule-service/internal/repository"
	"schedule-service/internal/service"
	"schedule-service/pkg/utils"

	"github.com/gorilla/schema"
)

type ScheduleHandler struct {
	service *service.ScheduleService
}

func NewScheduleHandler(db *repository.DB) *ScheduleHandler{
	repo := repository.NewScheduleRepository(db)
	return &ScheduleHandler{service: service.NewScheduleService(&repo)}
}

func (h *ScheduleHandler) AddGroupSchedule(w http.ResponseWriter, r *http.Request) {
	req := &dto.AddGroupScheduleRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("AddGroupSchedule handler error: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")	
		return
	}

	gs, err := h.service.AddGroupSchedule(req)
	if err != nil {
		log.Printf("AddGroupSchedule handler error: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return 
	}

	resp := dto.GroupScheduleResponse{
		GroupName: gs.Group.Name,
		AcademicYear: gs.StudyPeriod.AcademicYear,
		HalfYear: gs.StudyPeriod.HalfYear,
		Semester: gs.Semester,
		ScheduleImgURL: gs.ScheduleImgURL,
		CreatedAt: gs.CreatedAt,
	}

	utils.SuccessResponse(w, "Group schedule added", resp)
}

func (h *ScheduleHandler) GetGroupSchedule(w http.ResponseWriter, r *http.Request) {
	gsQueryParams := dto.GroupScheduleQueryParams{}
	if err := schema.NewDecoder().Decode(&gsQueryParams, r.URL.Query()); err != nil {
		log.Printf("GetGroupSchedule handler error: %v", err)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid query params")
		return
	}

	gs, err := h.service.GetGroupSchedule(&gsQueryParams)
	if err != nil {
		log.Printf("GetGroupSchedule handler error: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := dto.GroupScheduleResponse{
		GroupName: gs.Group.Name,
		AcademicYear: gs.StudyPeriod.AcademicYear,
		HalfYear: gs.StudyPeriod.HalfYear,
		Semester: gs.Semester,
		ScheduleImgURL: gs.ScheduleImgURL,
		CreatedAt: gs.CreatedAt,
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

	err := h.service.RemoveScheduleService(&gsQueryParams)
	if err != nil {
		log.Printf("RemoveGroupSchedule handler error: %v", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, "Group schedule removed", nil)
}
