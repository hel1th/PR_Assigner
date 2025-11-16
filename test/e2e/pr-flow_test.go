package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestE2E_PRFlow(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	teamPayload := map[string]interface{}{
		"team_name": "backend",
		"members": []map[string]interface{}{
			{"user_id": "u1", "username": "Alice", "is_active": true},
			{"user_id": "u2", "username": "Bob", "is_active": true},
		},
	}
	body, _ := json.Marshal(teamPayload)
	resp, err := client.Post(baseURL+"/team/add", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status creating team: %d", resp.StatusCode)
	}

	prPayload := map[string]interface{}{
		"pull_request_id":   "pr-1001",
		"pull_request_name": "Add search",
		"author_id":         "u1",
	}
	body, _ = json.Marshal(prPayload)
	resp, err = client.Post(baseURL+"/pullRequest/create", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected status creating PR: %d", resp.StatusCode)
	}

	var prResp struct {
		Pr struct {
			PullRequestID     string   `json:"pull_request_id"`
			AssignedReviewers []string `json:"assigned_reviewers"`
		} `json:"pr"`
	}
	json.NewDecoder(resp.Body).Decode(&prResp)
	if len(prResp.Pr.AssignedReviewers) == 0 {
		t.Fatal("no reviewers assigned")
	}

	req, _ := http.NewRequest("GET", baseURL+"/users/getReview?user_id=u2", nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status fetching user reviews: %d", resp.StatusCode)
	}

	var reviewsResp struct {
		UserID       string `json:"user_id"`
		PullRequests []struct {
			PullRequestID string `json:"pull_request_id"`
			Status        string `json:"status"`
		} `json:"pull_requests"`
	}
	json.NewDecoder(resp.Body).Decode(&reviewsResp)
	if len(reviewsResp.PullRequests) == 0 {
		t.Fatal("no PRs returned for reviewer")
	}

	mergePayload := map[string]string{"pull_request_id": "pr-1001"}
	body, _ = json.Marshal(mergePayload)
	resp, err = client.Post(baseURL+"/pullRequest/merge", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status merging PR: %d", resp.StatusCode)
	}
}
