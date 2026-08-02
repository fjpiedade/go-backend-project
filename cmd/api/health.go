package main

import (
	"net/http"

	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("social-api")

// @Summary      Health check
// @Description  Returns server status and version
// @Tags         Health
// @Produce      json
// @Success      200 {object} map[string]string
// @Router       /health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	ctx, span := tracer.Start(r.Context(), "healthCheckHandler")
	defer span.End()

	_ = ctx // usar ctx para propagar para chamadas downstream

	app.logger.Info("health check requested",
		"method", r.Method,
		"path", r.URL.Path,
	)

	data := map[string]string{
		"status":  "available",
		"version": "1.0.0",
	}

	if err := app.writeJSON(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
