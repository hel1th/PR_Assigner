package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hel1th/PR_Assigner/internal/dto"
	"github.com/hel1th/PR_Assigner/internal/service"
)

type StatsHandler struct {
	statsService service.StatsService
}

func NewStatsHandler(statsService service.StatsService) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
	}
}

func (h *StatsHandler) RegisterRoutes(r chi.Router) {
	r.Route("/stats", func(r chi.Router) {
		r.Get("/users", h.GetUserStats)
	})
}

func (h *StatsHandler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	var teamNamePtr *string
	if teamName != "" {
		teamNamePtr = &teamName
	}

	stats, err := h.statsService.GetUserStats(r.Context(), teamNamePtr)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	resp := dto.GetUserStatsResponse{
		Stats: dto.UserStatsFromDomain(stats),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
