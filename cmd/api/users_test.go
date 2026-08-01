// package main

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"social/internal/store"
// 	"testing"
// 	"time"
// )

// // ─── Helpers ────────────────────────────────────────────────────────────────

// func makeRequest(t *testing.T, app *application, method, path string, body any) *httptest.ResponseRecorder {
// 	t.Helper()

// 	var buf bytes.Buffer
// 	if body != nil {
// 		if err := json.NewEncoder(&buf).Encode(body); err != nil {
// 			t.Fatalf("failed to encode body: %v", err)
// 		}
// 	}

// 	req := httptest.NewRequest(method, path, &buf)
// 	req.Header.Set("Content-Type", "application/json")
// 	rr := httptest.NewRecorder()
// 	app.mount().ServeHTTP(rr, req)
// 	return rr
// }

// func assertStatus(t *testing.T, rr *httptest.ResponseRecorder, want int) {
// 	t.Helper()
// 	if rr.Code != want {
// 		t.Errorf("expected status %d, got %d — body: %s", want, rr.Code, rr.Body.String())
// 	}
// }

// func assertBodyContains(t *testing.T, rr *httptest.ResponseRecorder, key, value string) {
// 	t.Helper()
// 	var result map[string]any
// 	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
// 		t.Fatalf("failed to decode response body: %v", err)
// 	}
// 	got, ok := result[key]
// 	if !ok {
// 		t.Errorf("key %q not found in response body", key)
// 		return
// 	}
// 	if got != value {
// 		t.Errorf("expected %q=%q, got %q", key, value, got)
// 	}
// }

// func assertBodyNotContains(t *testing.T, rr *httptest.ResponseRecorder, key string) {
// 	t.Helper()
// 	var result map[string]any
// 	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
// 		t.Fatalf("failed to decode response body: %v", err)
// 	}
// 	if _, ok := result[key]; ok {
// 		t.Errorf("key %q should NOT be present in response body", key)
// 	}
// }

// // ─── Fixtures ───────────────────────────────────────────────────────────────

// func fakeUser() *store.User {
// 	return &store.User{
// 		ID:        1,
// 		FirstName: "Ensei",
// 		LastName:  "Tankado",
// 		Username:  "ensei.tankado",
// 		Email:     "ensei@example.com",
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 	}
// }

// // ─── POST /v1/users ──────────────────────────────────────────────────────────

// func TestCreateUserHandler_Success(t *testing.T) {
// 	app := newTestApplication(&mockUserStore{
// 		createFn: func(ctx context.Context, u *store.User) error {
// 			u.ID = 1
// 			u.CreatedAt = time.Now()
// 			u.UpdatedAt = time.Now()
// 			return nil
// 		},
// 	})

// 	body := map[string]string{
// 		"first_name": "Ensei",
// 		"last_name":  "Tankado",
// 		"username":   "ensei.tankado",
// 		"email":      "ensei@example.com",
// 		"password":   "secret123",
// 	}

// 	rr := makeRequest(t, app, http.MethodPost, "/v1/users", body)
// 	assertStatus(t, rr, http.StatusCreated)
// }

// func TestCreateUserHandler_PasswordNotExposed(t *testing.T) {
// 	app := newTestApplication(&mockUserStore{
// 		createFn: func(ctx context.Context, u *store.User) error {
// 			u.ID = 1
// 			u.CreatedAt = time.Now()
// 			u.UpdatedAt = time.Now()
// 			return nil
// 		},
// 	})

// 	body := map[string]string{
// 		"first_name": "Ensei",
// 		"last_name":  "Tankado",
// 		"username":   "ensei.tankado",
// 		"email":      "ensei@example.com",
// 		"password":   "secret123",
// 	}

// 	rr := makeRequest(t, app, http.MethodPost, "/v1/users", body)
// 	assertStatus(t, rr, http.StatusCreated)
// 	assertBodyNotContains(t, rr, "password")
// }

// func TestCreateUserHandler_MissingFields(t *testing.T) {
// 	app := newTestApplication(&mockUserStore{})

// 	body := map[string]string{
// 		"first_name": "Ensei",
// 		// faltam os outros campos
// 	}

// 	rr := makeRequest(t, app, http.MethodPost, "/v1/users", body)
// 	assertStatus(t, rr, http.StatusBadRequest)
// }

