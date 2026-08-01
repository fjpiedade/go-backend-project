package main

import (
	"encoding/json"
	"log"
	"net/http"
	"social/internal/store"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type application struct {
	logger *log.Logger
	config config
	store  store.Storage
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Route("/users", func(r chi.Router) {
			r.Post("/", app.createUserHandler)
			r.Get("/", app.listUsersHandler)
			r.Get("/{id}", app.getUserHandler)
			r.Patch("/{id}", app.updateUserHandler)
			r.Delete("/{id}", app.deleteUserHandler)
		})

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)
			r.Get("/{id}", app.getPostHandler)
			r.Patch("/{id}", app.updatePostHandler)
			r.Delete("/{id}", app.deletePostHandler)
		})
	})
	return r
}

func (app *application) run(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server has started at %s", app.config.addr)
	return srv.ListenAndServe()
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_576 //1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

// standard of error on the application
func (app *application) writeJSONError(w http.ResponseWriter, status int, message string) error {
	return app.writeJSON(w, status, map[string]string{"error": message})
}

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Printf("internal server error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	app.writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Printf("bad request error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	app.writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Printf("not found error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	app.writeJSONError(w, http.StatusNotFound, "not found")
}

func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Printf("conflict error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	app.writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) unprocessableEntityError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Printf("unprocessable entity error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	app.writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
}
