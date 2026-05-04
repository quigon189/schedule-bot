package handlers

import (
	"core/internal/dto"
	"core/internal/llm"
	"core/internal/services"
	"core/pkg/fileparser"
	"core/pkg/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type GroupHandler struct {
	scheduleService *services.ScheduleService
	llmAgent        *llm.LLMAgent
}

func NewGroupHandler(scheduleService *services.ScheduleService, client llm.Client) *GroupHandler {
	return &GroupHandler{
		scheduleService: scheduleService,
		llmAgent:        llm.NewAgent(client, *scheduleService, ""),
	}
}

func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
		return
	}

	if err := dto.ValidateStruct(req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
		return
	}

	group, err := h.scheduleService.CreateGroup(r.Context(), &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to create group: %v", err))
		return
	}
	utils.SuccessResponse(w, "group created", group)
}

func (h *GroupHandler) GetGroup(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid id: %v", err))
		return
	}
	group, err := h.scheduleService.GetGroup(r.Context(), id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get group: %v", err))
		return
	}
	if group == nil {
		utils.ErrorResponse(w, http.StatusNotFound, "group not found")
		return
	}
	utils.SuccessResponse(w, "group retrieved", group)
}

func (h *GroupHandler) GetAllGroups(w http.ResponseWriter, r *http.Request) {
	groupName := r.URL.Query().Get("group_name")

	if groupName == "" {
		groups, err := h.scheduleService.GetAllGroups(r.Context())
		if err != nil {
			utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get groups: %v", err))
			return
		}
		utils.SuccessResponse(w, "groups retrieved", groups)
		return
	}

	group, err := h.scheduleService.GetGroupByName(r.Context(), groupName)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get group: %v", err))
		return
	}
	utils.SuccessResponse(w, "group retrieved", group)
}

func (h *GroupHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid id: %v", err))
		return
	}
	var req dto.UpdateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
		return
	}

	if err := dto.ValidateStruct(req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
		return
	}

	group, err := h.scheduleService.UpdateGroup(r.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to update group: %v", err))
		return
	}

	utils.SuccessResponse(w, "group updated", group)
}

func (h *GroupHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid id: %v", err))
		return
	}
	if err := h.scheduleService.DeleteGroup(r.Context(), id); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to delete group: %v", err))
		return
	}
	utils.SuccessResponse(w, "group deleted", nil)
}

func (h *GroupHandler) CreateGroupWtihCurriculum(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateGroupWithCurriculumRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
		return
	}

	if err := dto.ValidateStruct(req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("validate: %v", err))
		return
	}

	resp, err := h.scheduleService.CreateGroupWithCurriculumAndStudents(r.Context(), &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to create group with curriculum: %v", err))
		return
	}

	utils.SuccessResponse(w, "group with curriculum and students created", resp)
}

func (h *GroupHandler) DownloadTemplate(w http.ResponseWriter, r *http.Request) {
	excelSvc := services.NewExcelService()
	data, err := excelSvc.GenerateGroupTemplate()
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("generate template: %v", err))
		return
	}
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=group_template.xlsx")
	w.Write(data)
}

func (h *GroupHandler) UploadGroupExcel(w http.ResponseWriter, r *http.Request) {
	// Ограничим размер файла (10 MB)
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

	resp, err := h.scheduleService.CreateGroupFromExcel(r.Context(), file)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("import failed: %v", err))
		return
	}

	utils.SuccessResponse(w, "group with curriculum and students created from Excel", resp)
}

func (h *GroupHandler) AIGroupTemplate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "file too large or invalid form")
		return
	}

	files := r.MultipartForm.File["files"]

	var parsedContent []fileparser.ParsedContent

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		pc, err := fileparser.Parse(file)
		if err != nil {
			log.Printf("%s error %v", fileHeader.Filename, err)
		}

		parsedContent = append(parsedContent, *pc)

	}

	var req dto.CreateGroupWithCurriculumRequest

	if err := h.llmAgent.GenerateStructuredData(r.Context(), "Ответь json структурой из системной строки", parsedContent, &req); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "fieled to generate llm response: " + err.Error())
		return
	}

	utils.SuccessResponse(w, "ok", req)
}
