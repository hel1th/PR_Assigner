package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hel1th/PR_Assigner/internal/dto"
	"github.com/hel1th/PR_Assigner/internal/service"
)

type TeamHandler struct {
	teamService service.TeamService
}

func NewTeamHandler(teamService service.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

func (h *TeamHandler) RegisterRoutes(r chi.Router) {
	r.Route("/team", func(r chi.Router) {
		r.Post("/add", h.AddTeam)
		r.Get("/get", h.GetTeam)
	})
}

func (h *TeamHandler) AddTeam(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTeamRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	members := req.ToDomainMembers()

	if err := h.teamService.CreateTeam(r.Context(), req.TeamName, members); err != nil {
		handleServiceError(w, err)
		return
	}

	team, err := h.teamService.GetTeam(r.Context(), req.TeamName)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	resp := dto.CreateTeamResponse{
		Team: dto.TeamFromDomain(team),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")

	if teamName == "" {
		respondError(w, http.StatusBadRequest, "INVALID_INPUT", "team_name is required")
		return
	}

	team, err := h.teamService.GetTeam(r.Context(), teamName)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	
	resp := dto.TeamFromDomain(team)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
