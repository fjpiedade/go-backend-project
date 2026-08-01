package main

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	app := newTestApplication(&mockUserStore{})

	rr := makeRequest(t, app, http.MethodGet, "/v1/health", nil)

	assertStatus(t, rr, http.StatusOK)

	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if response["status"] != "available" {
		t.Errorf("expected status %q, got %q", "available", response["status"])
	}

	if response["version"] != "1.0.0" {
		t.Errorf("expected version %q, got %q", "1.0.0", response["version"])
	}
}
