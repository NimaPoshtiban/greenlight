package main

import (
	"net/http"
	"time"

	"github.com/nimaposhtiban/greenlight/internal/data"
	"github.com/nimaposhtiban/greenlight/internal/validator"
)

// createMovieRequest represents the request body to create a movie.

// createMovieHandler handles the HTTP POST request to create a new movie.
// @Summary Create a new movie
// @Description Create a new movie with the provided details
// @BasePath /
// @Tags movies
// @Accept json
// @Produce json
// @Param request body createMovieRequest true "Request body to create a movie"
// @Success 201 "Created"
// @Failure 400 "Bad Request"
// @Router /v1/movies [post]
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input createMovieRequest

	err := app.readJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	app.writeJson(w, http.StatusCreated, nil, nil)
}

// @Summary Create a new movie
// @Description Create a new movie with the provided details
// @BasePath /
// @Tags movies
// @Accept json
// @Produce json
// @Success 200  {object}  getMovieResult
// @Failure 400 "Bad Request"
// @Failure 404 "not found"
// @Failure 500 "internal server error"
// @Router /v1/movies/{id} [get]
// @Param id   path int true "id"
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
	err = app.writeJson(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
