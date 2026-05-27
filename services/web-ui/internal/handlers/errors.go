package handlers

import (
	"log"
	"net/http"
	"web-ui/views/pages"
)

func NotFound(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	if err := pages.NotFoundPage().Render(req.Context(), w); err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func RenderInternalError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal error: %v", err)
	w.WriteHeader(http.StatusInternalServerError)
	if renderErr := pages.InternalServerErrorPage().Render(r.Context(), w); renderErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
