package handlers

import (
	"core/internal/storage"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type FileHandler struct {
	storage *storage.MemoryStorage
}

func NewFileHandler(strorage *storage.MemoryStorage) *FileHandler {
	return &FileHandler{storage: strorage}
}

func (h *FileHandler) GetFile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "file_id")

	file, ok := h.storage.GetFile(id)
	if !ok {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", file.Mime)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Filename))
	w.Write(file.Data)

}