// func TestCreateUserHandler_InvalidJSON(t *testing.T) {
// 	app := newTestApplication(&mockUserStore{})

// 	req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewBufferString(`{invalid json}`))
// 	req.Header.Set("Content-Type", "application/json")
// 	rr := httptest.NewRecorder()
// 	app.mount().ServeHTTP(rr, req)

// 	assertStatus(t, rr, http.StatusBadRequest)
// }

// func TestCreateUserHandler_Conflict(t *testing.T) {
// 	app := newTestApplication(&mockUserStore{
// 		createFn: func(ctx context.Context, u *store.User) error {
// 			return store.ErrConflict
// 		},
// 	})

// 	body := map[string]string{
// 		"first_name": "Ensei",
// 		"last_name":  "Tankado",
// 		"username":   "ensei.tankado",
// 		"email":      "ensei@example.com",
// 		"password":   "secret123",
// 	}

// 	rr := makeRequest(t, app, http.MethodPost, "/v1/users", body)
// 	assertStatus(t, rr, http.StatusConflict)
// }

// // ─── GET /v1/users ───────────────────────────────────────────────────────────

// func TestListUsersHandler_Success(t *testing.T) {
// 	app := newTestApplication(&mockUserStore{
// 		getAllFn: func(ctx context.Context) ([]store.User, error) {
// 			return []store.User{*fakeUser()}, nil
// 		},
// 	})

// 	rr := makeRequest(t, app, http.MethodGet, "/v1/users", nil)
// 	assertStatus(t, rr, http.StatusOK)
// }

// func TestListUsersHandler_Empty(t *testing.T) {
// 	app := newTestApplication(&mockUserStore{
// 		getAllFn: func(ctx context.Context) ([]store.User, error) {
// 			return []store.User{}, nil
// 		},
// 	})

// 	rr := makeRequest(t, app, http.MethodGet, "/v1/users", nil)
// 	assertStatus(t, rr, http.StatusOK)
// }

// // ─── GET /v1/users/{id} ──────────────────────────────────────────────────────

// func TestGetUserHandler_Success(t *testing.T) {
// 	app := newTestApplication(&mockUserStore{
// 		getByIDFn: func(ctx context.Context, id int64) (*store.User, error) {
// 			return fakeUser(), nil
// 		},
// 	})

// 	rr := makeRequest(t, app, http.MethodGet, "/v1/users/1", nil)
// 	assertStatus(t, rr, http.StatusOK)
// }

// func TestGetUserHandler_NotFound(t *testing.T) {
// 	app := newTestApplication(&mockUserStore{
// 		getByIDFn: func(ctx context.Context, id int64) (*store.User, error) {
// 			return nil, store.ErrNotFound
// 		},
// 	})

// 	rr := makeRequest(t, app, http.MethodGet, "/v1/users/999", nil)
// 	assertStatus(t, rr, http.StatusNotFound)
// }

// func TestGetUserHandler_InvalidID(t *testing.T) {
// 	app := newTestApplication(&mockUserStore{})

// 	rr := makeRequest(t, app, http.MethodGet, "/v1/users/abc", nil)
// 	assertStatus(t, rr, http.StatusBadRequest)
// }

// // ─── DELETE /v1/users/{id} ───────────────────────────────────────────────────

// func TestDeleteUserHandler_Success(t *testing.T) {
// 	app := newTestApplication(&mockUserStore{
// 		deleteFn: func(ctx context.Context, id int64) error {
// 			return nil
// 		},
// 	})

// 	rr := makeRequest(t, app, http.MethodDelete, "/v1/users/1", nil)
// 	assertStatus(t, rr, http.StatusNoContent)
// }

// func TestDeleteUserHandler_NotFound(t *testing.T) {
// 	app := newTestApplication(&mockUserStore{
// 		deleteFn: func(ctx context.Context, id int64) error {
// 			return store.ErrNotFound
// 		},
// 	})

// 	rr := makeRequest(t, app, http.MethodDelete, "/v1/users/999", nil)
// 	assertStatus(t, rr, http.StatusNotFound)
// }

// // ─── PATCH /v1/users/{id} ────────────────────────────────────────────────────

