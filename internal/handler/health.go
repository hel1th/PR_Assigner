package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) RegisterRoutes(r chi.Router) {
	r.Get("/health", h.Health)
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"status": "ok"}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
