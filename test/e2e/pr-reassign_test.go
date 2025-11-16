package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

const baseURL = "http://localhost:8081"

func TestE2E_Reassign(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	teamPayload := map[string]interface{}{
		"team_name": "frontend",
		"members": []map[string]interface{}{
			{"user_id": "u11", "username": "Alice", "is_active": true},
			{"user_id": "u21", "username": "Bob", "is_active": true},
			{"user_id": "u31", "username": "Charlie", "is_active": true},
			{"user_id": "u41", "username": "Diana", "is_active": true},
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

	body, _ = json.Marshal(teamPayload)
	resp, err = client.Post(baseURL+"/team/add", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected status creating team: %d", resp.StatusCode)
	}

	deactivatePayload := map[string]interface{}{
		"user_id":   "u41",
		"is_active": false,
	}
	body, _ = json.Marshal(deactivatePayload)
	resp, err = client.Post(baseURL+"/users/setIsActive", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("failed to deactivate user: %d", resp.StatusCode)
	}

	prPayload := map[string]interface{}{
		"pull_request_id":   "pr-10011",
		"pull_request_name": "Add search1",
		"author_id":         "u11",
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
	for _, r := range prResp.Pr.AssignedReviewers {
		if r == "u41" {
			t.Fatal("deactivated user should not be assigned")
		}
	}

	reassignPayload := map[string]string{
		"pull_request_id": "pr-10011",
		"old_user_id":     "u41",
	}
	body, _ = json.Marshal(reassignPayload)
	resp, err = client.Post(baseURL+"/pullRequest/reassign", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409 when reassigning inactive or unassigned user, got: %d", resp.StatusCode)
	}
}
