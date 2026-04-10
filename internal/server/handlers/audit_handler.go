package handlers

import (
	"net/http"
	"strconv"

	"github.com/bhartiyaanshul/envault/internal/server/middleware"
	"github.com/bhartiyaanshul/envault/internal/service"
)

type AuditHandler struct {
	auditSvc *service.AuditService
}

func NewAuditHandler(auditSvc *service.AuditService) *AuditHandler {
	return &AuditHandler{auditSvc: auditSvc}
}

func (h *AuditHandler) List(w http.ResponseWriter, r *http.Request) {
	project := middleware.GetProjectFromContext(r.Context())
	if project == nil {
		RespondError(w, http.StatusNotFound, "project not found")
		return
	}

	action := r.URL.Query().Get("action")
	limit := queryParamInt(r, "limit", 50)
	offset := queryParamInt(r, "offset", 0)

	if limit > 100 {
		limit = 100
	}

	logs, total, err := h.auditSvc.ListAuditLogs(project.Slug, action, limit, offset)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, map[string]interface{}{
		"data":   logs,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func queryParamInt(r *http.Request, key string, fallback int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return fallback
	}
	return n
}
