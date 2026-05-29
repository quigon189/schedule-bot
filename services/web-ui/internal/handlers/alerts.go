package handlers

import (
	"net/http"
	"time"
	"web-ui/internal/api"
	"web-ui/internal/models"
	"web-ui/pkg/cache"
	"web-ui/views/components"
)

type AlertsHandler struct {
	alertCache *cache.MemCache[[]models.Alert]
	ttl        time.Duration
}

var alerts AlertsHandler

func NewAlertsHandler(ttl, cleanup time.Duration) *AlertsHandler {
	alertCache := cache.NewMemCache[[]models.Alert](cleanup)
	return &AlertsHandler{alertCache: alertCache, ttl: ttl}
}

func SetupAlertsHandler(ttl, cleanup time.Duration) *AlertsHandler {
	alertCache := cache.NewMemCache[[]models.Alert](cleanup)
	alerts.alertCache = alertCache
	alerts.ttl = ttl

	return &alerts
}

// GET /alerts – возвращает HTML-фрагменты всех непрочитанных алертов и очищает их
func (h *AlertsHandler) GetAlerts(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value("session").(*api.Session)
	if !ok || session.SessionID == "" {
		w.WriteHeader(http.StatusOK)
		return
	}
	alerts, found := h.alertCache.Get(session.SessionID)
	if !found {
		alerts = []models.Alert{}
	} else {
		h.alertCache.Delete(session.SessionID)
	}

	for _, alert := range alerts {
		components.RenderAlert(alert).Render(r.Context(), w)
	}
}

// Вспомогательные методы для добавления алертов (используются в других хэндлерах)
func (h *AlertsHandler) AddSuccess(w http.ResponseWriter, sessionID, message string) {
	w.Header().Set("HX-Trigger-After-Swap", "alertAdded")
	h.addAllert(sessionID, models.AlertSuccess, message)
}

func (h *AlertsHandler) AddError(w http.ResponseWriter, sessionID, message string) {
	w.Header().Set("HX-Trigger-After-Swap", "alertAdded")
	h.addAllert(sessionID, models.AlertError, message)
}

func (h *AlertsHandler) AddInfo(w http.ResponseWriter, sessionID, message string) {
	w.Header().Set("HX-Trigger-After-Swap", "alertAdded")
	h.addAllert(sessionID, models.AlertInfo, message)
}

func (h *AlertsHandler) AddWarning(w http.ResponseWriter, sessionID, message string) {
	w.Header().Set("HX-Trigger-After-Swap", "alertAdded")
	h.addAllert(sessionID, models.AlertWarning, message)
}

func (h *AlertsHandler) addAllert(sessionID string, t models.AlertType, message string) {
	alert := models.Alert{
		Type:    t,
		Message: message,
	}
	alerts, found := h.alertCache.Get(sessionID)
	if !found {
		alerts = []models.Alert{}
	}
	alerts = append(alerts, alert)
	h.alertCache.Set(sessionID, alerts, h.ttl)

}
