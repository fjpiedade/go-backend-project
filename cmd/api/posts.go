package main

import (
	"errors"
	"net/http"
	"social/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Request types para o Swagger
type CreatePostRequest struct {
	Title   string   `json:"title"   example:"Primeiro Post"`
	Content string   `json:"content" example:"Conteúdo do post"`
	UserID  int64    `json:"user_id" example:"1"`
	Tags    []string `json:"tags"    example:"go,api"`
}

type UpdatePostRequest struct {
	Title   string   `json:"title,omitempty"   example:"Título Atualizado"`
	Content string   `json:"content,omitempty" example:"Novo conteúdo"`
	Tags    []string `json:"tags,omitempty"    example:"go,api,updated"`
}

// @Summary      Create a post
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Param        post body CreatePostRequest true "Post payload"
// @Success      201 {object} store.Post
// @Failure      400 {object} ErrorResponse
// @Failure      422 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /posts [post]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Title   string   `json:"title"`
		Content string   `json:"content"`
		UserID  int64    `json:"user_id"`
		Tags    []string `json:"tags"`
	}

	if err := app.readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if payload.Title == "" || payload.Content == "" || payload.UserID == 0 {
		app.badRequestError(w, r, errors.New("title, content and user_id are required"))
		return
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		UserID:  payload.UserID,
		Tags:    payload.Tags,
	}

	if post.Tags == nil {
		post.Tags = []string{}
	}

	if err := app.store.Posts.Create(r.Context(), post); err != nil {
		if errors.Is(err, store.ErrForeignKey) {
			app.unprocessableEntityError(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, post)
}

// @Summary      Get a post by ID
// @Tags         Posts
// @Produce      json
// @Param        id path int true "Post ID"
// @Success      200 {object} store.Post
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /posts/{id} [get]
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.badRequestError(w, r, errors.New("invalid post id"))
		return
	}

	post, err := app.store.Posts.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			app.notFoundError(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, post)
}

// @Summary      Update a post
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Param        id   path int              true "Post ID"
// @Param        post body UpdatePostRequest true "Fields to update"
// @Success      200 {object} store.Post
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /posts/{id} [patch]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.badRequestError(w, r, errors.New("invalid post id"))
		return
	}

	post, err := app.store.Posts.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			app.notFoundError(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	var payload struct {
		Title   *string  `json:"title"`
		Content *string  `json:"content"`
		Tags    []string `json:"tags"`
	}

	if err := app.readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}
	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Tags != nil {
		post.Tags = payload.Tags
	}

	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			app.notFoundError(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, post)
}

// @Summary      Delete a post
// @Tags         Posts
// @Param        id path int true "Post ID"
// @Success      204
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /posts/{id} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.badRequestError(w, r, errors.New("invalid post id"))
		return
	}

	if err := app.store.Posts.Delete(r.Context(), id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			app.notFoundError(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