// func TestUpdateUserHandler_Success(t *testing.T) {
// 	app := newTestApplication(&mockUserStore{
// 		getByIDFn: func(ctx context.Context, id int64) (*store.User, error) {
// 			return fakeUser(), nil
// 		},
// 		updateFn: func(ctx context.Context, u *store.User) error {
// 			return nil
// 		},
// 	})

// 	body := map[string]string{"first_name": "Updated"}
// 	rr := makeRequest(t, app, http.MethodPatch, "/v1/users/1", body)
// 	assertStatus(t, rr, http.StatusOK)
// }

// func TestUpdateUserHandler_NotFound(t *testing.T) {
// 	app := newTestApplication(&mockUserStore{
// 		getByIDFn: func(ctx context.Context, id int64) (*store.User, error) {
// 			return nil, store.ErrNotFound
// 		},
// 	})

// 	body := map[string]string{"first_name": "Updated"}
// 	rr := makeRequest(t, app, http.MethodPatch, "/v1/users/999", body)
// 	assertStatus(t, rr, http.StatusNotFound)
// }

// func TestUpdateUserHandler_Conflict(t *testing.T) {
// 	app := newTestApplication(&mockUserStore{
// 		getByIDFn: func(ctx context.Context, id int64) (*store.User, error) {
// 			return fakeUser(), nil
// 		},
// 		updateFn: func(ctx context.Context, user *store.User) error {
// 			return store.ErrConflict
// 		},
// 	})

// 	body := map[string]string{
// 		"username": "existing.username",
// 	}

// 	rr := makeRequest(t, app, http.MethodPatch, "/v1/users/1", body)

// 	assertStatus(t, rr, http.StatusConflict)
// }

// func TestCreateUserHandler_InternalServerError(t *testing.T) {
// 	app := newTestApplication(&mockUserStore{
// 		createFn: func(ctx context.Context, user *store.User) error {
// 			return errors.New("database unavailable")
// 		},
// 	})

// 	body := map[string]string{
// 		"first_name": "Ensei",
// 		"last_name":  "Tankado",
// 		"username":   "ensei.tankado",
// 		"email":      "ensei@example.com",
// 		"password":   "secret123",
// 	}

// 	rr := makeRequest(t, app, http.MethodPost, "/v1/users", body)

// 	assertStatus(t, rr, http.StatusInternalServerError)
// }

// // ─── Casos de erro interno em falta ─────────────────────────────────────────

// func TestListUsersHandler_InternalServerError(t *testing.T) {
// 	app := newTestApplication(&mockUserStore{
// 		getAllFn: func(ctx context.Context) ([]store.User, error) {
// 			return nil, errors.New("db down")
// 		},
// 	})

// 	rr := makeRequest(t, app, http.MethodGet, "/v1/users", nil)
// 	assertStatus(t, rr, http.StatusInternalServerError)
// }

// func TestGetUserHandler_InternalServerError(t *testing.T) {
// 	app := newTestApplication(&mockUserStore{
// 		getByIDFn: func(ctx context.Context, id int64) (*store.User, error) {
// 			return nil, errors.New("db down")
// 		},
// 	})

// 	rr := makeRequest(t, app, http.MethodGet, "/v1/users/1", nil)
// 	assertStatus(t, rr, http.StatusInternalServerError)
// }

// func TestDeleteUserHandler_InvalidID(t *testing.T) {
// 	app := newTestApplication(&mockUserStore{})

// 	rr := makeRequest(t, app, http.MethodDelete, "/v1/users/abc", nil)
// 	assertStatus(t, rr, http.StatusBadRequest)
// }

// func TestDeleteUserHandler_InternalServerError(t *testing.T) {
// 	app := newTestApplication(&mockUserStore{
// 		deleteFn: func(ctx context.Context, id int64) error {
// 			return errors.New("db down")
// 		},
// 	})

// 	rr := makeRequest(t, app, http.MethodDelete, "/v1/users/1", nil)
// 	assertStatus(t, rr, http.StatusInternalServerError)
// }

// func TestUpdateUserHandler_InvalidID(t *testing.T) {
// 	app := newTestApplication(&mockUserStore{})

// 	body := map[string]string{"first_name": "Updated"}
// 	rr := makeRequest(t, app, http.MethodPatch, "/v1/users/abc", body)
// 	assertStatus(t, rr, http.StatusBadRequest)
// }

