package main

import (
	"net/http"
)

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
