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
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
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
// @Failure 404 "Not found"
// @Failure 409 "Edit conflict"
// @Failure 500 "Internal Server Error"
// @Router /v1/movies/{id} [patch]
// @Param id   path int true "id"
// @Param request body updateMovieRequest true "Request body to update a movie"
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

	var input updateMovieRequest

	err = app.readJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}

	if input.Year != nil {
		movie.Year = *input.Year
	}
	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}
	if input.Genres != nil {
		movie.Genres = input.Genres
	}

	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Update(movie)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)

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

// Delete handles the HTTP Delete request to delete a movie.
// @Summary Delete a movie
// @Description Delete a  movie with the provided id
// @BasePath /
// @Tags movies
// @Produce json
// @Success 200 "Ok"
// @Failure 400 "Bad Request"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Param id   path int true "id"
// @Router /v1/movies/{id} [delete]
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

// @Summary Return a list of movies
// @Description returns a list of movies based on provided query string
// @BasePath /
// @Tags movies
// @Produce json
// @Success 200 "Ok"
// @Failure 422 "Invalid request"
// @Failure 500 "Internal Server Error"
// @Param title   query string false "title"
// @Param genres   query string false "genres"
// @Param page   query int false "page"
// @Param page_size   query int false "page_size"
// @Param sort   query string false "sort"
// @Router /v1/movies [get]
func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	var input listMoviesRequest

	v := validator.New()

	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	input.Filters.Sort = app.readString(qs, "sort", "id")

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	movies, metadata, err := app.models.Movies.GetAll(input.Title, input.Genres, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envelope{"movies": movies, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
