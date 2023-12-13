package main

import (
	"net/http"
)

// @title healthcheck
// @version 1.0
// @description this endpoint shows the status of API
// @BasePath /
// @produce json
// @schemes http https
// @Tags Health
// @Success 200
// @router /v1/healthcheck [get]
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	status := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}
	err := app.writeJson(w, http.StatusAccepted, envelope{"health": status}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
