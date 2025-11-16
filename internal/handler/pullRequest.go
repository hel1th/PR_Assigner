package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hel1th/PR_Assigner/internal/dto"
	"github.com/hel1th/PR_Assigner/internal/service"
)

type PullRequestHandler struct {
	prSvc service.PullRequestService
}

func NewPullRequestHandler(prSvc service.PullRequestService) *PullRequestHandler {
	return &PullRequestHandler{prSvc: prSvc}
}

func (h *PullRequestHandler) RegisterRoutes(r chi.Router) {
	r.Route("/pullRequest", func(r chi.Router) {
		r.Post("/create", h.CreatePullReq)
		r.Post("/merge", h.MergePullReq)
		r.Post("/reassign", h.ReassignPullReq)
	})
}

func (h *PullRequestHandler) CreatePullReq(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePullRequestRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	pr, err := h.prSvc.CreatePullReq(r.Context(), req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	resp := dto.CreatePullRequestResponse{PR: dto.PullRequestFromDomain(pr)}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *PullRequestHandler) MergePullReq(w http.ResponseWriter, r *http.Request) {
	var req dto.MergePullRequestRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	pr, err := h.prSvc.MergePullReq(r.Context(), req.PullRequestID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	resp := dto.MergePullRequestResponse{
		PR: dto.PullRequestFromDomain(pr),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *PullRequestHandler) ReassignPullReq(w http.ResponseWriter, r *http.Request) {
	var req dto.ReassignReviewerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	if err := h.prSvc.ReassignReviewer(r.Context(), req.PullRequestID, req.OldUserID); err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "reviewer reassigned successfully",
	})
}
