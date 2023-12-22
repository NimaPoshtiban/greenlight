package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Movies interface {
		Insert(*Movie) error
		Get(int64) (*Movie, error)
		Update(*Movie) error
		Delete(int64) error
		GetAll(string, []string, Filters) ([]*Movie, Metadata, error)
	}
	Users interface {
		Insert(*User) error
		Update(*User) error
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
		Users:  UserModel{DB: db},
	}
}
