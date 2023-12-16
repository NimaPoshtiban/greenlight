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

// this model is used for put & get
type createMovieRequest struct {
	Title   string       `json:"title" validate:"required" maximum:"500"`
	Year    int32        `json:"year" validate:"required" minimum:"1888" `
	Runtime data.Runtime `json:"runtime" validate:"required" swaggertype:"string" format:"utf-8" example:"128 mins"`
	Genres  []string     `json:"genres" validate:"required" maximum:"5"`
}
