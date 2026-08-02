package main

import (
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"social/internal/store"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type application struct {
	logger  *slog.Logger
	config  config
	store   store.Storage
	metrics *metrics
}

type config struct {
	addr         string
	metricsAddr  string
	otelEndpoint string
	db           dbConfig
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
	r.Use(app.slogMiddleware)
	r.Use(app.prometheusMiddleware)
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

	// envolve o router com instrumentação OTel
	return otelhttp.NewHandler(r, "social-api")
}

func (app *application) slogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()

		defer func() {
			app.logger.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"bytes", ww.BytesWritten(),
				"duration", time.Since(start).String(),
				"request_id", middleware.GetReqID(r.Context()),
			)
		}()

		next.ServeHTTP(ww, r)
	})
}

// func (app *application) run(mux http.Handler) error {

// 	srv := &http.Server{
// 		Addr:         app.config.addr,
// 		Handler:      mux,
// 		WriteTimeout: time.Second * 30,
// 		ReadTimeout:  time.Second * 10,
// 		IdleTimeout:  time.Minute,
// 	}

// 	//log.Printf("Server has started at %s", app.config.addr)
// 	app.logger.Info("server started", "addr", app.config.addr)
// 	return srv.ListenAndServe()
// }

func (app *application) run(mux http.Handler) error {
	metricsListener, err := net.Listen("tcp", app.config.metricsAddr)
	if err != nil {
		return fmt.Errorf("start metrics listener: %w", err)
	}

	metricsServer := &http.Server{
		Handler:           app.metrics.handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		app.logger.Info(
			"metrics server started",
			"addr", app.config.metricsAddr,
			"endpoint", "/metrics",
		)

		if err := metricsServer.Serve(metricsListener); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			app.logger.Error("metrics server error", "error", err)
		}
	}()

	apiServer := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	app.logger.Info("API server started", "addr", app.config.addr)

	return apiServer.ListenAndServe()
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
	app.logger.Error("internal server error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err,
	)
	app.writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn("bad request error",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err,
	)
	app.writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn("not found",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err,
	)
	app.writeJSONError(w, http.StatusNotFound, "not found")
}

func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn("conflict",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err,
	)
	app.writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) unprocessableEntityError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn("unprocessable entity",
		"method", r.Method,
		"path", r.URL.Path,
		"error", err,
	)
	app.writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
}
