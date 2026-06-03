package handlers

import (
	"io"
	"net/http"
	"sort"
	"strconv"
	"web-ui/internal/api"
	"web-ui/internal/models"
	"web-ui/views/components"
	"web-ui/views/pages"

	"github.com/go-chi/chi/v5"
)

type GroupsHandler struct {
	coreClient *api.CoreClient
}

func NewGroupsHandler(client *api.CoreClient) *GroupsHandler {
	return &GroupsHandler{coreClient: client}
}

// GET /admin/groups – страница со списком групп (карточки)
func (h *GroupsHandler) ListGroupsPage(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)
	user, _ := r.Context().Value("user").(*models.User)

	groups, err := h.coreClient.GetGroups(r.Context(), session)
	if err != nil {
		RenderInternalError(w, r, err)
		return
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})

	// CSRF-токен пока передаём пустым, как в admin_users
	pages.AdminGroupsPage("", user, groups).Render(r.Context(), w)
}

// GET /admin/groups/list – HTMX-фрагмент со списком групп (карточки)
func (h *GroupsHandler) ListGroupsFragment(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)

	groups, err := h.coreClient.GetGroups(r.Context(), session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})

	components.GroupCards(groups).Render(r.Context(), w)
}

// GET /admin/groups/new – модальная форма создания группы
func (h *GroupsHandler) NewGroupForm(w http.ResponseWriter, r *http.Request) {
	components.GroupForm(nil).Render(r.Context(), w)
}

// POST /admin/groups – создание группы
func (h *GroupsHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)

	name := r.FormValue("name")
	specialty := r.FormValue("specialty")
	admissionYearStr := r.FormValue("admission_year")
	admissionYear, err := strconv.Atoi(admissionYearStr)
	if err != nil {
		alerts.AddError(w, session.SessionID, "Год поступления должен быть числом")
		w.Header().Set("HX-Redirect", "/admin/groups")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req := api.CreateGroupRequest{
		Name:          name,
		Specialty:     specialty,
		AdmissionYear: admissionYear,
	}
	_, err = h.coreClient.CreateGroup(r.Context(), session, req)
	if err != nil {
		alerts.AddError(w, session.SessionID, "Ошибка создания: "+err.Error())
		w.Header().Set("HX-Redirect", "/admin/groups")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	alerts.AddSuccess(w, session.SessionID, "Группа успешно создана")
	w.Header().Set("HX-Redirect", "/admin/groups")
	w.WriteHeader(http.StatusOK)
}

// DELETE /admin/groups/{id} – удаление группы
func (h *GroupsHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = h.coreClient.DeleteGroup(r.Context(), session, id)
	if err != nil {
		alerts.AddError(w, session.SessionID, "Ошибка удаления: "+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	alerts.AddSuccess(w, session.SessionID, "Группа удалена")
	w.WriteHeader(http.StatusOK)
}

// GET /admin/groups/upload-form – модальная форма загрузки Excel
func (h *GroupsHandler) UploadGroupsForm(w http.ResponseWriter, r *http.Request) {
	components.GroupUploadForm().Render(r.Context(), w)
}

// GET /admin/groups/template – скачать шаблон Excel
func (h *GroupsHandler) DownloadGroupTemplate(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)
	data, err := h.coreClient.DownloadGroupTemplate(r.Context(), session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=groups_template.xlsx")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Write(data)
}

// POST /admin/groups/upload – загрузка Excel с группами
func (h *GroupsHandler) UploadGroupsExcel(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)

	file, header, err := r.FormFile("file")
	if err != nil {
		alerts.AddError(w, session.SessionID, "Файл не загружен")
		w.Header().Set("HX-Redirect", "/admin/groups")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileData := make([]byte, header.Size)
	_, err = file.Read(fileData)
	if err != nil {
		alerts.AddError(w, session.SessionID, "Ошибка чтения файла")
		w.Header().Set("HX-Redirect", "/admin/groups")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := h.coreClient.UploadGroupExcel(r.Context(), session, fileData, header.Filename)
	if err != nil {
		alerts.AddError(w, session.SessionID, "Ошибка загрузки: "+err.Error())
	} else {
		alerts.AddSuccess(w, session.SessionID, "Группы успешно загружены")
		components.GroupUploadResults(result.Students).Render(r.Context(), w)
		return
	}
	w.Header().Set("HX-Redirect", "/admin/groups")
	w.WriteHeader(http.StatusOK)
}

// файл: internal/handlers/groups.go (добавить новый метод)

// POST /admin/groups/ai-template – генерация шаблона через AI
func (h *GroupsHandler) GenerateAITemplate(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*api.Session)
	
	message := r.FormValue("message")
	
	// Получаем файлы
	var files []api.File
	if err := r.ParseMultipartForm(32 << 20); err != nil { // maxMemory 32MB
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	formFiles := r.MultipartForm.File["files"]
	for _, fileHeader := range formFiles {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Failed to open file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		
		fileData, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to read file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		
		files = append(files, api.File{
			Filename: fileHeader.Filename,
			FileData: fileData,
		})
	}
	
	// Вызываем API для генерации шаблона
	fileData, err := h.coreClient.DownloadAIGroupTemplate(r.Context(), session, message, files)
	if err != nil {
		http.Error(w, "AI generation failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Отдаём файл пользователю
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=ai_generated_groups_template.xlsx")
	w.Header().Set("Content-Length", strconv.Itoa(len(fileData)))
	w.Write(fileData)
}
