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

type StudentHandler struct {
	scheduleService *services.ScheduleService
}

func NewStudentHandler(scheduleService *services.ScheduleService) *StudentHandler {
	return &StudentHandler{scheduleService: scheduleService}
}

func (h *StudentHandler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
		return
	}

	if err := dto.ValidateStruct(req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
	}

	student, err := h.scheduleService.CreateStudent(r.Context(), &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to create student: %v", err))
		return
	}
	utils.SuccessResponse(w, "student created", student)
}

func (h *StudentHandler) GetStudent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid id: %v", err))
		return
	}
	student, err := h.scheduleService.GetStudent(r.Context(), id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get student: %v", err))
		return
	}
	if student == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "student not found")
		return
	}
	utils.SuccessResponse(w, "student retrieved", student)
}

func (h *StudentHandler) GetAllStudents(w http.ResponseWriter, r *http.Request) {
	students, err := h.scheduleService.GetAllStudents(r.Context())
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get students: %v", err))
		return
	}
	utils.SuccessResponse(w, "students retrieved", students)
}

func (h *StudentHandler) UpdateStudentGroup(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid id: %v", err))
		return
	}

	var req struct {
		GroupID int `json:"group_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
		return
	}

	if err := h.scheduleService.UpdateStudentGroup(r.Context(), id, req.GroupID); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to update student group: %v", err))
		return
	}
	utils.SuccessResponse(w, "student group updated", nil)
}

func (h *StudentHandler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid id: %v", err))
		return
	}
	if err := h.scheduleService.DeleteStudent(r.Context(), id); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to delete student: %v", err))
		return
	}
	utils.SuccessResponse(w, "student deleted", nil)
}

// DownloadStudentTemplate скачивает Excel-шаблон для студентов
func (h *StudentHandler) DownloadTemplate(w http.ResponseWriter, r *http.Request) {
    excelSvc := services.NewExcelService()
    data, err := excelSvc.GenerateStudentTemplate()
    if err != nil {
        utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("generate template: %v", err))
        return
    }
    w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
    w.Header().Set("Content-Disposition", "attachment; filename=students_template.xlsx")
    w.Write(data)
}

// UploadStudentsExcel загружает Excel и создаёт студентов
func (h *StudentHandler) UploadExcel(w http.ResponseWriter, r *http.Request) {
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
    
    results, err := h.scheduleService.CreateStudentsFromExcel(r.Context(), file)
    if err != nil {
        utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("import failed: %v", err))
        return
    }
    utils.SuccessResponse(w, "students created from Excel", results)
}
