package handlers

import (
	"net/http"

	"github.com/bhartiyaanshul/envault/internal/repository"
	"github.com/bhartiyaanshul/envault/internal/server/middleware"
	"github.com/bhartiyaanshul/envault/internal/service"
)

type CreateProjectRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

type ProjectHandler struct {
	projectSvc *service.ProjectService
	userRepo   *repository.UserRepository
}

func NewProjectHandler(projectSvc *service.ProjectService, userRepo *repository.UserRepository) *ProjectHandler {
	return &ProjectHandler{projectSvc: projectSvc, userRepo: userRepo}
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateProjectRequest
	if err := DecodeAndValidate(r, &req); err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	jwtUser := middleware.GetUserFromContext(r.Context())
	if jwtUser == nil {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.userRepo.FindOrCreate(jwtUser.SupabaseUID, jwtUser.Email)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "user lookup failed: "+err.Error())
		return
	}

	project, token, err := h.projectSvc.CreateProject(r.Context(), req.Name, user.ID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"project":     project,
		"vault_token": token,
	})
}

func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	jwtUser := middleware.GetUserFromContext(r.Context())
	if jwtUser == nil {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.userRepo.FindOrCreate(jwtUser.SupabaseUID, jwtUser.Email)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "user lookup failed: "+err.Error())
		return
	}

	projects, err := h.projectSvc.ListProjects(user.ID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, projects)
}

func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	project := middleware.GetProjectFromContext(r.Context())
	if project == nil {
		RespondError(w, http.StatusNotFound, "project not found")
		return
	}
	RespondJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	project := middleware.GetProjectFromContext(r.Context())
	user := middleware.GetUserFromContext(r.Context())
	if project == nil || user == nil {
		RespondError(w, http.StatusNotFound, "not found")
		return
	}

	if err := h.projectSvc.DeleteProject(r.Context(), project.Slug, user.ID); err != nil {
		RespondError(w, http.StatusForbidden, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