// func TestUpdateUserHandler_GetByIDInternalError(t *testing.T) {
// 	app := newTestApplication(&mockUserStore{
// 		getByIDFn: func(ctx context.Context, id int64) (*store.User, error) {
// 			return nil, errors.New("db down")
// 		},
// 	})

// 	body := map[string]string{"first_name": "Updated"}
// 	rr := makeRequest(t, app, http.MethodPatch, "/v1/users/1", body)
// 	assertStatus(t, rr, http.StatusInternalServerError)
// }

// func TestUpdateUserHandler_UpdateInternalError(t *testing.T) {
// 	app := newTestApplication(&mockUserStore{
// 		getByIDFn: func(ctx context.Context, id int64) (*store.User, error) {
// 			return fakeUser(), nil
// 		},
// 		updateFn: func(ctx context.Context, u *store.User) error {
// 			return errors.New("db down")
// 		},
// 	})

// 	body := map[string]string{"first_name": "Updated"}
// 	rr := makeRequest(t, app, http.MethodPatch, "/v1/users/1", body)
// 	assertStatus(t, rr, http.StatusInternalServerError)
// }

// func TestCreateUserHandler_PasswordTooLong(t *testing.T) {
// 	app := newTestApplication(&mockUserStore{})

// 	body := map[string]string{
// 		"first_name": "Ensei",
// 		"last_name":  "Tankado",
// 		"username":   "ensei.tankado",
// 		"email":      "ensei@example.com",
// 		"password":   string(make([]byte, 73)), // > 72 bytes → bcrypt falha
// 	}

// 	rr := makeRequest(t, app, http.MethodPost, "/v1/users", body)
// 	assertStatus(t, rr, http.StatusInternalServerError)
// }

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"social/internal/store"
	"testing"
	"time"
)

// ─── Fixtures ────

