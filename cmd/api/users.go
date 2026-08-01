package main

import (
	"errors"
	"net/http"
	"social/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Request/Response types para o Swagger
type CreateUserRequest struct {
	FirstName string `json:"first_name" example:"Ensei"`
	LastName  string `json:"last_name"  example:"Tankado"`
	Username  string `json:"username"   example:"ensei.tankado"`
	Email     string `json:"email"      example:"ensei@example.com"`
	Password  string `json:"password"   example:"secret123"`
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name,omitempty" example:"John"`
	LastName  string `json:"last_name,omitempty"  example:"Phi"`
	Username  string `json:"username,omitempty"   example:"john.phi"`
	Email     string `json:"email,omitempty"      example:"john.phi@example.com"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"not found"`
}

// @Summary      Create a user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user body CreateUserRequest true "User payload"
// @Success      201 {object} store.User
// @Failure      400 {object} ErrorResponse
// @Failure      409 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /users [post]
func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	if err := app.readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if payload.FirstName == "" || payload.LastName == "" ||
		payload.Username == "" || payload.Email == "" || payload.Password == "" {
		app.badRequestError(w, r, errors.New("all fields are required"))
		return
	}

	hashed, err := store.HashPassword(payload.Password)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	user := &store.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Username:  payload.Username,
		Email:     payload.Email,
		Password:  hashed,
	}

	if err := app.store.Users.Create(r.Context(), user); err != nil {
		if errors.Is(err, store.ErrConflict) {
			app.conflictError(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, user)
}

// @Summary      List all users
// @Tags         Users
// @Produce      json
// @Success      200 {array}  store.User
// @Failure      500 {object} ErrorResponse
// @Router       /users [get]
func (app *application) listUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := app.store.Users.GetAll(r.Context())
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	app.writeJSON(w, http.StatusOK, users)
}

// @Summary      Get a user by ID
// @Tags         Users
// @Produce      json
// @Param        id path int true "User ID"
// @Success      200 {object} store.User
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /users/{id} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.badRequestError(w, r, errors.New("invalid user id"))
		return
	}

	user, err := app.store.Users.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			app.notFoundError(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, user)
}

// @Summary      Update a user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path int              true "User ID"
// @Param        user body UpdateUserRequest true "Fields to update"
// @Success      200 {object} store.User
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      409 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /users/{id} [patch]
func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.badRequestError(w, r, errors.New("invalid user id"))
		return
	}

	user, err := app.store.Users.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			app.notFoundError(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	var payload struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Username  *string `json:"username"`
		Email     *string `json:"email"`
	}

	if err := app.readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Atualiza apenas os campos enviados (PATCH parcial)
	if payload.FirstName != nil {
		user.FirstName = *payload.FirstName
	}
	if payload.LastName != nil {
		user.LastName = *payload.LastName
	}
	if payload.Username != nil {
		user.Username = *payload.Username
	}
	if payload.Email != nil {
		user.Email = *payload.Email
	}

	if err := app.store.Users.Update(r.Context(), user); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			app.notFoundError(w, r, err)
			return
		}

		if errors.Is(err, store.ErrConflict) {
			app.conflictError(w, r, err)
			return
		}

		app.internalServerError(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, user)
}

// @Summary      Delete a user
// @Tags         Users
// @Param        id path int true "User ID"
// @Success      204
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /users/{id} [delete]
func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.badRequestError(w, r, errors.New("invalid user id"))
		return
	}

	if err := app.store.Users.Delete(r.Context(), id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			app.notFoundError(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
