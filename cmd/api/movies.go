package main

import (
	"errors"
	"fmt"
	"github.com/nimaposhtiban/greenlight/internal/data"
	"github.com/nimaposhtiban/greenlight/internal/validator"
	"net/http"
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
// @Failure 500 "Internal Server Error"
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
	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = app.writeJson(w, http.StatusCreated, envelope{"movie": movie}, headers)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
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
	movie, err := app.models.Movies.Get(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJson(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Update a  movie
// @Description Update a movie with the provided details
// @BasePath /
// @Tags movies
// @Accept json
// @Produce json
// @Success 200  "Ok"
// @Failure 400 "Bad Request"
// @Failure 404 "not found"
// @Failure 500 "Internal Server Error"
// @Router /v1/movies/{id} [put]
// @Param id   path int true "id"
// @Param request body createMovieRequest true "Request body to update a movie"
func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	movie, err := app.models.Movies.Get(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)

		default:
			app.serverErrorResponse(w, r, err)

		}
		return
	}

	var input createMovieRequest

	err = app.readJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie.Title = input.Title
	movie.Year = input.Year
	movie.Runtime = input.Runtime
	movie.Genres = input.Genres

	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Update(movie)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJson(w, http.StatusOK, envelope{"movie": movie}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Delete handles the HTTP Delete request to delete a movie.
// @Summary delete a movie
// @Description Delete a  movie with the provided id
// @BasePath /
// @Tags movies
// @Produce json
// @Success 200 "Ok"
// @Failure 400 "Bad Request"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /v1/movies/{id} [delete]
// @Param id   path int true "id"
func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)

		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJson(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
