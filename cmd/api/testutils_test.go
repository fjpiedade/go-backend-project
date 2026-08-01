package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http/httptest"
	"os"
	"social/internal/store"
	"testing"
)

func newTestApplication(userStore store.UserStorageInterface) *application {
	return &application{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		config: config{addr: ":9090"},
		store: store.Storage{
			Users: userStore,
			Posts: &mockPostStore{},
		},
	}
}

func newTestAppWithPost(postStore store.PostStorageInterface) *application {
	return &application{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		config: config{addr: ":9090"},
		store: store.Storage{
			Users: &mockUserStore{},
			Posts: postStore,
		},
	}
}

func makeRequest(t *testing.T, app *application, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	app.mount().ServeHTTP(rr, req)
	return rr
}

func makeRequestRaw(t *testing.T, app *application, method, path string, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	app.mount().ServeHTTP(rr, req)
	return rr
}

func assertStatus(t *testing.T, rr *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if rr.Code != expected {
		t.Errorf("expected status %d, got %d — body: %s", expected, rr.Code, rr.Body.String())
	}
}
