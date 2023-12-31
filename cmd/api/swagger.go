// types for swaggo API are defined here
package main

import (
	"time"

	"github.com/nimaposhtiban/greenlight/internal/data"
)

type getMovieResult struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   int32     `json:"runtime,string"`
	Genres    []string  `json:"genres"`
	Version   int32     `json:"version"`
}

// this model is used for  get
type createMovieRequest struct {
	Title   string       `json:"title" validate:"required" maximum:"500"`
	Year    int32        `json:"year" validate:"required" minimum:"1888" `
	Runtime data.Runtime `json:"runtime" validate:"required" swaggertype:"string" format:"utf-8" example:"128 mins"`
	Genres  []string     `json:"genres" validate:"required" maximum:"5"`
}

type updateMovieRequest struct {
	Title   *string       `json:"title" maximum:"500"`
	Year    *int32        `json:"year" minimum:"1888"`
	Runtime *data.Runtime `json:"runtime" swaggertype:"string" format:"utf-8" example:"128 mins"`
	Genres  []string      `json:"genres" maximum:"5"`
}

type listMoviesRequest struct {
	Title  string
	Genres []string
	data.Filters
}

type registerUserRequest struct {
	Name     string `json:"name" validate:"required" maximum:"500"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required" maximum:"72"`
}
