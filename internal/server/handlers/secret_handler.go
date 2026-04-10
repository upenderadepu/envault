package handlers

import (
	"net/http"

	"github.com/bhartiyaanshul/envault/internal/repository"
	"github.com/bhartiyaanshul/envault/internal/server/middleware"
	"github.com/bhartiyaanshul/envault/internal/service"
	"github.com/go-chi/chi/v5"
)

type SetSecretRequest struct {
	Environment string `json:"environment" validate:"required,oneof=development staging production"`
	Key         string `json:"key" validate:"required,min=1,max=256,keyname"`
	Value       string `json:"value" validate:"required"`
}

type BulkSetSecretsRequest struct {
	Environment string            `json:"environment" validate:"required,oneof=development staging production"`
	Secrets     map[string]string `json:"secrets" validate:"required,min=1"`
}

type SecretHandler struct {
	secretSvc *service.SecretService
	userRepo  *repository.UserRepository
}

func NewSecretHandler(secretSvc *service.SecretService, userRepo *repository.UserRepository) *SecretHandler {
	return &SecretHandler{secretSvc: secretSvc, userRepo: userRepo}
}

// List returns secret metadata only — NO values.
func (h *SecretHandler) List(w http.ResponseWriter, r *http.Request) {
	project := middleware.GetProjectFromContext(r.Context())
	if project == nil {
		RespondError(w, http.StatusNotFound, "project not found")
		return
	}

	envName := r.URL.Query().Get("environment")
	if envName == "" {
		envName = "development"
	}

	metas, err := h.secretSvc.ListKeys(project.Slug, envName)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, metas)
}

// Get is the ONLY endpoint that returns a secret value.
func (h *SecretHandler) Get(w http.ResponseWriter, r *http.Request) {
	project := middleware.GetProjectFromContext(r.Context())
	user := middleware.GetUserFromContext(r.Context())
	if project == nil || user == nil {
		RespondError(w, http.StatusNotFound, "not found")
		return
	}

	keyName := chi.URLParam(r, "key")
	envName := r.URL.Query().Get("environment")
	if envName == "" {
		envName = "development"
	}

	value, meta, err := h.secretSvc.GetSecret(r.Context(), project.Slug, envName, keyName, user.ID)
	if err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, map[string]interface{}{
		"key":              meta.KeyName,
		"value":            value,
		"version":          meta.VaultVersion,
		"last_modified_at": meta.LastModifiedAt,
	})
}

func (h *SecretHandler) Set(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	project := middleware.GetProjectFromContext(r.Context())
	if user == nil || project == nil {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req SetSecretRequest
	if err := DecodeAndValidate(r, &req); err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	meta, err := h.secretSvc.SetSecret(r.Context(), project.Slug, req.Environment, req.Key, req.Value, user.ID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, meta)
}

func (h *SecretHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	project := middleware.GetProjectFromContext(r.Context())
	if user == nil || project == nil {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	keyName := chi.URLParam(r, "key")
	envName := r.URL.Query().Get("environment")
	if envName == "" {
		envName = "development"
	}

	if err := h.secretSvc.DeleteSecret(r.Context(), project.Slug, envName, keyName, user.ID); err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SecretHandler) BulkSet(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	project := middleware.GetProjectFromContext(r.Context())
	if user == nil || project == nil {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req BulkSetSecretsRequest
	if err := DecodeAndValidate(r, &req); err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	metas, err := h.secretSvc.BulkSetSecrets(r.Context(), project.Slug, req.Environment, req.Secrets, user.ID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, metas)
}