func fakeUser() *store.User {
	return &store.User{
		ID:        1,
		FirstName: "Ensei",
		LastName:  "Tankado",
		Username:  "ensei.tankado",
		Email:     "ensei@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func assertBodyContains(t *testing.T, rr *httptest.ResponseRecorder, key, value string) {
	t.Helper()
	var result map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	got, ok := result[key]
	if !ok {
		t.Errorf("key %q not found in response body", key)
		return
	}
	if got != value {
		t.Errorf("expected %q=%q, got %q", key, value, got)
	}
}

func assertBodyNotContains(t *testing.T, rr *httptest.ResponseRecorder, key string) {
	t.Helper()
	var result map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if _, ok := result[key]; ok {
		t.Errorf("key %q should NOT be present in response body", key)
	}
}

// ─── POST /v1/users ────

func TestCreateUserHandler_Success(t *testing.T) {
	app := newTestApplication(&mockUserStore{
		createFn: func(ctx context.Context, u *store.User) error {
			u.ID = 1
			u.CreatedAt = time.Now()
			u.UpdatedAt = time.Now()
			return nil
		},
	})

	body := map[string]string{
		"first_name": "Ensei",
		"last_name":  "Tankado",
		"username":   "ensei.tankado",
		"email":      "ensei@example.com",
		"password":   "secret123",
	}

	rr := makeRequest(t, app, http.MethodPost, "/v1/users", body)
	assertStatus(t, rr, http.StatusCreated)
}

func TestCreateUserHandler_PasswordNotExposed(t *testing.T) {
	app := newTestApplication(&mockUserStore{
		createFn: func(ctx context.Context, u *store.User) error {
			u.ID = 1
			u.CreatedAt = time.Now()
			u.UpdatedAt = time.Now()
			return nil
		},
	})

	body := map[string]string{
		"first_name": "Ensei",
		"last_name":  "Tankado",
		"username":   "ensei.tankado",
		"email":      "ensei@example.com",
		"password":   "secret123",
	}

	rr := makeRequest(t, app, http.MethodPost, "/v1/users", body)
	assertStatus(t, rr, http.StatusCreated)
	assertBodyNotContains(t, rr, "password")
}

func TestCreateUserHandler_MissingFields(t *testing.T) {
	app := newTestApplication(&mockUserStore{})

	body := map[string]string{"first_name": "Ensei"}
	rr := makeRequest(t, app, http.MethodPost, "/v1/users", body)
	assertStatus(t, rr, http.StatusBadRequest)
}

func TestCreateUserHandler_InvalidJSON(t *testing.T) {
	app := newTestApplication(&mockUserStore{})

	req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewBufferString(`{invalid json}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	app.mount().ServeHTTP(rr, req)

	assertStatus(t, rr, http.StatusBadRequest)
}

func TestCreateUserHandler_Conflict(t *testing.T) {
	app := newTestApplication(&mockUserStore{
		createFn: func(ctx context.Context, u *store.User) error {
			return store.ErrConflict
		},
	})

	body := map[string]string{
		"first_name": "Ensei",
		"last_name":  "Tankado",
		"username":   "ensei.tankado",
		"email":      "ensei@example.com",
		"password":   "secret123",
	}

	rr := makeRequest(t, app, http.MethodPost, "/v1/users", body)
	assertStatus(t, rr, http.StatusConflict)
}

// ─── GET /v1/users ────

func TestListUsersHandler_Success(t *testing.T) {
	app := newTestApplication(&mockUserStore{
		getAllFn: func(ctx context.Context) ([]store.User, error) {
			return []store.User{*fakeUser()}, nil
		},
	})

	rr := makeRequest(t, app, http.MethodGet, "/v1/users", nil)
	assertStatus(t, rr, http.StatusOK)
}

func TestListUsersHandler_Empty(t *testing.T) {
	app := newTestApplication(&mockUserStore{
		getAllFn: func(ctx context.Context) ([]store.User, error) {
			return []store.User{}, nil
		},
	})

	rr := makeRequest(t, app, http.MethodGet, "/v1/users", nil)
	assertStatus(t, rr, http.StatusOK)
}

// ─── GET /v1/users/{id} ────

func TestGetUserHandler_Success(t *testing.T) {
	app := newTestApplication(&mockUserStore{
		getByIDFn: func(ctx context.Context, id int64) (*store.User, error) {
			return fakeUser(), nil
		},
	})

	rr := makeRequest(t, app, http.MethodGet, "/v1/users/1", nil)
	assertStatus(t, rr, http.StatusOK)
}

func TestGetUserHandler_NotFound(t *testing.T) {
	app := newTestApplication(&mockUserStore{
		getByIDFn: func(ctx context.Context, id int64) (*store.User, error) {
			return nil, store.ErrNotFound
		},
	})

	rr := makeRequest(t, app, http.MethodGet, "/v1/users/999", nil)
	assertStatus(t, rr, http.StatusNotFound)
}

func TestGetUserHandler_InvalidID(t *testing.T) {
	app := newTestApplication(&mockUserStore{})

	rr := makeRequest(t, app, http.MethodGet, "/v1/users/abc", nil)
	assertStatus(t, rr, http.StatusBadRequest)
}

// ─── DELETE /v1/users/{id} ────

func TestDeleteUserHandler_Success(t *testing.T) {
	app := newTestApplication(&mockUserStore{
		deleteFn: func(ctx context.Context, id int64) error { return nil },
	})

	rr := makeRequest(t, app, http.MethodDelete, "/v1/users/1", nil)
	assertStatus(t, rr, http.StatusNoContent)
}

func TestDeleteUserHandler_NotFound(t *testing.T) {
	app := newTestApplication(&mockUserStore{
		deleteFn: func(ctx context.Context, id int64) error { return store.ErrNotFound },
	})

	rr := makeRequest(t, app, http.MethodDelete, "/v1/users/999", nil)
	assertStatus(t, rr, http.StatusNotFound)
}

func TestDeleteUserHandler_InvalidID(t *testing.T) {
	app := newTestApplication(&mockUserStore{})

	rr := makeRequest(t, app, http.MethodDelete, "/v1/users/abc", nil)
	assertStatus(t, rr, http.StatusBadRequest)
}

func TestDeleteUserHandler_InternalServerError(t *testing.T) {
	app := newTestApplication(&mockUserStore{
		deleteFn: func(ctx context.Context, id int64) error { return errors.New("db down") },
	})

	rr := makeRequest(t, app, http.MethodDelete, "/v1/users/1", nil)
	assertStatus(t, rr, http.StatusInternalServerError)
}

// ─── PATCH /v1/users/{id} ────

func TestUpdateUserHandler_Success(t *testing.T) {
	app := newTestApplication(&mockUserStore{
		getByIDFn: func(ctx context.Context, id int64) (*store.User, error) { return fakeUser(), nil },
		updateFn:  func(ctx context.Context, u *store.User) error { return nil },
	})

	rr := makeRequest(t, app, http.MethodPatch, "/v1/users/1", map[string]string{"first_name": "Updated"})
	assertStatus(t, rr, http.StatusOK)
}

func TestUpdateUserHandler_NotFound(t *testing.T) {
	app := newTestApplication(&mockUserStore{
		getByIDFn: func(ctx context.Context, id int64) (*store.User, error) { return nil, store.ErrNotFound },
	})

	rr := makeRequest(t, app, http.MethodPatch, "/v1/users/999", map[string]string{"first_name": "Updated"})
	assertStatus(t, rr, http.StatusNotFound)
}

func TestUpdateUserHandler_Conflict(t *testing.T) {
	app := newTestApplication(&mockUserStore{
		getByIDFn: func(ctx context.Context, id int64) (*store.User, error) { return fakeUser(), nil },
		updateFn:  func(ctx context.Context, u *store.User) error { return store.ErrConflict },
	})

	rr := makeRequest(t, app, http.MethodPatch, "/v1/users/1", map[string]string{"username": "existing"})
	assertStatus(t, rr, http.StatusConflict)
}

func TestUpdateUserHandler_InvalidID(t *testing.T) {
	app := newTestApplication(&mockUserStore{})

	rr := makeRequest(t, app, http.MethodPatch, "/v1/users/abc", map[string]string{"first_name": "Updated"})
	assertStatus(t, rr, http.StatusBadRequest)
}

func TestUpdateUserHandler_GetByIDInternalError(t *testing.T) {
	app := newTestApplication(&mockUserStore{
		getByIDFn: func(ctx context.Context, id int64) (*store.User, error) { return nil, errors.New("db down") },
	})

	rr := makeRequest(t, app, http.MethodPatch, "/v1/users/1", map[string]string{"first_name": "Updated"})
	assertStatus(t, rr, http.StatusInternalServerError)
}

func TestUpdateUserHandler_UpdateInternalError(t *testing.T) {
	app := newTestApplication(&mockUserStore{
		getByIDFn: func(ctx context.Context, id int64) (*store.User, error) { return fakeUser(), nil },
		updateFn:  func(ctx context.Context, u *store.User) error { return errors.New("db down") },
	})

	rr := makeRequest(t, app, http.MethodPatch, "/v1/users/1", map[string]string{"first_name": "Updated"})
	assertStatus(t, rr, http.StatusInternalServerError)
}

func TestCreateUserHandler_InternalServerError(t *testing.T) {
	app := newTestApplication(&mockUserStore{
		createFn: func(ctx context.Context, u *store.User) error { return errors.New("database unavailable") },
	})

	body := map[string]string{
		"first_name": "Ensei", "last_name": "Tankado",
		"username": "ensei.tankado", "email": "ensei@example.com", "password": "secret123",
	}

	rr := makeRequest(t, app, http.MethodPost, "/v1/users", body)
	assertStatus(t, rr, http.StatusInternalServerError)
}

func TestListUsersHandler_InternalServerError(t *testing.T) {
	app := newTestApplication(&mockUserStore{
		getAllFn: func(ctx context.Context) ([]store.User, error) { return nil, errors.New("db down") },
	})

	rr := makeRequest(t, app, http.MethodGet, "/v1/users", nil)
	assertStatus(t, rr, http.StatusInternalServerError)
}

func TestGetUserHandler_InternalServerError(t *testing.T) {
	app := newTestApplication(&mockUserStore{
		getByIDFn: func(ctx context.Context, id int64) (*store.User, error) { return nil, errors.New("db down") },
	})

	rr := makeRequest(t, app, http.MethodGet, "/v1/users/1", nil)
	assertStatus(t, rr, http.StatusInternalServerError)
}

func TestCreateUserHandler_PasswordTooLong(t *testing.T) {
	app := newTestApplication(&mockUserStore{})

	body := map[string]string{
		"first_name": "Ensei", "last_name": "Tankado",
		"username": "ensei.tankado", "email": "ensei@example.com",
		"password": string(make([]byte, 73)),
	}

	rr := makeRequest(t, app, http.MethodPost, "/v1/users", body)
	assertStatus(t, rr, http.StatusInternalServerError)
}
