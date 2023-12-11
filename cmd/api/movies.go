package main

import (
	"net/http"
	"time"

	"github.com/nimaposhtiban/greenlight/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create movie using post"))
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIdParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// this is only for testing purposes, remove this ASAP
	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Need for speed",
		Runtime:   102,
		Genres:    []string{"Action", "Drama"},
		Version:   1,
	}
	err = app.writeJson(w, http.StatusFound, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
