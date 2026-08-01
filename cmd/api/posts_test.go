package main

import (
	"context"
	"errors"
	"net/http"
	"social/internal/store"
	"testing"
	"time"
)

// ─── Fixture ────────────────────────────────────────────────────────────────

func fakePost() *store.Post {
	return &store.Post{
		ID:        1,
		Title:     "Primeiro Post",
		Content:   "Conteúdo do post",
		UserID:    1,
		Tags:      []string{"go", "api"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// func newTestAppWithPost(postStore store.PostStorageInterface) *application {
// 	return newTestApplication(&mockUserStore{}, postStore)
// }

// ─── POST /v1/posts ──────────────────────────────────────────────────────────

func TestCreatePostHandler_Success(t *testing.T) {
	app := newTestAppWithPost(&mockPostStore{
		createFn: func(ctx context.Context, p *store.Post) error {
			p.ID = 1
			p.CreatedAt = time.Now()
			p.UpdatedAt = time.Now()
			return nil
		},
	})

	body := map[string]any{
		"title":   "Primeiro Post",
		"content": "Conteúdo do post",
		"user_id": 1,
		"tags":    []string{"go", "api"},
	}

	rr := makeRequest(t, app, http.MethodPost, "/v1/posts", body)
	assertStatus(t, rr, http.StatusCreated)
}

func TestCreatePostHandler_MissingFields(t *testing.T) {
	app := newTestAppWithPost(&mockPostStore{})

	body := map[string]any{"title": "Só título"}
	rr := makeRequest(t, app, http.MethodPost, "/v1/posts", body)
	assertStatus(t, rr, http.StatusBadRequest)
}

func TestCreatePostHandler_InvalidJSON(t *testing.T) {
	app := newTestAppWithPost(&mockPostStore{})

	rr := makeRequestRaw(t, app, http.MethodPost, "/v1/posts", `{invalid}`)
	assertStatus(t, rr, http.StatusBadRequest)
}

func TestCreatePostHandler_InternalServerError(t *testing.T) {
	app := newTestAppWithPost(&mockPostStore{
		createFn: func(ctx context.Context, p *store.Post) error {
			return errors.New("db down")
		},
	})

	body := map[string]any{
		"title":   "Post",
		"content": "Conteúdo",
		"user_id": 1,
	}

	rr := makeRequest(t, app, http.MethodPost, "/v1/posts", body)
	assertStatus(t, rr, http.StatusInternalServerError)
}

// ─── GET /v1/posts/{id} ──────────────────────────────────────────────────────

func TestGetPostHandler_Success(t *testing.T) {
	app := newTestAppWithPost(&mockPostStore{
		getByIDFn: func(ctx context.Context, id int64) (*store.Post, error) {
			return fakePost(), nil
		},
	})

	rr := makeRequest(t, app, http.MethodGet, "/v1/posts/1", nil)
	assertStatus(t, rr, http.StatusOK)
}

func TestGetPostHandler_NotFound(t *testing.T) {
	app := newTestAppWithPost(&mockPostStore{
		getByIDFn: func(ctx context.Context, id int64) (*store.Post, error) {
			return nil, store.ErrNotFound
		},
	})

	rr := makeRequest(t, app, http.MethodGet, "/v1/posts/999", nil)
	assertStatus(t, rr, http.StatusNotFound)
}

func TestGetPostHandler_InvalidID(t *testing.T) {
	app := newTestAppWithPost(&mockPostStore{})

	rr := makeRequest(t, app, http.MethodGet, "/v1/posts/abc", nil)
	assertStatus(t, rr, http.StatusBadRequest)
}

func TestGetPostHandler_InternalServerError(t *testing.T) {
	app := newTestAppWithPost(&mockPostStore{
		getByIDFn: func(ctx context.Context, id int64) (*store.Post, error) {
			return nil, errors.New("db down")
		},
	})

	rr := makeRequest(t, app, http.MethodGet, "/v1/posts/1", nil)
	assertStatus(t, rr, http.StatusInternalServerError)
}

// ─── PATCH /v1/posts/{id} ────────────────────────────────────────────────────

func TestUpdatePostHandler_Success(t *testing.T) {
	app := newTestAppWithPost(&mockPostStore{
		getByIDFn: func(ctx context.Context, id int64) (*store.Post, error) {
			return fakePost(), nil
		},
		updateFn: func(ctx context.Context, p *store.Post) error {
			return nil
		},
	})

	body := map[string]any{"title": "Título Atualizado"}
	rr := makeRequest(t, app, http.MethodPatch, "/v1/posts/1", body)
	assertStatus(t, rr, http.StatusOK)
}

func TestUpdatePostHandler_NotFound(t *testing.T) {
	app := newTestAppWithPost(&mockPostStore{
		getByIDFn: func(ctx context.Context, id int64) (*store.Post, error) {
			return nil, store.ErrNotFound
		},
	})

	body := map[string]any{"title": "Título"}
	rr := makeRequest(t, app, http.MethodPatch, "/v1/posts/999", body)
	assertStatus(t, rr, http.StatusNotFound)
}

func TestUpdatePostHandler_InvalidID(t *testing.T) {
	app := newTestAppWithPost(&mockPostStore{})

	body := map[string]any{"title": "Título"}
	rr := makeRequest(t, app, http.MethodPatch, "/v1/posts/abc", body)
	assertStatus(t, rr, http.StatusBadRequest)
}

func TestUpdatePostHandler_GetByIDInternalError(t *testing.T) {
	app := newTestAppWithPost(&mockPostStore{
		getByIDFn: func(ctx context.Context, id int64) (*store.Post, error) {
			return nil, errors.New("db down")
		},
	})

	body := map[string]any{"title": "Título"}
	rr := makeRequest(t, app, http.MethodPatch, "/v1/posts/1", body)
	assertStatus(t, rr, http.StatusInternalServerError)
}

func TestUpdatePostHandler_UpdateInternalError(t *testing.T) {
	app := newTestAppWithPost(&mockPostStore{
		getByIDFn: func(ctx context.Context, id int64) (*store.Post, error) {
			return fakePost(), nil
		},
		updateFn: func(ctx context.Context, p *store.Post) error {
			return errors.New("db down")
		},
	})

	body := map[string]any{"title": "Título"}
	rr := makeRequest(t, app, http.MethodPatch, "/v1/posts/1", body)
	assertStatus(t, rr, http.StatusInternalServerError)
}

// ─── DELETE /v1/posts/{id} ───────────────────────────────────────────────────

func TestDeletePostHandler_Success(t *testing.T) {
	app := newTestAppWithPost(&mockPostStore{
		deleteFn: func(ctx context.Context, id int64) error {
			return nil
		},
	})

	rr := makeRequest(t, app, http.MethodDelete, "/v1/posts/1", nil)
	assertStatus(t, rr, http.StatusNoContent)
}

func TestDeletePostHandler_NotFound(t *testing.T) {
	app := newTestAppWithPost(&mockPostStore{
		deleteFn: func(ctx context.Context, id int64) error {
			return store.ErrNotFound
		},
	})

	rr := makeRequest(t, app, http.MethodDelete, "/v1/posts/999", nil)
	assertStatus(t, rr, http.StatusNotFound)
}

func TestDeletePostHandler_InvalidID(t *testing.T) {
	app := newTestAppWithPost(&mockPostStore{})

	rr := makeRequest(t, app, http.MethodDelete, "/v1/posts/abc", nil)
	assertStatus(t, rr, http.StatusBadRequest)
}

func TestDeletePostHandler_InternalServerError(t *testing.T) {
	app := newTestAppWithPost(&mockPostStore{
		deleteFn: func(ctx context.Context, id int64) error {
			return errors.New("db down")
		},
	})

	rr := makeRequest(t, app, http.MethodDelete, "/v1/posts/1", nil)
	assertStatus(t, rr, http.StatusInternalServerError)
}
