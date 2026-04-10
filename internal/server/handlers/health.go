package handlers

import (
	"net/http"

	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *HealthHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	sqlDB, err := h.db.DB()
	if err != nil {
		RespondError(w, http.StatusServiceUnavailable, "database unavailable")
		return
	}
	if err := sqlDB.Ping(); err != nil {
		RespondError(w, http.StatusServiceUnavailable, "database ping failed")
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
