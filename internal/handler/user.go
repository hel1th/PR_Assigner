package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hel1th/PR_Assigner/internal/dto"
	"github.com/hel1th/PR_Assigner/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/setIsActive", h.SetIsActive)
		r.Get("/getReview", h.GetReview)
	})
}

func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	var req dto.SetUserActiveRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	if err := h.userService.SetUserActive(r.Context(), req.UserID, req.IsActive); err != nil {
		handleServiceError(w, err)
		return
	}

	user, err := h.userService.GetUser(r.Context(), req.UserID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	resp := dto.SetUserActiveResponse{
		User: dto.UserFromDomain(user),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	if userID == "" {
		respondError(w, http.StatusBadRequest, "INVALID_INPUT", "user_id is required")
		return
	}

	prs, err := h.userService.GetUserReviewPRs(r.Context(), userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	resp := dto.GetUserReviewPRsResponse{
		UserID:       userID,
		PullRequests: dto.PullRequestsShortFromDomain(prs),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
