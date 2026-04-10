package handlers

import (
	"net/http"

	"github.com/bhartiyaanshul/envault/internal/server/middleware"
	"github.com/bhartiyaanshul/envault/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AddMemberRequest struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=admin developer ci"`
}

type MemberHandler struct {
	memberSvc *service.MemberService
}

func NewMemberHandler(memberSvc *service.MemberService) *MemberHandler {
	return &MemberHandler{memberSvc: memberSvc}
}

func (h *MemberHandler) List(w http.ResponseWriter, r *http.Request) {
	project := middleware.GetProjectFromContext(r.Context())
	if project == nil {
		RespondError(w, http.StatusNotFound, "project not found")
		return
	}

	members, err := h.memberSvc.ListMembers(project.Slug)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, members)
}

func (h *MemberHandler) Add(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	project := middleware.GetProjectFromContext(r.Context())
	if user == nil || project == nil {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req AddMemberRequest
	if err := DecodeAndValidate(r, &req); err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	member, token, err := h.memberSvc.AddMember(r.Context(), project.Slug, req.Email, req.Role, user.ID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"member":      member,
		"vault_token": token,
	})
}

func (h *MemberHandler) Remove(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	project := middleware.GetProjectFromContext(r.Context())
	if user == nil || project == nil {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	memberIDStr := chi.URLParam(r, "id")
	memberID, err := uuid.Parse(memberIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid member ID")
		return
	}

	if err := h.memberSvc.RemoveMember(r.Context(), project.Slug, memberID, user.ID); err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MemberHandler) Rotate(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	project := middleware.GetProjectFromContext(r.Context())
	if user == nil || project == nil {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	token, err := h.memberSvc.RotateCredentials(r.Context(), project.Slug, user.ID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"vault_token": token})
}
