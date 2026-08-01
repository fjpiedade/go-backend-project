package main

import "net/http"

// @Summary      Health check
// @Description  Returns server status and version
// @Tags         Health
// @Produce      json
// @Success      200 {object} map[string]string
// @Router       /health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "available",
		"version": "1.0.0",
	}

	if err := app.writeJSON(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
