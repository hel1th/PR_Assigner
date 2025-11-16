package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hel1th/PR_Assigner/internal/apperrors"
	"github.com/hel1th/PR_Assigner/internal/dto"
)

func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apperrors.NotFound):
		respondError(w, http.StatusNotFound, "NOT_FOUND", "resource not found")

	case errors.Is(err, apperrors.InvalidInput):
		respondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid input")

	case errors.Is(err, apperrors.NotAssigned):
		respondError(w, http.StatusConflict, "NOT_ASSIGNED", "current reviewer is not assigned")

	case errors.Is(err, apperrors.TeamExists):
		respondError(w, http.StatusBadRequest, "TEAM_EXISTS", "team_name already exists")

	case errors.Is(err, apperrors.NoCandidate):
		respondError(w, http.StatusUnprocessableEntity, "NO_CANDIDATE", "no candidates available")

	case errors.Is(err, apperrors.PRExists):
		respondError(w, http.StatusBadRequest, "PR_EXISTS", "pull request already exists")

	case errors.Is(err, apperrors.PRMerged):
		respondError(w, http.StatusBadRequest, "PR_MERGED", "pull request already merged")

	case errors.Is(err, apperrors.ReviewerAlreadyAssigned):
		respondError(w, http.StatusBadRequest, "REVIEWER_ALREADY_ASSIGNED", "reviewer already assigned")

	case errors.Is(err, apperrors.MaxReviewers):
		respondError(w, http.StatusBadRequest, "MAX_REVIEWERS_REACHED", "maximum reviewers limit reached")

	default:
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
	}
}

func respondError(w http.ResponseWriter, httpCode int, errorCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(dto.ErrorResponse{
		Error: dto.ErrorDetail{
			Code:    errorCode,
			Message: message,
		},
	})
}
